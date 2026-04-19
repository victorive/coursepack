package seeders

import (
	"github.com/victorive/klaro-backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type rolesSeeder struct{}

func (rolesSeeder) Name() string { return "roles" }

func (rolesSeeder) Run(db *gorm.DB) error {
	roles := []domain.Role{
		{Name: string(domain.RoleSuperAdmin)},
		{Name: string(domain.RoleStudent)},
		{Name: string(domain.RoleClassRep)},
		{Name: string(domain.RoleCourseRep)},
	}
	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&roles).Error
}
