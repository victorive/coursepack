package domain

import "time"

type File struct {
	ID         int64      `gorm:"column:id;primaryKey;autoIncrement"`
	UploadedBy int64      `gorm:"column:uploaded_by;not null"`
	CourseID   int64      `gorm:"column:course_id;not null"`
	LectureID  *int64     `gorm:"column:lecture_id"`
	Category   string     `gorm:"column:category;not null"`
	Title      string     `gorm:"column:title;not null"`
	StorageKey string     `gorm:"column:storage_key;uniqueIndex;not null"`
	SizeBytes  int64      `gorm:"column:size_bytes;not null"`
	MimeType   string     `gorm:"column:mime_type;not null"`
	YearLabel  *string    `gorm:"column:year_label"`
	Status     string     `gorm:"column:status;not null;default:pending"`
	ReviewedBy *int64     `gorm:"column:reviewed_by"`
	ReviewNote *string    `gorm:"column:review_note"`
	ReviewedAt *time.Time `gorm:"column:reviewed_at"`
	CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}
