package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/pkg/response"
)

const LocalsUserId = "userId"

// UserUUID extracts the authenticated user's UUID from Fiber context locals.
func UserUUID(c *fiber.Ctx) (string, bool) {
	uuid, ok := c.Locals(LocalsUserId).(string)
	if !ok || uuid == "" {
		return "", false
	}
	return uuid, true
}

type TokenValidator interface {
	ValidateToken(ctx context.Context, token string) (sub string, err error)
}

func Auth(v TokenValidator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get(fiber.HeaderAuthorization)
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			return response.JSON(c, http.StatusUnauthorized, false, "unauthenticated", nil)
		}

		tokenString := parts[1]
		sub, err := v.ValidateToken(c.Context(), tokenString)
		if err != nil {
			return response.JSON(c, http.StatusUnauthorized, false, err.Error(), nil)
		}

		c.Locals(LocalsUserId, sub) //uuid
		return c.Next()
	}
}
