package seeders

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/victorive/klaro-backend/internal/domain"
	"github.com/victorive/klaro-backend/internal/pkg/hash"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type superAdminSeeder struct{}

func (superAdminSeeder) Name() string { return "super_admin" }

func (superAdminSeeder) Run(db *gorm.DB) error {
	ph, err := hash.HashPassword("Password@123!")
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	u := domain.User{
		UUID:          uuid.New().String(),
		Email:         "admin@klaro.app",
		FullName:      "Super Admin",
		Password:      &ph,
		EmailVerified: true,
	}

	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"full_name", "email_verified"}),
	}).Create(&u).Error; err != nil {
		return fmt.Errorf("upsert super admin user: %w", err)
	}

	if err := db.Where("email = ?", u.Email).First(&u).Error; err != nil {
		return fmt.Errorf("fetch super admin user: %w", err)
	}

	var role domain.Role
	if err := db.Where("name = ?", string(domain.RoleSuperAdmin)).First(&role).Error; err != nil {
		return fmt.Errorf("fetch super_admin role (run roles seeder first): %w", err)
	}

	ur := domain.UserRole{
		UserID:    u.ID,
		RoleID:    role.ID,
		ScopeType: "global",
	}

	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&ur).Error; err != nil {
		return fmt.Errorf("assign super_admin role: %w", err)
	}

	return nil
}
