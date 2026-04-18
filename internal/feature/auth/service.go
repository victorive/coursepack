package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/victorive/klaro-backend/internal/cache"
	"github.com/victorive/klaro-backend/internal/domain"
	"github.com/victorive/klaro-backend/internal/email"
	"github.com/victorive/klaro-backend/internal/pkg/apperr"
	"github.com/victorive/klaro-backend/internal/pkg/hash"
	"github.com/victorive/klaro-backend/internal/pkg/job"
	jwt "github.com/victorive/klaro-backend/internal/pkg/jwt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	OTPExpirationMinutes     int
	OTPMaxAttempts           int
	RefreshExpirationMinutes int
	GoogleClientID           string
	GoogleClientSecret       string
	GoogleRedirectURL        string
}

type Service interface {
	SignUp(ctx context.Context, req *SignUpRequest) (*TokensResponse, error)
	SignIn(ctx context.Context, req *SignInRequest) (*TokensResponse, error)
	GoogleInitiate(ctx context.Context) (redirectURL string, err error)
	GoogleCallback(ctx context.Context, code, state string) (*TokensResponse, error)
	VerifyOTP(ctx context.Context, req *VerifyOTPRequest) error
	ResendOTP(ctx context.Context, req *ResendOTPRequest) error
	Refresh(ctx context.Context, req *RefreshRequest) (*TokensResponse, error)
	ValidateToken(ctx context.Context, token string) (sub string, err error)
	Logout(ctx context.Context, accessToken string, req *LogoutRequest) error
	ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req *ResetPasswordRequest) error
	Me(ctx context.Context, uuid string) (*UserResponse, error)
}

type service struct {
	repo   Repository
	cache  cache.Cache
	mailer email.Mailer
	jwt    *jwt.Manager
	cfg    Config
	logger *slog.Logger
}

func NewService(
	repo Repository,
	c cache.Cache,
	mailer email.Mailer,
	jwt *jwt.Manager,
	cfg Config,
	logger *slog.Logger,
) Service {
	return &service{
		repo:   repo,
		cache:  c,
		mailer: mailer,
		jwt:    jwt,
		cfg:    cfg,
		logger: logger,
	}
}

const (
	otpKeyPrefix         = "otp:"
	otpAttemptsKeyPrefix = "otp_attempts:"
	otpResendKeyPrefix   = "otp_resend:"
	refreshKeyPrefix     = "refresh_token:"
	pwdResetKeyPrefix    = "pwd_reset:"
	googleStateKeyPrefix = "google_state:"
)

func (s *service) issueTokens(ctx context.Context, u *domain.User) (*TokensResponse, error) {
	access, err := s.jwt.Sign(u.UUID)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}

	refresh, err := hash.GenerateSecureToken()
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}

	ttl := time.Duration(s.cfg.RefreshExpirationMinutes) * time.Minute
	userIDStr := strconv.FormatInt(u.ID, 10)
	if err := s.cache.Set(ctx, refreshKeyPrefix+refresh, userIDStr, ttl); err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}

	return &TokensResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresAt:    int(time.Now().Add(s.jwt.GetExpirationTime()).Unix()),
	}, nil
}

func (s *service) sendOTPAsync(ctx context.Context, userEmail string) {
	expiry := s.cfg.OTPExpirationMinutes
	job.FireAndForget(ctx, "send_otp", func(bgCtx context.Context) error {
		otp, err := hash.GenerateOTP()
		if err != nil {
			return fmt.Errorf("generate otp: %w", err)
		}
		ttl := time.Duration(expiry) * time.Minute
		if err := s.cache.Set(bgCtx, otpKeyPrefix+userEmail, otp, ttl); err != nil {
			return fmt.Errorf("store otp: %w", err)
		}
		if err := s.mailer.Send(bgCtx, userEmail, "Your verification code", "otp.html", map[string]any{
			"OTP":    otp,
			"Expiry": expiry,
		}); err != nil {
			return fmt.Errorf("send otp email: %w", err)
		}
		return nil
	})
}

func (s *service) SignUp(ctx context.Context, req *SignUpRequest) (*TokensResponse, error) {
	existing, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if existing != nil {
		return nil, apperr.New("email already registered", http.StatusConflict)
	}

	ph, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}

	u := &domain.User{
		UUID:     uuid.New().String(),
		Email:    req.Email,
		FullName: req.FullName,
		Password: &ph,
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}

	s.sendOTPAsync(ctx, u.Email)

	return s.issueTokens(ctx, u)
}

