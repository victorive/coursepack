package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/app"
)

func Register(r fiber.Router, c *app.Container) {
	v1 := r.Group("/v1")
	authRoutes(v1, c)
	userRoutes(v1, c)
	departmentRoutes(v1, c)
	courseRoutes(v1, c)
	fileRoutes(v1, c)
	miscRoutes(v1, c)
}
