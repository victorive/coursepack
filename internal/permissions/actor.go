package permissions

import (
	"context"
	"net/http"

	"github.com/victorive/coursepack/internal/domain"
	"github.com/victorive/coursepack/internal/pkg/apperr"
)

// Actor carries the authenticated user's identity and permissions for the
// lifetime of a single request. It is created once by middleware and passed
// down through context or explicitly as a function argument.
type Actor struct {
	UserID       int64
	UUID         string
	UniversityID *int64
	Roles        []domain.RoleName
	Permissions  []domain.PermissionName
}

// HasPermission reports whether the actor holds a global permission.
func (a Actor) HasPermission(perm domain.PermissionName) bool {
	for _, p := range a.Permissions {
		if p == perm {
			return true
		}
	}
	return false
}

// HasRole reports whether the actor holds a global role.
func (a Actor) HasRole(role domain.RoleName) bool {
	for _, r := range a.Roles {
		if r == role {
			return true
		}
	}
	return false
}

type actorKey struct{}

// WithActor stores an Actor on the context.
func WithActor(ctx context.Context, actor Actor) context.Context {
	return context.WithValue(ctx, actorKey{}, actor)
}

// ActorFrom returns the Actor stored on the context, or an unauthorized error.
func ActorFrom(ctx context.Context) (Actor, error) {
	actor, ok := ctx.Value(actorKey{}).(Actor)
	if !ok {
		return Actor{}, apperr.New("unauthenticated", http.StatusUnauthorized)
	}
	return actor, nil
}
