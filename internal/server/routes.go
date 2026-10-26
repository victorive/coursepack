package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/coursepack/internal/app"
	"github.com/victorive/coursepack/internal/config"
	"github.com/victorive/coursepack/internal/pkg/response"
	v1 "github.com/victorive/coursepack/internal/server/routes/v1"
)

func RegisterRoutes(appFiber *fiber.App, cfg *config.Config, container *app.Container) {
	api := appFiber.Group("/api")
	v1.Register(api, container)

	appFiber.Use(func(c *fiber.Ctx) error {
	return response.Error(c, http.StatusNotFound, "route not found")
	})
}
