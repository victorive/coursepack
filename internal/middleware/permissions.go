package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/coursepack/internal/domain"
	"github.com/victorive/coursepack/internal/permissions"
	"github.com/victorive/coursepack/internal/pkg/response"
)

func RequireGlobalPermission(perm domain.PermissionName) fiber.Handler {
	return func(c *fiber.Ctx) error {
		actor, err := permissions.ActorFrom(c.UserContext())
		if err != nil {
			return err
		}
		if !actor.HasPermission(perm) {
			return response.Error(c, http.StatusForbidden, "insufficient permissions")
		}
		return c.Next()
	}
}

func RequireClassPermission(repo permissions.Service, perm domain.PermissionName, paramName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		actor, err := permissions.ActorFrom(c.UserContext())
		if err != nil {
			return err
		}
		classUUID := c.Params(paramName)
		if classUUID == "" {
			return response.Error(c, http.StatusBadRequest, "invalid class id")
		}
		allowed, err := repo.HasClassPermissionByUUID(c.UserContext(), actor.UUID, classUUID, perm)
		if err != nil {
			return response.Error(c, http.StatusInternalServerError, "an unexpected error occurred")
		}
		if !allowed {
			return response.Error(c, http.StatusForbidden, "insufficient permissions")
		}
		return c.Next()
	}
}
