package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/coursepack/internal/app"
	"github.com/victorive/coursepack/internal/feature/university"
)

func universityRoutes(r fiber.Router, c *app.Container) {
	h := university.NewHandler(c.UniversityService)

	universities := r.Group("/universities")
	universities.Get("", h.List)
	universities.Get("/:id", h.GetByUUID)
}
