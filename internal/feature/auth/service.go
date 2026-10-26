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
	"github.com/victorive/coursepack/internal/cache"
	"github.com/victorive/coursepack/internal/domain"
	"github.com/victorive/coursepack/internal/email"
	"github.com/victorive/coursepack/internal/feature/university"
	"github.com/victorive/coursepack/internal/permissions"
	"github.com/victorive/coursepack/internal/pkg/apperr"
	"github.com/victorive/coursepack/internal/pkg/hash"
	"github.com/victorive/coursepack/internal/pkg/job"
	jwt "github.com/victorive/coursepack/internal/pkg/jwt"
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
	FrontendURL              string
}

type Service interface {
	SignUp(ctx context.Context, req *SignUpRequest) (*TokensResponse, error)
	SignIn(ctx context.Context, req *SignInRequest) (*TokensResponse, error)
	GoogleInitiate(ctx context.Context) (redirectURL string, err error)
	GoogleCallback(ctx context.Context, code, state string) (oneTimeCode string, err error)
	ExchangeGoogleCode(ctx context.Context, code string) (*TokensResponse, error)
	VerifyOTP(ctx context.Context, req *VerifyOTPRequest) error
	ResendOTP(ctx context.Context, req *ResendOTPRequest) error
	Refresh(ctx context.Context, req *RefreshRequest) (*TokensResponse, error)
	ValidateToken(ctx context.Context, token string) (sub string, err error)
	Logout(ctx context.Context, accessToken string, req *LogoutRequest) error
	ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req *ResetPasswordRequest) error
	Me(ctx context.Context, actor permissions.Actor) (*UserResponse, error)
	CheckEmailVerified(ctx context.Context, uuid string) (emailVerified bool, isOAuthUser bool, err error)
}

type service struct {
	repo    Repository
	cache   cache.Cache
	mailer  email.Mailer
	jwt     *jwt.Manager
	cfg     Config
	permSvc permissions.Service
	logger  *slog.Logger
}

func NewService(
	repo Repository,
	c cache.Cache,
	mailer email.Mailer,
	jwt *jwt.Manager,
	cfg Config,
	permSvc permissions.Service,
	logger *slog.Logger,
) Service {
	return &service{
		repo:    repo,
		cache:   c,
		mailer:  mailer,
		jwt:     jwt,
		cfg:     cfg,
		permSvc: permSvc,
		logger:  logger,
	}
}

const (
	otpKeyPrefix         = "otp:"
	otpAttemptsKeyPrefix = "otp_attempts:"
	otpResendKeyPrefix   = "otp_resend:"
	refreshKeyPrefix     = "refresh_token:"
	pwdResetKeyPrefix    = "pwd_reset:"
	googleStateKeyPrefix = "google_state:"
	googleCodeKeyPrefix  = "google_code:"
)

func (s *service) issueTokens(ctx context.Context, u *domain.User) (*TokensResponse, error) {
	access, err := s.jwt.Sign(u.UUID)
	if err != nil {
		return nil, apperr.Internal(err)
	}

	refresh, err := hash.GenerateSecureToken()
	if err != nil {
		return nil, apperr.Internal(err)
	}

	ttl := time.Duration(s.cfg.RefreshExpirationMinutes) * time.Minute
	userIDStr := strconv.FormatInt(u.ID, 10)
	if err := s.cache.Set(ctx, refreshKeyPrefix+refresh, userIDStr, ttl); err != nil {
		return nil, apperr.Internal(err)
	}

	return &TokensResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresAt:    int(time.Now().Add(s.jwt.GetExpirationTime()).Unix()),
	}, nil
}

func (s *service) buildUserResponse(ctx context.Context, u *domain.User) (*UserResponse, error) {
	roles, err := s.permSvc.GetUserRoles(ctx, u.UUID)
	if err != nil {
		return nil, apperr.Internal(err)
	}

	perms, err := s.permSvc.GetUserPermissions(ctx, u.UUID)
	if err != nil {
		return nil, apperr.Internal(err)
	}

	var uniInfo *university.UniversityResponse
	if u.University != nil {
		uniInfo = university.ToUniversityResponse(u.University)
	}

	return &UserResponse{
		UUID:          u.UUID,
		Email:         u.Email,
		FullName:      u.FullName,
		Phone:         u.Phone,
		AvatarURL:     u.AvatarURL,
		Provider:      u.Provider,
		EmailVerified: u.EmailVerified,
		University:    uniInfo,
		Roles:         roles,
		Permissions:   perms,
	}, nil
}

