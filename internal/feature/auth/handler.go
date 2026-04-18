package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/middleware"
	"github.com/victorive/klaro-backend/internal/pkg/apperr"
	"github.com/victorive/klaro-backend/internal/pkg/response"
	"github.com/victorive/klaro-backend/internal/pkg/validator"
)

type Handler struct {
	svc               Service
	googleFrontendURL string
}

func NewHandler(svc Service, googleFrontendURL string) *Handler {
	return &Handler{svc: svc, googleFrontendURL: googleFrontendURL}
}

// SignUp godoc
// @Summary      Register a new user
// @Description  Creates a new local account and returns access + refresh tokens. An OTP is sent via email asynchronously.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body SignUpRequest true "Sign-up payload"
// @Success      201 {object} response.Response{data=TokensResponse}
// @Failure      409 {object} response.Response "Email already taken"
// @Failure      422 {object} response.Response "Validation error"
// @Router       /auth/signup [post]
func (h *Handler) SignUp(c *fiber.Ctx) error {
	var req SignUpRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	tokens, err := h.svc.SignUp(c.Context(), &req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusCreated, true, "account created", tokens)
}

// SignIn godoc
// @Summary      Sign in with email and password
// @Description  Authenticates a local-provider user and returns tokens.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body SignInRequest true "Sign-in payload"
// @Success      200 {object} response.Response{data=TokensResponse}
// @Failure      401 {object} response.Response "Invalid credentials"
// @Failure      422 {object} response.Response "Validation error"
// @Router       /auth/signin [post]
func (h *Handler) SignIn(c *fiber.Ctx) error {
	var req SignInRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	tokens, err := h.svc.SignIn(c.Context(), &req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "login successful", tokens)
}

// GoogleInitiate godoc
// @Summary      Start Google OAuth2
// @Description  Redirects the user to Google's authorization page.
// @Tags         Auth
// @Success      302
// @Failure      500 {object} response.Response
// @Router       /auth/google [get]
func (h *Handler) GoogleInitiate(c *fiber.Ctx) error {
	redirectURL, err := h.svc.GoogleInitiate(c.Context())
	if err != nil {
		return err
	}
	return c.Redirect(redirectURL, http.StatusTemporaryRedirect)
}

// GoogleCallback godoc
// @Summary      Google OAuth2 callback
// @Description  Exchanges the authorization code for JWT tokens. Redirects to GOOGLE_FRONTEND_URL with tokens if configured; returns JSON otherwise.
// @Tags         Auth
// @Param        code  query string true "Authorization code from Google"
// @Param        state query string true "CSRF state token"
// @Success      200   {object} response.Response{data=TokensResponse}
// @Failure      401   {object} response.Response
// @Router       /auth/google/callback [get]
func (h *Handler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")
	if code == "" || state == "" {
		return apperr.New("missing code or state parameter", http.StatusBadRequest)
	}
	tokens, err := h.svc.GoogleCallback(c.Context(), code, state)
	if err != nil {
		return err
	}
	if h.googleFrontendURL != "" {
		dest := fmt.Sprintf("%s?access_token=%s&refresh_token=%s",
			h.googleFrontendURL,
			url.QueryEscape(tokens.AccessToken),
			url.QueryEscape(tokens.RefreshToken))
		return c.Redirect(dest, http.StatusTemporaryRedirect)
	}
	return response.JSON(c, http.StatusOK, true, "login successful", tokens)
}

// VerifyOTP godoc
// @Summary      Verify email OTP
// @Description  Validates the 6-digit OTP sent to the user's email and marks the email as verified.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body VerifyOTPRequest true "OTP payload"
// @Success      200 {object} response.Response
// @Failure      422 {object} response.Response "Invalid or expired OTP"
// @Failure      429 {object} response.Response "Too many attempts"
// @Router       /auth/verify-otp [post]
func (h *Handler) VerifyOTP(c *fiber.Ctx) error {
	var req VerifyOTPRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	if err := h.svc.VerifyOTP(c.Context(), &req); err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "email verified", nil)
}

// ResendOTP godoc
// @Summary      Resend email verification OTP
// @Description  Sends a new OTP to the given email. Rate-limited to 3 requests per minute.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body ResendOTPRequest true "Email payload"
// @Success      200 {object} response.Response
// @Failure      429 {object} response.Response "Rate limit exceeded"
// @Router       /auth/resend-otp [post]
func (h *Handler) ResendOTP(c *fiber.Ctx) error {
	var req ResendOTPRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	if err := h.svc.ResendOTP(c.Context(), &req); err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "OTP sent if account exists", nil)
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Exchanges a valid refresh token for a new access + refresh token pair. Old token is revoked.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body RefreshRequest true "Refresh token payload"
// @Success      200 {object} response.Response{data=TokensResponse}
// @Failure      401 {object} response.Response "Invalid or expired refresh token"
// @Router       /auth/refresh [post]
func (h *Handler) Refresh(c *fiber.Ctx) error {
	var req RefreshRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	tokens, err := h.svc.Refresh(c.Context(), &req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "token refreshed", tokens)
}

// Logout godoc
// @Summary      Logout
// @Description  Revokes the provided refresh token.
// @Tags         Auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body LogoutRequest false "Optional: refresh token to revoke"
// @Success      200 {object} response.Response
// @Router       /auth/logout [post]
func (h *Handler) Logout(c *fiber.Ctx) error {
	var req LogoutRequest
	_ = c.BodyParser(&req)

	accessToken := strings.TrimPrefix(c.Get(fiber.HeaderAuthorization), "Bearer ")
	if err := h.svc.Logout(c.Context(), accessToken, &req); err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "logout successful", nil)
}

// ForgotPassword godoc
// @Summary      Request password reset
// @Description  Sends a password-reset email if the address exists. Always returns 200 to prevent enumeration.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body ForgotPasswordRequest true "Email payload"
// @Success      200 {object} response.Response
// @Router       /auth/forgot-password [post]
func (h *Handler) ForgotPassword(c *fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	if err := h.svc.ForgotPassword(c.Context(), &req); err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "if that address exists, a reset link was sent", nil)
}

// ResetPassword godoc
// @Summary      Reset password with token
// @Description  Sets a new password using the token sent via the forgot-password flow.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body ResetPasswordRequest true "Reset payload"
// @Success      200 {object} response.Response
// @Failure      401 {object} response.Response "Invalid token"
// @Router       /auth/reset-password [post]
func (h *Handler) ResetPassword(c *fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	if err := h.svc.ResetPassword(c.Context(), &req); err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "password updated", nil)
}

// Me godoc
// @Summary      Get current user
// @Description  Returns the authenticated user's profile.
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.Response{data=UserResponse}
// @Failure      401 {object} response.Response
// @Router       /auth/me [get]
func (h *Handler) Me(c *fiber.Ctx) error {
	uuid, _ := c.Locals(middleware.LocalsUserId).(string)
	u, err := h.svc.Me(c.Context(), uuid)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "profile retrieved", u)
}
