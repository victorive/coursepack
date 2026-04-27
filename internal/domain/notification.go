package domain

import "time"

type Notification struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID       int64     `gorm:"column:user_id;not null"`
	Title        string    `gorm:"column:title;not null"`
	Body         string    `gorm:"column:body;not null"`
	IsRead       bool      `gorm:"column:is_read;not null;default:false"`
	Kind         string    `gorm:"column:kind;not null"`
	ResourceType *string   `gorm:"column:resource_type"`
	ResourceID   *int64    `gorm:"column:resource_id"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

type Suggestion struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    int64     `gorm:"column:user_id;not null"`
	Body      string    `gorm:"column:body;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

type Report struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement"`
	ReportedBy int64     `gorm:"column:reported_by;not null"`
	TargetType string    `gorm:"column:target_type;not null"`
	TargetID   int64     `gorm:"column:target_id;not null"`
	Reason     string    `gorm:"column:reason;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}
