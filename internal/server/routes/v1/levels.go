package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/app"
	"github.com/victorive/klaro-backend/internal/feature/level"
	"github.com/victorive/klaro-backend/internal/middleware"
)

func levelRoutes(r fiber.Router, c *app.Container) {
	h := level.NewHandler(c.LevelService)
	auth := []fiber.Handler{middleware.Auth(c.AuthService), middleware.OTPVerified(c.AuthService)}

	levels := r.Group("/levels", auth...)
	levels.Post("", h.Create)
	levels.Get("", h.List)
	levels.Get("/:id", h.Get)
	levels.Delete("/:id", h.Delete)

	levels.Post("/join", h.RequestJoin)
	levels.Get("/:id/join-requests", h.ListPendingJoinRequests)
	levels.Patch("/:id/join-requests/:requestId", h.ReviewJoinRequest)

	levels.Get("/:id/members", h.ListMembers)
	levels.Delete("/:id/members/:userId", h.RemoveMember)
}
