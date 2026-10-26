package rbac

import (
	"context"
	"fmt"

	"github.com/victorive/coursepack/internal/domain"
	"github.com/victorive/coursepack/internal/pkg/dbtx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	HasGlobalRole(ctx context.Context, userID int64, role domain.RoleName) (bool, error)
	HasClassRole(ctx context.Context, userID, classID int64, role domain.RoleName) (bool, error)
	AssignGlobalRole(ctx context.Context, userID int64, role domain.RoleName) error
	AssignClassRole(ctx context.Context, userID, classID int64, role domain.RoleName) error
	HasGlobalPermission(ctx context.Context, userUUID string, perm domain.PermissionName) (bool, error)
	HasClassPermission(ctx context.Context, userUUID string, classID int64, perm domain.PermissionName) (bool, error)
	HasClassPermissionByUUID(ctx context.Context, userUUID, classUUID string, perm domain.PermissionName) (bool, error)
	GetUserRoles(ctx context.Context, userUUID string) ([]string, error)
	GetUserPermissions(ctx context.Context, userUUID string) ([]string, error)
}

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) HasGlobalRole(ctx context.Context, userID int64, role domain.RoleName) (bool, error) {
	var count int64

	err := dbtx.Conn(ctx, r.db).
		Table("user_roles ur").
		Joins("JOIN roles r ON r.id = ur.role_id").
		Where("ur.user_id = ? AND ur.scope_type = 'global' AND r.name = ?", userID, string(role)).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("rbac.HasGlobalRole: %w", err)
	}

	return count > 0, nil
}

func (r *gormRepository) HasClassRole(ctx context.Context, userID, classID int64, role domain.RoleName) (bool, error) {
	var count int64

	err := dbtx.Conn(ctx, r.db).
		Table("user_roles ur").
		Joins("JOIN roles r ON r.id = ur.role_id").
		Where("ur.user_id = ? AND ur.scope_type = 'class' AND ur.scope_id = ? AND r.name = ?", userID, classID, string(role)).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("rbac.HasClassRole: %w", err)
	}

	return count > 0, nil
}

func (r *gormRepository) AssignGlobalRole(ctx context.Context, userID int64, role domain.RoleName) error {
	roleID, err := r.roleIDByName(ctx, role)
	if err != nil {
		return err
	}

	ur := domain.UserRole{
		UserID:    userID,
		RoleID:    roleID,
		ScopeType: "global",
		ScopeID:   nil,
	}

	if err := dbtx.Conn(ctx, r.db).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&ur).Error; err != nil {
		return fmt.Errorf("rbac.AssignGlobalRole: %w", err)
	}

	return nil
}

func (r *gormRepository) AssignClassRole(ctx context.Context, userID, classID int64, role domain.RoleName) error {
	roleID, err := r.roleIDByName(ctx, role)
	if err != nil {
		return err
	}

	ur := domain.UserRole{
		UserID:    userID,
		RoleID:    roleID,
		ScopeType: "class",
		ScopeID:   &classID,
	}

	if err := dbtx.Conn(ctx, r.db).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&ur).Error; err != nil {
		return fmt.Errorf("rbac.AssignClassRole: %w", err)
	}

	return nil
}

func (r *gormRepository) HasGlobalPermission(ctx context.Context, userUUID string, perm domain.PermissionName) (bool, error) {
	var count int64

	err := dbtx.Conn(ctx, r.db).
		Table("user_roles ur").
		Joins("JOIN users u ON u.id = ur.user_id").
		Joins("JOIN roles ro ON ro.id = ur.role_id").
		Joins("JOIN role_permissions rp ON rp.role_id = ro.id").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Where("u.uuid = ? AND ur.scope_type = 'global' AND p.name = ?", userUUID, string(perm)).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("rbac.HasGlobalPermission: %w", err)
	}

	return count > 0, nil
}

func (r *gormRepository) HasClassPermission(ctx context.Context, userUUID string, classID int64, perm domain.PermissionName) (bool, error) {
	var count int64

	err := dbtx.Conn(ctx, r.db).
		Table("user_roles ur").
		Joins("JOIN users u ON u.id = ur.user_id").
		Joins("JOIN roles ro ON ro.id = ur.role_id").
		Joins("JOIN role_permissions rp ON rp.role_id = ro.id").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Where(
			"u.uuid = ? AND p.name = ? AND (ur.scope_type = 'global' OR (ur.scope_type = 'class' AND ur.scope_id = ?))",
			userUUID, string(perm), classID,
		).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("rbac.HasClassPermission: %w", err)
	}

	return count > 0, nil
}

func (r *gormRepository) HasClassPermissionByUUID(ctx context.Context, userUUID, classUUID string, perm domain.PermissionName) (bool, error) {
	var count int64

	err := dbtx.Conn(ctx, r.db).
		Table("user_roles ur").
		Joins("JOIN users u ON u.id = ur.user_id").
		Joins("JOIN roles ro ON ro.id = ur.role_id").
		Joins("JOIN role_permissions rp ON rp.role_id = ro.id").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Joins("LEFT JOIN classes c ON c.id = ur.scope_id").
		Where(
			"u.uuid = ? AND p.name = ? AND (ur.scope_type = 'global' OR (ur.scope_type = 'class' AND c.uuid = ?))",
			userUUID, string(perm), classUUID,
		).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("rbac.HasClassPermissionByUUID: %w", err)
	}

	return count > 0, nil
}

func (r *gormRepository) GetUserRoles(ctx context.Context, userUUID string) ([]string, error) {
	var roles []string

	err := dbtx.Conn(ctx, r.db).
		Table("user_roles ur").
		Joins("JOIN users u ON u.id = ur.user_id").
		Joins("JOIN roles ro ON ro.id = ur.role_id").
		Where("u.uuid = ?", userUUID).
		Distinct().
		Pluck("ro.name", &roles).Error
	if err != nil {
		return nil, fmt.Errorf("rbac.GetUserRoles: %w", err)
	}

	return roles, nil
}

func (r *gormRepository) GetUserPermissions(ctx context.Context, userUUID string) ([]string, error) {
	var permissions []string

	err := dbtx.Conn(ctx, r.db).
		Table("user_roles ur").
		Joins("JOIN users u ON u.id = ur.user_id").
		Joins("JOIN roles ro ON ro.id = ur.role_id").
		Joins("JOIN role_permissions rp ON rp.role_id = ro.id").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Where("u.uuid = ?", userUUID).
		Distinct().
		Pluck("p.name", &permissions).Error
	if err != nil {
		return nil, fmt.Errorf("rbac.GetUserPermissions: %w", err)
	}

	return permissions, nil
}

func (r *gormRepository) roleIDByName(ctx context.Context, name domain.RoleName) (int64, error) {
	var role domain.Role

	if err := dbtx.Conn(ctx, r.db).Where("name = ?", string(name)).First(&role).Error; err != nil {
		return 0, fmt.Errorf("rbac.roleIDByName(%s): %w", name, err)
	}

	return role.ID, nil
}
