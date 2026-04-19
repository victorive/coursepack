package academic

import (
	"time"

	"github.com/victorive/klaro-backend/internal/domain"
)

type CreateSemesterRequest struct {
	ClassID   int64      `json:"class_id"   validate:"required,min=1"`
	Session   string     `json:"session"    validate:"required,min=9,max=9"`
	Kind      string     `json:"kind"       validate:"required,oneof=first second summer"`
	IsCurrent bool       `json:"is_current"`
	StartsAt  *time.Time `json:"starts_at"`
	EndsAt    *time.Time `json:"ends_at"`
}

type UpdateSemesterRequest struct {
	Kind      string     `json:"kind"       validate:"required,oneof=first second summer"`
	IsCurrent bool       `json:"is_current"`
	StartsAt  *time.Time `json:"starts_at"`
	EndsAt    *time.Time `json:"ends_at"`
}

type SemesterResponse struct {
	ID        int64      `json:"id"`
	ClassID   int64      `json:"class_id"`
	Session   string     `json:"session"`
	Kind      string     `json:"kind"`
	IsCurrent bool       `json:"is_current"`
	StartsAt  *time.Time `json:"starts_at"`
	EndsAt    *time.Time `json:"ends_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func toSemesterResponse(s *domain.Semester) *SemesterResponse {
	return &SemesterResponse{
		ID:        s.ID,
		ClassID:   s.ClassID,
		Session:   s.Session,
		Kind:      s.Kind,
		IsCurrent: s.IsCurrent,
		StartsAt:  s.StartsAt,
		EndsAt:    s.EndsAt,
		CreatedAt: s.CreatedAt,
	}
}
