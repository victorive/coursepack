package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/app"
	"github.com/victorive/klaro-backend/internal/feature/department"
	"github.com/victorive/klaro-backend/internal/middleware"
)

func departmentRoutes(r fiber.Router, c *app.Container) {
	h := department.NewHandler(c.DepartmentService)

	protected := r.Group("/departments", middleware.Auth(c.AuthService), middleware.OTPVerified(c.AuthService))
	protected.Post("", h.Create)
	protected.Get("", h.List)
	protected.Get("/:id", h.Get)
	protected.Patch("/:id", h.Update)
	protected.Delete("/:id", h.Delete)
}