func (s *service) SignIn(ctx context.Context, req *SignInRequest) (*TokensResponse, error) {
	u, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if u == nil {
		return nil, apperr.New("invalid email or password", http.StatusUnauthorized)
	}
	if u.Provider != nil {
		return nil, apperr.New("account registered with a different sign-in method", http.StatusUnprocessableEntity)
	}
	if u.Password == nil {
		return nil, apperr.New("invalid email or password", http.StatusUnauthorized)
	}

	ok, err := hash.VerifyPassword(*u.Password, req.Password)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if !ok {
		return nil, apperr.New("invalid email or password", http.StatusUnauthorized)
	}

	return s.issueTokens(ctx, u)
}

func (s *service) googleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     s.cfg.GoogleClientID,
		ClientSecret: s.cfg.GoogleClientSecret,
		RedirectURL:  s.cfg.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func (s *service) GoogleInitiate(ctx context.Context) (string, error) {
	state, err := hash.GenerateSecureToken()
	if err != nil {
		return "", apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if err := s.cache.Set(ctx, googleStateKeyPrefix+state, "1", 10*time.Minute); err != nil {
		return "", apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	return s.googleOAuthConfig().AuthCodeURL(state, oauth2.AccessTypeOnline), nil
}

type googleUserInfo struct {
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func (s *service) GoogleCallback(ctx context.Context, code, state string) (*TokensResponse, error) {
	val, err := s.cache.Get(ctx, googleStateKeyPrefix+state)
	if err != nil || val == "" {
		return nil, apperr.New("invalid or expired state", http.StatusUnauthorized)
	}
	_ = s.cache.Delete(ctx, googleStateKeyPrefix+state)

	oauthCfg := s.googleOAuthConfig()
	googleToken, err := oauthCfg.Exchange(ctx, code)
	if err != nil {
		return nil, apperr.New("failed to exchange authorization code", http.StatusUnauthorized).Wrap(err)
	}

	client := oauthCfg.Client(ctx, googleToken)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apperr.New("failed to fetch Google user info", http.StatusInternalServerError)
	}

	var info googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}

	u, err := s.repo.FindByEmail(ctx, info.Email)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}

	if u == nil {
		picture := (*string)(nil)
		if info.Picture != "" {
			picture = &info.Picture
		}
		provider := string(domain.ProviderGoogle)
		u = &domain.User{
			UUID:          uuid.New().String(),
			Email:         info.Email,
			FullName:      info.Name,
			AvatarURL:     picture,
			Provider:      &provider,
			EmailVerified: info.VerifiedEmail,
		}
		if err := s.repo.Create(ctx, u); err != nil {
			return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
		}
	} else if u.Provider == nil || *u.Provider != string(domain.ProviderGoogle) {
		return nil, apperr.New("account registered with a different sign-in method", http.StatusUnprocessableEntity)
	}

	return s.issueTokens(ctx, u)
}

