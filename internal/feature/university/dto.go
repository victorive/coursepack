package university

import (
	"time"

	"github.com/victorive/klaro-backend/internal/domain"
)

type UniversityResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Acronym   string    `json:"acronym"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToUniversityResponse(u *domain.University) *UniversityResponse {
	if u == nil {
		return nil
	}
	return &UniversityResponse{
		ID:        u.ID,
		Name:      u.Name,
		Acronym:   u.Acronym,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func ToUniversityResponses(universities []*domain.University) []*UniversityResponse {
	responses := make([]*UniversityResponse, len(universities))
	for i, u := range universities {
		responses[i] = ToUniversityResponse(u)
	}
	return responses
}
