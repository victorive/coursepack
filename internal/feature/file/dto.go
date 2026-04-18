package file

import (
	"time"

	"github.com/victorive/klaro-backend/internal/domain"
)

type UploadFileRequest struct {
	CourseID  int64   `form:"course_id" validate:"required"`
	LectureID *int64  `form:"lecture_id"`
	Category  string  `form:"category" validate:"required"`
	Title     string  `form:"title" validate:"required"`
	YearLabel *string `form:"year_label"`
}

type ReviewFileRequest struct {
	Action     string  `json:"action" validate:"required,oneof=approve reject"`
	ReviewNote *string `json:"review_note"`
}

type FileResponse struct {
	ID         int64      `json:"id"`
	UploadedBy int64      `json:"uploaded_by"`
	CourseID   int64      `json:"course_id"`
	LectureID  *int64     `json:"lecture_id,omitempty"`
	Category   string     `json:"category"`
	Title      string     `json:"title"`
	StorageKey string     `json:"storage_key"`
	SizeBytes  int64      `json:"size_bytes"`
	MimeType   string     `json:"mime_type"`
	YearLabel  *string    `json:"year_label,omitempty"`
	Status     string     `json:"status"`
	ReviewedBy *int64     `json:"reviewed_by,omitempty"`
	ReviewNote *string    `json:"review_note,omitempty"`
	ReviewedAt *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

func toFileResponse(f *domain.File) *FileResponse {
	return &FileResponse{
		ID:         f.ID,
		UploadedBy: f.UploadedBy,
		CourseID:   f.CourseID,
		LectureID:  f.LectureID,
		Category:   f.Category,
		Title:      f.Title,
		StorageKey: f.StorageKey,
		SizeBytes:  f.SizeBytes,
		MimeType:   f.MimeType,
		YearLabel:  f.YearLabel,
		Status:     f.Status,
		ReviewedBy: f.ReviewedBy,
		ReviewNote: f.ReviewNote,
		ReviewedAt: f.ReviewedAt,
		CreatedAt:  f.CreatedAt,
	}
}
