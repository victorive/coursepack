package course

import (
	"time"

	"github.com/victorive/klaro-backend/internal/domain"
)

type CreateCourseRequest struct {
	CohortID   int64  `json:"cohort_id"   validate:"required,min=1"`
	SemesterID int64  `json:"semester_id" validate:"required,min=1"`
	Code       string `json:"code"        validate:"required,min=1,max=20"`
	Title      string `json:"title"       validate:"required,min=2,max=200"`
	Units      int16  `json:"units"       validate:"required,min=1,max=10"`
}

type UpdateCourseRequest struct {
	Code  string `json:"code"  validate:"required,min=1,max=20"`
	Title string `json:"title" validate:"required,min=2,max=200"`
	Units int16  `json:"units" validate:"required,min=1,max=10"`
}

type CourseResponse struct {
	ID         int64     `json:"id"`
	CohortID   int64     `json:"cohort_id"`
	SemesterID int64     `json:"semester_id"`
	Code       string    `json:"code"`
	Title      string    `json:"title"`
	Units      int16     `json:"units"`
	CreatedAt  time.Time `json:"created_at"`
}

func toCourseResponse(c *domain.Course) *CourseResponse {
	return &CourseResponse{
		ID:         c.ID,
		CohortID:   c.CohortID,
		SemesterID: c.SemesterID,
		Code:       c.Code,
		Title:      c.Title,
		Units:      c.Units,
		CreatedAt:  c.CreatedAt,
	}
}

type CreateLectureRequest struct {
	CourseID     int64   `json:"course_id"     validate:"required,min=1"`
	Topic        string  `json:"topic"         validate:"required,min=2,max=300"`
	Description  *string `json:"description"`
	LecturerName *string `json:"lecturer_name"`
	WeekNumber   *int16  `json:"week_number"   validate:"omitempty,min=1,max=52"`
	Kind         string  `json:"kind"          validate:"required,oneof=lecture test exam lab"`
}

type UpdateLectureRequest struct {
	Topic        string  `json:"topic"         validate:"required,min=2,max=300"`
	Description  *string `json:"description"`
	LecturerName *string `json:"lecturer_name"`
	WeekNumber   *int16  `json:"week_number"   validate:"omitempty,min=1,max=52"`
	Kind         string  `json:"kind"          validate:"required,oneof=lecture test exam lab"`
}

type LectureResponse struct {
	ID           int64     `json:"id"`
	CourseID     int64     `json:"course_id"`
	Topic        string    `json:"topic"`
	Description  *string   `json:"description"`
	LecturerName *string   `json:"lecturer_name"`
	WeekNumber   *int16    `json:"week_number"`
	Kind         string    `json:"kind"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

func toLectureResponse(l *domain.Lecture) *LectureResponse {
	return &LectureResponse{
		ID:           l.ID,
		CourseID:     l.CourseID,
		Topic:        l.Topic,
		Description:  l.Description,
		LecturerName: l.LecturerName,
		WeekNumber:   l.WeekNumber,
		Kind:         l.Kind,
		Status:       l.Status,
		CreatedAt:    l.CreatedAt,
	}
}
