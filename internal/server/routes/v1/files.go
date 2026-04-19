package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/app"
	"github.com/victorive/klaro-backend/internal/domain"
	"github.com/victorive/klaro-backend/internal/feature/file"
	"github.com/victorive/klaro-backend/internal/middleware"
)

func fileRoutes(r fiber.Router, c *app.Container) {
	h := file.NewHandler(c.FileService)

	auth := []fiber.Handler{middleware.Auth(c.AuthService), middleware.OTPVerified(c.AuthService)}

	files := r.Group("/files", auth...)
	files.Post("", middleware.RequireGlobalPermission(c.RBAC, domain.PermFilesUpload), h.Upload)
	files.Get("/:id", h.GetFile)
	files.Get("/:id/download", h.GetDownloadURL)
	files.Patch("/:id/review", middleware.RequireGlobalPermission(c.RBAC, domain.PermFilesApprove), h.Review)
	files.Patch("/:id/archive", middleware.RequireGlobalPermission(c.RBAC, domain.PermFilesApprove), h.Archive)
}

