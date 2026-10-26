package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/coursepack/internal/app"
	"github.com/victorive/coursepack/internal/domain"
	"github.com/victorive/coursepack/internal/feature/file"
	"github.com/victorive/coursepack/internal/middleware"
)

func fileRoutes(r fiber.Router, c *app.Container) {
	h := file.NewHandler(c.FileService)

	auth := []fiber.Handler{middleware.Auth(c.AuthService), middleware.OTPVerified(c.AuthService)}

	files := r.Group("/files", auth...)
	files.Post("", middleware.RequireGlobalPermission(domain.PermFilesUpload), h.Upload)
	files.Get("/:id", h.GetFile)
	files.Get("/:id/download", h.GetDownloadURL)
	files.Patch("/:id/review", h.Review)
	files.Patch("/:id/archive", h.Archive)
}
