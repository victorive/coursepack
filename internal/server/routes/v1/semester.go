package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/coursepack/internal/app"
	"github.com/victorive/coursepack/internal/domain"
	"github.com/victorive/coursepack/internal/feature/semester"
	"github.com/victorive/coursepack/internal/middleware"
)

func semesterRoutes(r fiber.Router, c *app.Container) {
	h := semester.NewHandler(c.AcademicService)
	auth := []fiber.Handler{middleware.Auth(c.AuthService), middleware.OTPVerified(c.AuthService)}

	classes := r.Group("/classes", auth...)
	classes.Post("/:classId/semesters", middleware.RequireClassPermission(c.PermissionsService, domain.PermManageSemesters, "classId"), h.CreateSemester)
	classes.Get("/:classId/semesters", h.ListSemestersByClass)
	classes.Get("/:classId/semesters/:semId", h.GetSemester)
	classes.Patch("/:classId/semesters/:semId", middleware.RequireClassPermission(c.PermissionsService, domain.PermManageSemesters, "classId"), h.UpdateSemester)
	classes.Delete("/:classId/semesters/:semId", middleware.RequireClassPermission(c.PermissionsService, domain.PermManageSemesters, "classId"), h.DeleteSemester)
}
