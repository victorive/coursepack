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

	classes := r.Group("/classes", auth...)
	classes.Post("/:classId/semesters", middleware.RequireClassPermission(c.RBAC, domain.PermManageSemesters, "classId"), h.CreateSemester)
	classes.Get("/:classId/semesters", h.ListSemestersByClass)
	classes.Get("/:classId/semesters/:semId", h.GetSemester)
	classes.Patch("/:classId/semesters/:semId", middleware.RequireClassPermission(c.RBAC, domain.PermManageSemesters, "classId"), h.UpdateSemester)
	classes.Delete("/:classId/semesters/:semId", middleware.RequireClassPermission(c.RBAC, domain.PermManageSemesters, "classId"), h.DeleteSemester)
}
