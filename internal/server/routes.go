package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/app"
	"github.com/victorive/klaro-backend/internal/config"
	"github.com/victorive/klaro-backend/internal/pkg/response"
	v1 "github.com/victorive/klaro-backend/internal/server/routes/v1"
)

func RegisterRoutes(appFiber *fiber.App, cfg *config.Config, container *app.Container) {
	api := appFiber.Group("/api")
	v1.Register(api, container)

	appFiber.Use(func(c *fiber.Ctx) error {
	return response.JSON(c, http.StatusNotFound, false, "route not found", nil)
	})
}
