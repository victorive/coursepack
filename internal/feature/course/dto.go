package course

import (
	"time"

	"github.com/victorive/klaro-backend/internal/domain"
	"github.com/victorive/klaro-backend/internal/pkg/pagination"
)

type CourseRegistryResponse struct {
	UUID        string `json:"uuid"`
	Code        string `json:"code"`
	Title       string `json:"title"`
	CreditUnits int16  `json:"credit_units"`
}

type CourseRegistryPaginatedResponse struct {
	Items []CourseRegistryResponse `json:"items"`
	Meta  pagination.CursorMeta    `json:"meta"`
}

func toCourseRegistryResponse(c *domain.Course) *CourseRegistryResponse {
	return &CourseRegistryResponse{
		UUID:        c.UUID,
		Code:        c.Code,
		Title:       c.Title,
		CreditUnits: c.CreditUnits,
	}
}

type AddClassCourseRequest struct {
	SemesterID int64  `json:"semester_id" validate:"required,min=1"`
	CourseUUID string `json:"course_uuid"  validate:"required"`
}

type ClassCourseResponse struct {
	UUID       string                  `json:"uuid"`
	ClassID    int64                   `json:"class_id"`
	SemesterID int64                   `json:"semester_id"`
	CourseID   int64                   `json:"course_id"`
	Course     *CourseRegistryResponse `json:"course,omitempty"`
	CreatedAt  time.Time               `json:"created_at"`
}

func toClassCourseResponse(cc *domain.ClassCourse) *ClassCourseResponse {
	resp := &ClassCourseResponse{
		UUID:       cc.UUID,
		ClassID:    cc.ClassID,
		SemesterID: cc.SemesterID,
		CourseID:   cc.CourseID,
		CreatedAt:  cc.CreatedAt,
	}
	if cc.Course != nil {
		resp.Course = toCourseRegistryResponse(cc.Course)
	}
	return resp
}

type CreateLectureRequest struct {
	Topic        string  `json:"topic"           validate:"required,min=2,max=300"`
	Description  *string `json:"description"`
	LecturerName *string `json:"lecturer_name"`
	WeekNumber   *int16  `json:"week_number"     validate:"omitempty,min=1,max=52"`
	Kind         string  `json:"kind"            validate:"required,oneof=lecture test exam lab"`
}

type UpdateLectureRequest struct {
	Topic        string  `json:"topic"         validate:"omitempty,min=2,max=300"`
	Description  *string `json:"description"`
	LecturerName *string `json:"lecturer_name"`
	WeekNumber   *int16  `json:"week_number"   validate:"omitempty,min=1,max=52"`
	Kind         string  `json:"kind"          validate:"omitempty,oneof=lecture test exam lab"`
}

type LectureResponse struct {
	UUID          string    `json:"uuid"`
	ClassCourseID int64     `json:"class_course_id"`
	Topic         string    `json:"topic"`
	Description   *string   `json:"description"`
	LecturerName  *string   `json:"lecturer_name"`
	WeekNumber    *int16    `json:"week_number"`
	Kind          string    `json:"kind"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

func toLectureResponse(l *domain.Lecture) *LectureResponse {
	return &LectureResponse{
		UUID:          l.UUID,
		ClassCourseID: l.ClassCourseID,
		Topic:         l.Topic,
		Description:   l.Description,
		LecturerName:  l.LecturerName,
		WeekNumber:    l.WeekNumber,
		Kind:          l.Kind,
		Status:        l.Status,
		CreatedAt:     l.CreatedAt,
	}
}
