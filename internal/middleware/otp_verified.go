package middleware

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/coursepack/internal/pkg/response"
)

type OTPVerifiedChecker interface {
	CheckEmailVerified(ctx context.Context, uuid string) (emailVerified bool, isOAuthUser bool, err error)
}

func OTPVerified(svc OTPVerifiedChecker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals(LocalsUserId).(string)
		if !ok || userID == "" {
			return response.JSON(c, http.StatusUnauthorized, false, "unauthenticated", nil)
		}

		verified, isOAuth, err := svc.CheckEmailVerified(c.UserContext(), userID)
		if err != nil {
			return response.JSON(c, http.StatusInternalServerError, false, "internal error", nil)
		}

		if !isOAuth && !verified {
			return response.JSON(c, http.StatusForbidden, false, "email not verified", nil)
		}

		return c.Next()
	}
}
