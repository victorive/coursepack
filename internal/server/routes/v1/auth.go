package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/coursepack/internal/app"
	"github.com/victorive/coursepack/internal/feature/auth"
	"github.com/victorive/coursepack/internal/middleware"
)

func authRoutes(r fiber.Router, c *app.Container) {
	h := auth.NewHandler(c.AuthService, c.Config.Google.FrontendURL)

	auth := r.Group("/auth")
	auth.Post("/signup", h.SignUp)
	auth.Post("/signin", h.SignIn)
	auth.Get("/google", h.GoogleInitiate)
	auth.Get("/google/callback", h.GoogleCallback)
	auth.Post("/google/exchange", h.ExchangeGoogleCode)
	auth.Post("/verify-otp", h.VerifyOTP)
	auth.Post("/resend-otp", h.ResendOTP)
	auth.Post("/refresh", h.Refresh)
	auth.Post("/forgot-password", h.ForgotPassword)
	auth.Post("/reset-password", h.ResetPassword)

	protected := auth.Group("", middleware.Auth(c.AuthService), middleware.OTPVerified(c.AuthService))
	protected.Post("/logout", h.Logout)
}
