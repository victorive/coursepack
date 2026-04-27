package academic

import (
	"github.com/victorive/klaro-backend/internal/domain"
)

type CreateSemesterRequest struct {
	ClassUUID string `json:"class_uuid"   validate:"required"`
	Session   string `json:"session"      validate:"required,min=9,max=9"`
	Kind      string `json:"kind"         validate:"required,oneof=first second"`
	IsCurrent bool   `json:"is_current"`
	StartsAt  string `json:"starts_at"    validate:"omitempty,datetime=2006-01-02,ltefield=EndsAt"`
	EndsAt    string `json:"ends_at"      validate:"omitempty,datetime=2006-01-02,gtefield=StartsAt"`
}

type UpdateSemesterRequest struct {
	Kind      string `json:"kind"       validate:"required,oneof=first second"`
	IsCurrent bool   `json:"is_current"`
	StartsAt  string `json:"starts_at"  validate:"omitempty,datetime=2006-01-02,ltefield=EndsAt"`
	EndsAt    string `json:"ends_at"    validate:"omitempty,datetime=2006-01-02,gtefield=StartsAt"`
}

type SemesterResponse struct {
	ID        int64   `json:"id"`
	ClassID   int64   `json:"class_id"`
	Session   string  `json:"session"`
	Kind      string  `json:"kind"`
	IsCurrent bool    `json:"is_current"`
	StartsAt  *string `json:"starts_at"`
	EndsAt    *string `json:"ends_at"`
	CreatedAt string  `json:"created_at"`
}

func toSemesterResponse(s *domain.Semester) *SemesterResponse {
	var startsAt *string
	if s.StartsAt != nil {
		v := s.StartsAt.Format("2006-01-02")
		startsAt = &v
	}

	var endsAt *string
	if s.EndsAt != nil {
		v := s.EndsAt.Format("2006-01-02")
		endsAt = &v
	}
	
	return &SemesterResponse{
		ID:        s.ID,
		ClassID:   s.ClassID,
		Session:   s.Session,
		Kind:      s.Kind,
		IsCurrent: s.IsCurrent,
		StartsAt:  startsAt,
		EndsAt:    endsAt,
		CreatedAt: s.CreatedAt.Format("2006-01-02"),
	}
}