func (s *service) toAuthResponse(ctx context.Context, u *domain.User) (*TokensResponse, error) {
	tokens, err := s.issueTokens(ctx, u)
	if err != nil {
		return nil, err
	}

	userResp, err := s.buildUserResponse(ctx, u)
	if err != nil {
		return nil, err
	}

	tokens.User = userResp

	return tokens, nil
}

const defaultOTPScope = "email_verification"

func (s *service) sendOTPAsync(ctx context.Context, userEmail, scope string) {
	expiry := s.cfg.OTPExpirationMinutes
	logger := s.logger
	otpKey := otpKeyPrefix + scope + ":" + userEmail
	job.FireAndForget(ctx, "send_otp", func(bgCtx context.Context) error {
		otp, err := hash.GenerateOTP()
		if err != nil {
			return fmt.Errorf("generate otp: %w", err)
		}
		// TODO: remove once SMTP is configured
		logger.Info("DEV OTP", "email", userEmail, "otp", otp)
		ttl := time.Duration(expiry) * time.Minute
		if err := s.cache.Set(bgCtx, otpKey, otp, ttl); err != nil {
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
		return nil, apperr.Internal(err)
	}
	if existing != nil {
		return nil, apperr.New("email already registered", http.StatusConflict)
	}

	ph, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, apperr.Internal(err)
	}

	u := &domain.User{
		UUID:     uuid.New().String(),
		Email:    req.Email,
		FullName: req.FullName,
		Password: &ph,
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, apperr.Internal(err)
	}

	if err := s.permSvc.AssignGlobalRole(ctx, u.ID, domain.RoleStudent); err != nil {
		return nil, apperr.Internal(err)
	}

	s.sendOTPAsync(ctx, u.Email, defaultOTPScope)

	return s.toAuthResponse(ctx, u)
}

func (s *service) SignIn(ctx context.Context, req *SignInRequest) (*TokensResponse, error) {
	u, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperr.Internal(err)
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
		return nil, apperr.Internal(err)
	}
	if !ok {
		return nil, apperr.New("invalid email or password", http.StatusUnauthorized)
	}

	return s.toAuthResponse(ctx, u)
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
	if s.cfg.GoogleClientID == "" || s.cfg.GoogleClientSecret == "" {
		return "", apperr.New("Google OAuth is not configured", http.StatusInternalServerError)
	}
	state, err := hash.GenerateSecureToken()
	if err != nil {
		return "", apperr.Internal(err)
	}
	if err := s.cache.Set(ctx, googleStateKeyPrefix+state, "1", 10*time.Minute); err != nil {
		return "", apperr.Internal(err)
	}
	return s.googleOAuthConfig().AuthCodeURL(state, oauth2.AccessTypeOnline), nil
}

type googleUserInfo struct {
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func (s *service) GoogleCallback(ctx context.Context, code, state string) (string, error) {
	val, err := s.cache.Get(ctx, googleStateKeyPrefix+state)
	if err != nil || val == "" {
		return "", apperr.New("invalid or expired state", http.StatusUnauthorized)
	}
	_ = s.cache.Delete(ctx, googleStateKeyPrefix+state)

	oauthCfg := s.googleOAuthConfig()
	googleToken, err := oauthCfg.Exchange(ctx, code)
	if err != nil {
		return "", apperr.New("failed to exchange authorization code", http.StatusUnauthorized).Wrap(err)
	}

	client := oauthCfg.Client(ctx, googleToken)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return "", apperr.Internal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", apperr.New("failed to fetch Google user info", http.StatusInternalServerError)
	}

	var info googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return "", apperr.Internal(err)
	}

	u, err := s.repo.FindByEmail(ctx, info.Email)
	if err != nil {
		return "", apperr.Internal(err)
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
			return "", apperr.Internal(err)
		}
		if err := s.permSvc.AssignGlobalRole(ctx, u.ID, domain.RoleStudent); err != nil {
			return "", apperr.Internal(err)
		}
	} else if u.Provider == nil || *u.Provider != string(domain.ProviderGoogle) {
		return "", apperr.New("account registered with a different sign-in method", http.StatusUnprocessableEntity)
	} else {
		picture := (*string)(nil)
		if info.Picture != "" {
			picture = &info.Picture
		}
		if err := s.repo.UpdateGoogleProfile(ctx, u.ID, info.Name, picture); err != nil {
			s.logger.Error("failed to update google profile", "user_id", u.ID, "error", err)
		}
	}

	tokens, err := s.toAuthResponse(ctx, u)
	if err != nil {
		return "", err
	}

	exchangeCode, err := hash.GenerateSecureToken()
	if err != nil {
		return "", apperr.Internal(err)
	}

	if err := s.cache.SetJSON(ctx, googleCodeKeyPrefix+exchangeCode, tokens, 5*time.Minute); err != nil {
		return "", apperr.Internal(err)
	}

	return exchangeCode, nil
}

