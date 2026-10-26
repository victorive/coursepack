package file

import (
	"time"

	"github.com/victorive/coursepack/internal/domain"
)

type UploadFileRequest struct {
	ClassCourseUUID string  `form:"class_course_uuid" validate:"required"`
	LectureUUID     *string `form:"lecture_uuid"`
	Category        string  `form:"category" validate:"required,oneof=lecture_note past_question material assignment"`
	Title           string  `form:"title" validate:"required"`
	YearLabel       *string `form:"year_label"`
}

type ReviewFileRequest struct {
	Action     string  `json:"action" validate:"required,oneof=approve reject"`
	ReviewNote *string `json:"review_note"`
}

type FileResponse struct {
	UUID          string     `json:"uuid"`
	UploadedBy    int64      `json:"uploaded_by"`
	ClassCourseID int64      `json:"class_course_id"`
	LectureID     *int64     `json:"lecture_id,omitempty"`
	Category      string     `json:"category"`
	Title         string     `json:"title"`
	SizeBytes     int64      `json:"size_bytes"`
	MimeType      string     `json:"mime_type"`
	YearLabel     *string    `json:"year_label,omitempty"`
	Status        string     `json:"status"`
	ReviewedBy    *int64     `json:"reviewed_by,omitempty"`
	ReviewNote    *string    `json:"review_note,omitempty"`
	ReviewedAt    *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

func toFileResponse(f *domain.File) *FileResponse {
	return &FileResponse{
		UUID:          f.UUID,
		UploadedBy:    f.UploadedBy,
		ClassCourseID: f.ClassCourseID,
		LectureID:     f.LectureID,
		Category:      f.Category,
		Title:         f.Title,
		SizeBytes:     f.SizeBytes,
		MimeType:      f.MimeType,
		YearLabel:     f.YearLabel,
		Status:        f.Status,
		ReviewedBy:    f.ReviewedBy,
		ReviewNote:    f.ReviewNote,
		ReviewedAt:    f.ReviewedAt,
		CreatedAt:     f.CreatedAt,
	}
}
