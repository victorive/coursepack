package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/coursepack/internal/app"
	"github.com/victorive/coursepack/internal/feature/department"
	"github.com/victorive/coursepack/internal/middleware"
)

func departmentRoutes(r fiber.Router, c *app.Container) {
	h := department.NewHandler(c.DepartmentService)
	auth := []fiber.Handler{middleware.Auth(c.AuthService), middleware.OTPVerified(c.AuthService), middleware.Actor(c.UserRepository, c.PermissionsService)}

	protected := r.Group("/departments", auth...)
	protected.Get("", h.List)
	protected.Get("/:id", h.Get)
}
