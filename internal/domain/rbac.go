package domain

import (
	"slices"
	"time"
)

type Role struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string    `gorm:"column:name;uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

type Permission struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string    `gorm:"column:name;uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

type RolePermission struct {
	RoleID       int64 `gorm:"column:role_id;primaryKey"`
	PermissionID int64 `gorm:"column:permission_id;primaryKey"`
}

type UserRole struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    int64     `gorm:"column:user_id;not null"`
	RoleID    int64     `gorm:"column:role_id;not null"`
	ScopeType string    `gorm:"column:scope_type;not null;default:global"`
	ScopeID   *int64    `gorm:"column:scope_id"`
	GrantedBy *int64    `gorm:"column:granted_by"`
	GrantedAt time.Time `gorm:"column:granted_at;autoCreateTime"`
}

type UserProfile struct {
	User
	Roles       []Role
	Permissions []Permission
}

func (p *UserProfile) HasRole(name RoleName) bool {
	return slices.ContainsFunc(p.Roles, func(r Role) bool {
		return RoleName(r.Name) == name
	})
}

func (p *UserProfile) HasPermission(name PermissionName) bool {
	return slices.ContainsFunc(p.Permissions, func(p Permission) bool {
		return PermissionName(p.Name) == name
	})
}

func (p *UserProfile) IsSuperAdmin() bool { return p.HasRole(RoleSuperAdmin) }
func (p *UserProfile) IsClassRep() bool   { return p.HasRole(RoleClassRep) }
func (p *UserProfile) IsCourseRep() bool  { return p.HasRole(RoleCourseRep) }
func (p *UserProfile) IsStudent() bool    { return p.HasRole(RoleStudent) }
