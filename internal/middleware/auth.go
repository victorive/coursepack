package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/coursepack/internal/pkg/apperr"
	"github.com/victorive/coursepack/internal/pkg/response"
)

const LocalsUserID = "userId"

// UserUUID returns the authenticated user's UUID from the Fiber locals.
func UserUUID(c *fiber.Ctx) (string, error) {
	uuid, ok := c.Locals(LocalsUserID).(string)
	if !ok || uuid == "" {
		return "", apperr.New("unauthenticated", http.StatusUnauthorized)
	}
	return uuid, nil
}

type TokenValidator interface {
	ValidateToken(ctx context.Context, token string) (sub string, err error)
}

func Auth(v TokenValidator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get(fiber.HeaderAuthorization)
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			return response.Error(c, http.StatusUnauthorized, "unauthenticated")
		}

		tokenString := parts[1]
		sub, err := v.ValidateToken(c.UserContext(), tokenString)
		if err != nil {
			return response.Error(c, http.StatusUnauthorized, "invalid or expired token")
		}

		c.Locals(LocalsUserID, sub)
		return c.Next()
	}
}
