package rbac

import (
	"context"
	"log/slog"

	"github.com/victorive/klaro-backend/internal/domain"
)

type Service interface {
	HasGlobalPermission(ctx context.Context, userUUID string, perm domain.PermissionName) (bool, error)
	HasLevelPermission(ctx context.Context, userUUID string, levelID int64, perm domain.PermissionName) (bool, error)
	AssignGlobalRole(ctx context.Context, userID int64, role domain.RoleName) error
	AssignLevelRole(ctx context.Context, userID, levelID int64, role domain.RoleName) error
	GetUserRoles(ctx context.Context, userUUID string) ([]string, error)
	GetUserPermissions(ctx context.Context, userUUID string) ([]string, error)
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{repo: repo, logger: logger}
}

func (s *service) HasGlobalPermission(ctx context.Context, userUUID string, perm domain.PermissionName) (bool, error) {
	return s.repo.HasGlobalPermission(ctx, userUUID, perm)
}

func (s *service) HasLevelPermission(ctx context.Context, userUUID string, levelID int64, perm domain.PermissionName) (bool, error) {
	return s.repo.HasLevelPermission(ctx, userUUID, levelID, perm)
}

func (s *service) AssignGlobalRole(ctx context.Context, userID int64, role domain.RoleName) error {
	return s.repo.AssignGlobalRole(ctx, userID, role)
}

func (s *service) AssignLevelRole(ctx context.Context, userID, levelID int64, role domain.RoleName) error {
	return s.repo.AssignLevelRole(ctx, userID, levelID, role)
}

func (s *service) GetUserRoles(ctx context.Context, userUUID string) ([]string, error) {
	return s.repo.GetUserRoles(ctx, userUUID)
}

func (s *service) GetUserPermissions(ctx context.Context, userUUID string) ([]string, error) {
	return s.repo.GetUserPermissions(ctx, userUUID)
}
