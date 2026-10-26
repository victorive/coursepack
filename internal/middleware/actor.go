package middleware

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/coursepack/internal/domain"
	"github.com/victorive/coursepack/internal/permissions"
	"github.com/victorive/coursepack/internal/pkg/apperr"
	"github.com/victorive/coursepack/internal/pkg/response"
)

// UserReader is the subset of the user repository needed to build an Actor.
type UserReader interface {
	FindByUUID(ctx context.Context, uuid string) (*domain.User, error)
}

// Actor loads the authenticated user's identity, roles, and permissions once
// per request and stores them on the request context.
func Actor(userRepo UserReader, permSvc permissions.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uuid, err := UserUUID(c)
		if err != nil {
			return err
		}

		u, err := userRepo.FindByUUID(c.UserContext(), uuid)
		if err != nil {
			return response.Error(c, http.StatusInternalServerError, "an unexpected error occurred")
		}
		if u == nil {
			return apperr.New("unauthenticated", http.StatusUnauthorized)
		}

		roleNames, err := permSvc.GetUserRoles(c.UserContext(), uuid)
		if err != nil {
			return response.Error(c, http.StatusInternalServerError, "an unexpected error occurred")
		}

		permNames, err := permSvc.GetUserPermissions(c.UserContext(), uuid)
		if err != nil {
			return response.Error(c, http.StatusInternalServerError, "an unexpected error occurred")
		}

		actor := permissions.Actor{
			UserID:       u.ID,
			UUID:         uuid,
			UniversityID: u.UniversityID,
			Roles:        toRoleNames(roleNames),
			Permissions:  toPermissionNames(permNames),
		}

		c.SetUserContext(permissions.WithActor(c.UserContext(), actor))
		return c.Next()
	}
}

func toRoleNames(names []string) []domain.RoleName {
	out := make([]domain.RoleName, len(names))
	for i, n := range names {
		out[i] = domain.RoleName(n)
	}
	return out
}

func toPermissionNames(names []string) []domain.PermissionName {
	out := make([]domain.PermissionName, len(names))
	for i, n := range names {
		out[i] = domain.PermissionName(n)
	}
	return out
}
