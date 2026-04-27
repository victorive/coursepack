package rbac

import (
	"context"
	"log/slog"

	"github.com/victorive/klaro-backend/internal/domain"
)

type Service interface {
	HasGlobalPermission(ctx context.Context, userUUID string, perm domain.PermissionName) (bool, error)
	HasClassPermission(ctx context.Context, userUUID string, classID int64, perm domain.PermissionName) (bool, error)
	HasClassPermissionByUUID(ctx context.Context, userUUID, classUUID string, perm domain.PermissionName) (bool, error)
	AssignGlobalRole(ctx context.Context, userID int64, role domain.RoleName) error
	AssignClassRole(ctx context.Context, userID, classID int64, role domain.RoleName) error
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

func (s *service) HasClassPermission(ctx context.Context, userUUID string, classID int64, perm domain.PermissionName) (bool, error) {
	return s.repo.HasClassPermission(ctx, userUUID, classID, perm)
}

func (s *service) HasClassPermissionByUUID(ctx context.Context, userUUID, classUUID string, perm domain.PermissionName) (bool, error) {
	return s.repo.HasClassPermissionByUUID(ctx, userUUID, classUUID, perm)
}

func (s *service) AssignGlobalRole(ctx context.Context, userID int64, role domain.RoleName) error {
	return s.repo.AssignGlobalRole(ctx, userID, role)
}

func (s *service) AssignClassRole(ctx context.Context, userID, classID int64, role domain.RoleName) error {
	return s.repo.AssignClassRole(ctx, userID, classID, role)
}

func (s *service) GetUserRoles(ctx context.Context, userUUID string) ([]string, error) {
	return s.repo.GetUserRoles(ctx, userUUID)
}

func (s *service) GetUserPermissions(ctx context.Context, userUUID string) ([]string, error) {
	return s.repo.GetUserPermissions(ctx, userUUID)
}
