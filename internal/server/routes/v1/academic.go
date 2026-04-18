package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/app"
	"github.com/victorive/klaro-backend/internal/feature/academic"
	"github.com/victorive/klaro-backend/internal/middleware"
)

func academicRoutes(r fiber.Router, c *app.Container) {
	h := academic.NewHandler(c.AcademicService)
	auth := []fiber.Handler{middleware.Auth(c.AuthService), middleware.OTPVerified(c.AuthService)}

	sessions := r.Group("/sessions", auth...)
	sessions.Post("", h.CreateSession)
	sessions.Get("", h.ListSessions)
	sessions.Get("/:id", h.GetSession)
	sessions.Patch("/:id", h.UpdateSession)

	semesters := r.Group("/semesters", auth...)
	semesters.Post("", h.CreateSemester)
	semesters.Get("/:id", h.GetSemester)
	semesters.Patch("/:id", h.UpdateSemester)
	semesters.Delete("/:id", h.DeleteSemester)

	r.Group("/cohorts", auth...).Get("/:cohortId/semesters", h.ListSemestersByCohort)
}
