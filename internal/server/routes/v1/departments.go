package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/app"
	"github.com/victorive/klaro-backend/internal/domain"
	"github.com/victorive/klaro-backend/internal/feature/department"
	"github.com/victorive/klaro-backend/internal/middleware"
)

func departmentRoutes(r fiber.Router, c *app.Container) {
	h := department.NewHandler(c.DepartmentService)
	auth := []fiber.Handler{middleware.Auth(c.AuthService), middleware.OTPVerified(c.AuthService)}

	protected := r.Group("/departments", auth...)
	protected.Post("", middleware.RequireGlobalPermission(c.RBAC, domain.PermManageDepartments), h.Create)
	protected.Get("", h.List)
	protected.Get("/:id", h.Get)
	protected.Patch("/:id", middleware.RequireGlobalPermission(c.RBAC, domain.PermManageDepartments), h.Update)
	protected.Delete("/:id", middleware.RequireGlobalPermission(c.RBAC, domain.PermManageDepartments), h.Delete)
}
