package rbac

import (
	"context"
	"fmt"

	"github.com/victorive/klaro-backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repository abstracts RBAC role checks and assignments.
type Repository interface {
	HasGlobalRole(ctx context.Context, userID int64, role domain.RoleName) (bool, error)
	HasCohortRole(ctx context.Context, userID, cohortID int64, role domain.RoleName) (bool, error)
	AssignGlobalRole(ctx context.Context, userID, grantedByID int64, role domain.RoleName) error
	AssignCohortRole(ctx context.Context, userID, cohortID, grantedByID int64, role domain.RoleName) error
}

type gormRepository struct {
	db *gorm.DB
}

// NewRepository returns a GORM-backed RBAC repository.
func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) HasGlobalRole(ctx context.Context, userID int64, role domain.RoleName) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("user_roles ur").
		Joins("JOIN roles r ON r.id = ur.role_id").
		Where("ur.user_id = ? AND ur.scope_type = 'global' AND r.name = ?", userID, string(role)).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("rbac.HasGlobalRole: %w", err)
	}
	return count > 0, nil
}

func (r *gormRepository) HasCohortRole(ctx context.Context, userID, cohortID int64, role domain.RoleName) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("user_roles ur").
		Joins("JOIN roles r ON r.id = ur.role_id").
		Where("ur.user_id = ? AND ur.scope_type = 'cohort' AND ur.scope_id = ? AND r.name = ?", userID, cohortID, string(role)).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("rbac.HasCohortRole: %w", err)
	}
	return count > 0, nil
}

func (r *gormRepository) AssignGlobalRole(ctx context.Context, userID, grantedByID int64, role domain.RoleName) error {
	roleID, err := r.roleIDByName(ctx, role)
	if err != nil {
		return err
	}
	ur := domain.UserRole{
		UserID:    userID,
		RoleID:    roleID,
		ScopeType: "global",
		ScopeID:   nil,
		GrantedBy: &grantedByID,
	}
	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&ur).Error; err != nil {
		return fmt.Errorf("rbac.AssignGlobalRole: %w", err)
	}
	return nil
}

func (r *gormRepository) AssignCohortRole(ctx context.Context, userID, cohortID, grantedByID int64, role domain.RoleName) error {
	roleID, err := r.roleIDByName(ctx, role)
	if err != nil {
		return err
	}
	ur := domain.UserRole{
		UserID:    userID,
		RoleID:    roleID,
		ScopeType: "cohort",
		ScopeID:   &cohortID,
		GrantedBy: &grantedByID,
	}
	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&ur).Error; err != nil {
		return fmt.Errorf("rbac.AssignCohortRole: %w", err)
	}
	return nil
}

func (r *gormRepository) roleIDByName(ctx context.Context, name domain.RoleName) (int64, error) {
	var role domain.Role
	if err := r.db.WithContext(ctx).Where("name = ?", string(name)).First(&role).Error; err != nil {
		return 0, fmt.Errorf("rbac.roleIDByName(%s): %w", name, err)
	}
	return role.ID, nil
}
