package domain

import "time"

type User struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UUID         string    `gorm:"column:uuid;uniqueIndex;not null"`
	Email        string    `gorm:"column:email;uniqueIndex;not null"`
	Password     *string   `gorm:"column:password"`
	FullName     string    `gorm:"column:full_name;not null"`
	Phone        *string   `gorm:"column:phone"`
	AvatarURL    *string   `gorm:"column:avatar_url"`
	Provider     *string   `gorm:"column:provider"`
	EmailVerified bool      `gorm:"column:email_verified;not null;default:false"`
	UniversityID *int64    `gorm:"column:university_id"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
