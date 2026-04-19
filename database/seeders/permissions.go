package seeders

import (
	"github.com/victorive/klaro-backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type permissionsSeeder struct{}

func (permissionsSeeder) Name() string { return "permissions" }

func (permissionsSeeder) Run(db *gorm.DB) error {
	perms := []domain.Permission{
		{Name: string(domain.PermMembersList)},
		{Name: string(domain.PermMembersManage)},
		{Name: string(domain.PermFilesUpload)},
		{Name: string(domain.PermFilesApprove)},
		{Name: string(domain.PermCoursesManage)},
		{Name: string(domain.PermLecturesManage)},
		{Name: string(domain.PermManageClasses)},
		{Name: string(domain.PermManageDepartments)},
		{Name: string(domain.PermManageSemesters)},
	}
	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&perms).Error
}