func (s *service) ExchangeGoogleCode(ctx context.Context, code string) (*TokensResponse, error) {
	var tokens TokensResponse
	if err := s.cache.GetJSON(ctx, googleCodeKeyPrefix+code, &tokens); err != nil {
		return nil, apperr.New("invalid or expired code", http.StatusUnauthorized)
	}

	_ = s.cache.Delete(ctx, googleCodeKeyPrefix+code)
	return &tokens, nil
}

func (s *service) VerifyOTP(ctx context.Context, req *VerifyOTPRequest) error {
	scope := defaultOTPScope
	attemptsKey := otpAttemptsKeyPrefix + scope + ":" + req.Email
	attemptsStr, _ := s.cache.Get(ctx, attemptsKey)
	attempts, _ := strconv.Atoi(attemptsStr)
	if attempts >= s.cfg.OTPMaxAttempts {
		return apperr.New("too many incorrect OTP attempts", http.StatusTooManyRequests)
	}

	otpKey := otpKeyPrefix + scope + ":" + req.Email
	stored, err := s.cache.Get(ctx, otpKey)
	if err != nil || stored == "" || stored != req.OTP {
		ttl := time.Duration(s.cfg.OTPExpirationMinutes) * time.Minute
		_ = s.cache.Set(ctx, attemptsKey, strconv.Itoa(attempts+1), ttl)
		return apperr.New("invalid or expired OTP", http.StatusUnprocessableEntity)
	}

	u, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return apperr.Internal(err)
	}
	if u == nil {
		return apperr.NotFound("user")
	}

	if err := s.repo.UpdateEmailVerified(ctx, u.ID, true); err != nil {
		return apperr.Internal(err)
	}
	_ = s.cache.Delete(ctx, otpKey)
	_ = s.cache.Delete(ctx, attemptsKey)
	return nil
}

func (s *service) ResendOTP(ctx context.Context, req *ResendOTPRequest) error {
	scope := defaultOTPScope

	u, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return apperr.Internal(err)
	}

	if u != nil && u.EmailVerified && scope == defaultOTPScope {
		return apperr.New("email already verified", http.StatusConflict)
	}

	resendKey := otpResendKeyPrefix + scope + ":" + req.Email
	countStr, _ := s.cache.Get(ctx, resendKey)
	count, _ := strconv.Atoi(countStr)
	if count >= 5 {
		return apperr.New("please wait before requesting a new OTP", http.StatusTooManyRequests)
	}

	newCount := strconv.Itoa(count + 1)
	_ = s.cache.Set(ctx, resendKey, newCount, time.Hour)

	if u != nil {
		s.sendOTPAsync(ctx, u.Email, scope)
	}

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
		return nil, apperr.Internal(err)
	}
	if u == nil {
		return nil, apperr.NotFound("user")
	}

	return s.toAuthResponse(ctx, u)
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
	frontendURL := s.cfg.FrontendURL
	logger := s.logger
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
		resetLink := frontendURL + "/reset-password?token=" + token
		// TODO: remove once SMTP is configured
		logger.Info("DEV password reset link", "email", u.Email, "link", resetLink)
		if err := s.mailer.Send(bgCtx, u.Email, "Reset your password", "password_reset.html", map[string]any{
			"ResetLink": resetLink,
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
		return apperr.Internal(err)
	}

	if err := s.repo.UpdatePassword(ctx, userID, ph); err != nil {
		return apperr.Internal(err)
	}

	_ = s.cache.Delete(ctx, key)
	return nil
}

func (s *service) CheckEmailVerified(ctx context.Context, userUUID string) (bool, bool, error) {
	u, err := s.repo.FindByUUID(ctx, userUUID)
	if err != nil {
		return false, false, err
	}
	if u == nil {
		return false, false, nil
	}

	return u.EmailVerified, u.Provider != nil, nil
}

func (s *service) Me(ctx context.Context, actor permissions.Actor) (*UserResponse, error) {
	u, err := s.repo.FindByUUID(ctx, actor.UUID)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if u == nil {
		return nil, apperr.NotFound("user")
	}
	return s.buildUserResponse(ctx, u)
}
