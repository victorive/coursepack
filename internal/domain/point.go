package domain

import "time"

type UserPoint struct {
	UserID      int64     `gorm:"column:user_id;primaryKey"`
	Balance     int       `gorm:"column:balance;not null;default:0"`
	TotalEarned int       `gorm:"column:total_earned;not null;default:0"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (UserPoint) TableName() string { return "user_points" }

type PointTransaction struct {
	ID           int64       `gorm:"column:id;primaryKey;autoIncrement"`
	UserID       int64       `gorm:"column:user_id;not null"`
	Delta        int         `gorm:"column:delta;not null"`
	BalanceAfter int         `gorm:"column:balance_after;not null"`
	SourceType   PointSource `gorm:"column:source_type;not null"`
	SourceID     *int64      `gorm:"column:source_id"`
	Note         *string     `gorm:"column:note"`
	CreatedBy    *int64      `gorm:"column:created_by"`
	CreatedAt    time.Time   `gorm:"column:created_at;autoCreateTime"`
}
