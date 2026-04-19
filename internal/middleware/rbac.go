package middleware

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/domain"
	"github.com/victorive/klaro-backend/internal/feature/rbac"
	"github.com/victorive/klaro-backend/internal/pkg/response"
)

func RequireGlobalPermission(repo rbac.Service, perm domain.PermissionName) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uuid, ok := UserUUID(c)
		if !ok {
			return response.JSON(c, http.StatusUnauthorized, false, "unauthenticated", nil)
		}
		ok, err := repo.HasGlobalPermission(c.UserContext(), uuid, perm)
		if err != nil {
			return response.JSON(c, http.StatusInternalServerError, false, "an unexpected error occurred", nil)
		}
		if !ok {
			return response.JSON(c, http.StatusForbidden, false, "unauthorized", nil)
		}
		return c.Next()
	}
}

func RequireLevelPermission(repo rbac.Service, perm domain.PermissionName, paramName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uuid, ok := UserUUID(c)
		if !ok {
			return response.JSON(c, http.StatusUnauthorized, false, "unauthenticated", nil)
		}
		levelID, err := strconv.ParseInt(c.Params(paramName), 10, 64)
		if err != nil {
			return response.JSON(c, http.StatusBadRequest, false, "invalid level id", nil)
		}
		ok, err = repo.HasLevelPermission(c.UserContext(), uuid, levelID, perm)
		if err != nil {
			return response.JSON(c, http.StatusInternalServerError, false, "an unexpected error occurred", nil)
		}
		if !ok {
			return response.JSON(c, http.StatusForbidden, false, "unauthorized", nil)
		}
		return c.Next()
	}
}
