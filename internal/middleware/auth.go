package middleware

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/victorive/klaro-backend/internal/pkg/jwt"
	"github.com/victorive/klaro-backend/internal/pkg/response"
)

const LocalsUserUUID = "user_uuid"

func Auth(mgr *jwt.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get(fiber.HeaderAuthorization)
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			return response.JSON(c, http.StatusUnauthorized, false, "missing or invalid authorization header", nil)
		}

		claims, err := mgr.Verify(parts[1])
		if err != nil {
			return response.JSON(c, http.StatusUnauthorized, false, "invalid token", nil)
		}

		c.Locals(LocalsUserUUID, claims.UUID)
		return c.Next()
	}
}
