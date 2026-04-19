package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/app"
	"github.com/victorive/klaro-backend/internal/domain"
	"github.com/victorive/klaro-backend/internal/feature/academic"
	"github.com/victorive/klaro-backend/internal/middleware"
)

func academicRoutes(r fiber.Router, c *app.Container) {
	h := academic.NewHandler(c.AcademicService)
	auth := []fiber.Handler{middleware.Auth(c.AuthService), middleware.OTPVerified(c.AuthService)}

	semesters := r.Group("/semesters", auth...)
	semesters.Post("", middleware.RequireGlobalPermission(c.RBAC, domain.PermManageSemesters), h.CreateSemester)
	semesters.Get("/:id", h.GetSemester)
	semesters.Patch("/:id", middleware.RequireGlobalPermission(c.RBAC, domain.PermManageSemesters), h.UpdateSemester)
	semesters.Delete("/:id", middleware.RequireGlobalPermission(c.RBAC, domain.PermManageSemesters), h.DeleteSemester)

	r.Group("/classes", auth...).Get("/:classId/semesters", h.ListSemestersByClass)
}
