package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/coursepack/internal/app"
	"github.com/victorive/coursepack/internal/feature/auth"
	"github.com/victorive/coursepack/internal/feature/user"
	"github.com/victorive/coursepack/internal/middleware"
)

func userRoutes(r fiber.Router, c *app.Container) {
	authHandler := auth.NewHandler(c.AuthService, c.Config.Google.FrontendURL)
	userHandler := user.NewHandler(c.UserService)

	protected := r.Group("/users", middleware.Auth(c.AuthService), middleware.OTPVerified(c.AuthService))
	protected.Get("/me", authHandler.Me)
	protected.Patch("/me", userHandler.UpdateProfile)
}
