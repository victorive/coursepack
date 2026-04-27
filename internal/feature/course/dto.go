package course

import (
	"time"

	"github.com/victorive/klaro-backend/internal/domain"
)

type CourseRegistryResponse struct {
	ID           int64     `json:"id"`
	UniversityID int64     `json:"university_id"`
	DepartmentID int64     `json:"department_id"`
	Code         string    `json:"code"`
	Title        string    `json:"title"`
	CreditUnits  int16     `json:"credit_units"`
	CreatedAt    time.Time `json:"created_at"`
}

func toCourseRegistryResponse(c *domain.Course) *CourseRegistryResponse {
	return &CourseRegistryResponse{
		ID:           c.ID,
		UniversityID: c.UniversityID,
		DepartmentID: c.DepartmentID,
		Code:         c.Code,
		Title:        c.Title,
		CreditUnits:  c.CreditUnits,
		CreatedAt:    c.CreatedAt,
	}
}

type AddClassCourseRequest struct {
	SemesterID int64 `json:"semester_id" validate:"required,min=1"`
	CourseID   int64 `json:"course_id"   validate:"required,min=1"`
}

type ClassCourseResponse struct {
	ID         int64                   `json:"id"`
	ClassID    int64                   `json:"class_id"`
	SemesterID int64                   `json:"semester_id"`
	CourseID   int64                   `json:"course_id"`
	Course     *CourseRegistryResponse `json:"course,omitempty"`
	CreatedAt  time.Time               `json:"created_at"`
}

func toClassCourseResponse(cc *domain.ClassCourse) *ClassCourseResponse {
	resp := &ClassCourseResponse{
		ID:         cc.ID,
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
	ID            int64     `json:"id"`
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
		ID:            l.ID,
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
