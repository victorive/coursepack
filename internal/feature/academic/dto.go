package academic

import (
	"time"

	"github.com/victorive/klaro-backend/internal/domain"
)

type CreateSessionRequest struct {
	Label     string `json:"label"      validate:"required,min=2,max=50"`
	IsCurrent bool   `json:"is_current"`
}

type UpdateSessionRequest struct {
	Label     string `json:"label"      validate:"required,min=2,max=50"`
	IsCurrent bool   `json:"is_current"`
}

type SessionResponse struct {
	ID        int64     `json:"id"`
	Label     string    `json:"label"`
	IsCurrent bool      `json:"is_current"`
	CreatedAt time.Time `json:"created_at"`
}

func toSessionResponse(s *domain.AcademicSession) *SessionResponse {
	return &SessionResponse{
		ID:        s.ID,
		Label:     s.Label,
		IsCurrent: s.IsCurrent,
		CreatedAt: s.CreatedAt,
	}
}

type CreateSemesterRequest struct {
	LevelID   int64      `json:"level_id"   validate:"required,min=1"`
	SessionID int64      `json:"session_id" validate:"required,min=1"`
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
	LevelID   int64      `json:"level_id"`
	SessionID int64      `json:"session_id"`
	Kind      string     `json:"kind"`
	IsCurrent bool       `json:"is_current"`
	StartsAt  *time.Time `json:"starts_at"`
	EndsAt    *time.Time `json:"ends_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func toSemesterResponse(s *domain.Semester) *SemesterResponse {
	return &SemesterResponse{
		ID:        s.ID,
		LevelID:   s.LevelID,
		SessionID: s.SessionID,
		Kind:      s.Kind,
		IsCurrent: s.IsCurrent,
		StartsAt:  s.StartsAt,
		EndsAt:    s.EndsAt,
		CreatedAt: s.CreatedAt,
	}
}
