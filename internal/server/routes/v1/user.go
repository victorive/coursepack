package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/app"
	"github.com/victorive/klaro-backend/internal/feature/user"
	"github.com/victorive/klaro-backend/internal/middleware"
)

func userRoutes(r fiber.Router, c *app.Container) {
	h := user.NewHandler(c.UserService)

	user := r.Group("/user")
	protected := user.Group("", middleware.Auth(c.JWT))
	protected.Patch("", h.UpdateProfile)
}
