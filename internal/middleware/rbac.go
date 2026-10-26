package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/coursepack/internal/domain"
	"github.com/victorive/coursepack/internal/feature/rbac"
	"github.com/victorive/coursepack/internal/pkg/response"
)

func RequireGlobalPermission(repo rbac.Service, perm domain.PermissionName) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uuid, err := UserUUID(c)
		if err != nil {
			return err
		}
		allowed, err := repo.HasGlobalPermission(c.UserContext(), uuid, perm)
		if err != nil {
			return response.JSON(c, http.StatusInternalServerError, false, "an unexpected error occurred", nil)
		}
		if !allowed {
			return response.JSON(c, http.StatusForbidden, false, "insufficient permissions", nil)
		}
		return c.Next()
	}
}

func RequireClassPermission(repo rbac.Service, perm domain.PermissionName, paramName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userUUID, err := UserUUID(c)
		if err != nil {
			return err
		}
		classUUID := c.Params(paramName)
		if classUUID == "" {
			return response.JSON(c, http.StatusBadRequest, false, "invalid class id", nil)
		}
		allowed, err := repo.HasClassPermissionByUUID(c.UserContext(), userUUID, classUUID, perm)
		if err != nil {
			return response.JSON(c, http.StatusInternalServerError, false, "an unexpected error occurred", nil)
		}
		if !allowed {
			return response.JSON(c, http.StatusForbidden, false, "insufficient permissions", nil)
		}
		return c.Next()
	}
}
