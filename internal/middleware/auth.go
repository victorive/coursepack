package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/pkg/response"
)

const LocalsUserId = "userId"

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
