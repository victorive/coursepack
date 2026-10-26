package seeders

import (
	"fmt"

	"github.com/victorive/coursepack/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type rolePermissionsSeeder struct{}

func (rolePermissionsSeeder) Name() string { return "role_permissions" }

var rolePerms = map[domain.RoleName][]domain.PermissionName{
	domain.RoleSuperAdmin: {
		domain.PermMembersList,
		domain.PermMembersManage,
		domain.PermFilesUpload,
		domain.PermFilesApprove,
		domain.PermCoursesManage,
		domain.PermLecturesManage,
		domain.PermManageClasses,
		domain.PermManageDepartments,
		domain.PermManageSemesters,
	},
	domain.RoleClassRep: {
		domain.PermMembersList,
		domain.PermMembersManage,
		domain.PermFilesUpload,
		domain.PermFilesApprove,
		domain.PermCoursesManage,
		domain.PermLecturesManage,
		domain.PermManageSemesters,
	},
	domain.RoleCourseRep: {
		domain.PermMembersList,
		domain.PermFilesUpload,
		domain.PermFilesApprove,
		domain.PermCoursesManage,
		domain.PermLecturesManage,
	},
	domain.RoleStudent: {
		domain.PermMembersList,
		domain.PermFilesUpload,
		domain.PermManageClasses,
	},
}

func (rolePermissionsSeeder) Run(db *gorm.DB) error {
	var roles []domain.Role
	if err := db.Find(&roles).Error; err != nil {
		return fmt.Errorf("load roles: %w", err)
	}
	var perms []domain.Permission
	if err := db.Find(&perms).Error; err != nil {
		return fmt.Errorf("load permissions: %w", err)
	}

	roleByName := make(map[string]int64, len(roles))
	for _, r := range roles {
		roleByName[r.Name] = r.ID
	}
	permByName := make(map[string]int64, len(perms))
	for _, p := range perms {
		permByName[p.Name] = p.ID
	}

	var rows []domain.RolePermission
	for roleName, permNames := range rolePerms {
		roleID, ok := roleByName[string(roleName)]
		if !ok {
			return fmt.Errorf("role not found: %s", roleName)
		}
		for _, permName := range permNames {
			permID, ok := permByName[string(permName)]
			if !ok {
				return fmt.Errorf("permission not found: %s", permName)
			}
			rows = append(rows, domain.RolePermission{RoleID: roleID, PermissionID: permID})
		}
	}

	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&rows).Error
}