func (s *service) VerifyOTP(ctx context.Context, req *VerifyOTPRequest) error {
	attemptsKey := otpAttemptsKeyPrefix + req.Email
	attemptsStr, _ := s.cache.Get(ctx, attemptsKey)
	attempts, _ := strconv.Atoi(attemptsStr)
	if attempts >= s.cfg.OTPMaxAttempts {
		return apperr.New("too many incorrect OTP attempts", http.StatusTooManyRequests)
	}

	stored, err := s.cache.Get(ctx, otpKeyPrefix+req.Email)
	if err != nil || stored == "" || stored != req.OTP {
		ttl := time.Duration(s.cfg.OTPExpirationMinutes) * time.Minute
		_ = s.cache.Set(ctx, attemptsKey, strconv.Itoa(attempts+1), ttl)
		return apperr.New("invalid or expired OTP", http.StatusUnprocessableEntity)
	}

	u, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if u == nil {
		return apperr.New("user not found", http.StatusNotFound)
	}

	if err := s.repo.UpdateEmailVerified(ctx, u.ID, true); err != nil {
		return apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	_ = s.cache.Delete(ctx, otpKeyPrefix+req.Email)
	_ = s.cache.Delete(ctx, attemptsKey)
	return nil
}

func (s *service) ResendOTP(ctx context.Context, req *ResendOTPRequest) error {
	resendKey := otpResendKeyPrefix + req.Email
	countStr, _ := s.cache.Get(ctx, resendKey)
	count, _ := strconv.Atoi(countStr)
	if count >= 3 {
		return apperr.New("please wait before requesting a new OTP", http.StatusTooManyRequests)
	}

	newCount := strconv.Itoa(count + 1)
	_ = s.cache.Set(ctx, resendKey, newCount, time.Minute)

	email := req.Email
	job.FireAndForget(ctx, "resend_otp", func(bgCtx context.Context) error {
		u, err := s.repo.FindByEmail(bgCtx, email)
		if err != nil {
			return fmt.Errorf("find user: %w", err)
		}
		if u == nil {
			return nil
		}
		s.sendOTPAsync(bgCtx, u.Email)
		return nil
	})

	return nil
}

func (s *service) Refresh(ctx context.Context, req *RefreshRequest) (*TokensResponse, error) {
	key := refreshKeyPrefix + req.RefreshToken
	userIDStr, err := s.cache.Get(ctx, key)
	if err != nil || userIDStr == "" {
		return nil, apperr.New("invalid token", http.StatusUnauthorized)
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return nil, apperr.New("invalid token", http.StatusUnauthorized)
	}

	_ = s.cache.Delete(ctx, key)

	u, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if u == nil {
		return nil, apperr.New("user not found", http.StatusNotFound)
	}

	return s.issueTokens(ctx, u)
}

func (s *service) ValidateToken(ctx context.Context, token string) (string, error) {
	claims, err := s.jwt.Verify(token)
	if err != nil {
		return "", apperr.New("unauthenticated", http.StatusUnauthorized)
	}

	blacklisted, _ := s.cache.Exists(ctx, "blacklist:"+token)
	if blacklisted {
		return "", apperr.New("unauthenticated", http.StatusUnauthorized)
	}
	
	return claims.Subject, nil
}

func (s *service) Logout(ctx context.Context, accessToken string, req *LogoutRequest) error {
	claims, err := s.jwt.Verify(accessToken)
	if err != nil {
		return apperr.New("unauthenticated", http.StatusUnauthorized)
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl > 0 {
		_, _ = s.cache.SetNX(ctx, "blacklist:"+accessToken, "1", ttl)
	}

	_ = s.cache.Delete(ctx, refreshKeyPrefix+req.RefreshToken)

	return nil
}

func (s *service) ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) error {
	reqEmail := req.Email
	job.FireAndForget(ctx, "forgot_password", func(bgCtx context.Context) error {
		u, err := s.repo.FindByEmail(bgCtx, reqEmail)
		if err != nil {
			return fmt.Errorf("find user: %w", err)
		}
		if u == nil {
			return nil
		}
		token, err := hash.GenerateSecureToken()
		if err != nil {
			return fmt.Errorf("generate reset token: %w", err)
		}
		if err := s.cache.Set(bgCtx, pwdResetKeyPrefix+token, strconv.FormatInt(u.ID, 10), 15*time.Minute); err != nil {
			return fmt.Errorf("store reset token: %w", err)
		}
		if err := s.mailer.Send(bgCtx, u.Email, "Reset your password", "password_reset.html", map[string]any{
			"Token": token,
		}); err != nil {
			return fmt.Errorf("send reset email: %w", err)
		}
		return nil
	})
	return nil
}

func (s *service) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	key := pwdResetKeyPrefix + req.Token
	userIDStr, err := s.cache.Get(ctx, key)
	if err != nil || userIDStr == "" {
		return apperr.New("invalid token", http.StatusUnauthorized)
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return apperr.New("invalid token", http.StatusUnauthorized)
	}

	ph, err := hash.HashPassword(req.Password)
	if err != nil {
		return apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}

	if err := s.repo.UpdatePassword(ctx, userID, ph); err != nil {
		return apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}

	_ = s.cache.Delete(ctx, key)
	return nil
}

func (s *service) Me(ctx context.Context, userUUID string) (*UserResponse, error) {
	u, err := s.repo.FindByUUID(ctx, userUUID)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if u == nil {
		return nil, apperr.New("user not found", http.StatusNotFound)
	}
	return toUserResponse(u), nil
}

func toUserResponse(u *domain.User) *UserResponse {
	return &UserResponse{
		ID:            u.ID,
		UUID:          u.UUID,
		Email:         u.Email,
		FullName:      u.FullName,
		Phone:         u.Phone,
		AvatarURL:     u.AvatarURL,
		Provider:      u.Provider,
		EmailVerified: u.EmailVerified,
	}
}
