package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/app"
)

func courseRoutes(r fiber.Router, _ *app.Container) {
	_ = r.Group("/lectures")
}
