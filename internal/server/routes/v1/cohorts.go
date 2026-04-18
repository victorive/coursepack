package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/app"
	"github.com/victorive/klaro-backend/internal/feature/cohort"
	"github.com/victorive/klaro-backend/internal/middleware"
)

func cohortRoutes(r fiber.Router, c *app.Container) {
	h := cohort.NewHandler(c.CohortService)
	auth := []fiber.Handler{middleware.Auth(c.AuthService), middleware.OTPVerified(c.AuthService)}

	cohorts := r.Group("/cohorts", auth...)
	cohorts.Post("", h.Create)
	cohorts.Get("", h.List)
	cohorts.Get("/:id", h.Get)
	cohorts.Delete("/:id", h.Delete)

	cohorts.Post("/join", h.RequestJoin)
	cohorts.Get("/:id/join-requests", h.ListPendingJoinRequests)
	cohorts.Patch("/:id/join-requests/:requestId", h.ReviewJoinRequest)

	cohorts.Get("/:id/members", h.ListMembers)
	cohorts.Delete("/:id/members/:userId", h.RemoveMember)
}
