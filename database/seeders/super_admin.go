package seeders

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/victorive/coursepack/internal/domain"
	"github.com/victorive/coursepack/internal/pkg/hash"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type superAdminSeeder struct{}

func (superAdminSeeder) Name() string { return "super_admin" }

func (superAdminSeeder) Run(db *gorm.DB) error {
	email := os.Getenv("SUPER_ADMIN_EMAIL")
	password := os.Getenv("SUPER_ADMIN_PASSWORD")
	if email == "" || password == "" {
		return fmt.Errorf("SUPER_ADMIN_EMAIL and SUPER_ADMIN_PASSWORD must be set to seed the super admin")
	}
	if len(password) < 12 {
		return fmt.Errorf("SUPER_ADMIN_PASSWORD must be at least 12 characters")
	}

	ph, err := hash.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	u := domain.User{
		UUID:          uuid.New().String(),
		Email:         email,
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
