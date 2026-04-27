package university

import "github.com/victorive/klaro-backend/internal/domain"

type UniversityResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Acronym string `json:"acronym"`
}

func ToUniversityResponse(u *domain.University) *UniversityResponse {
	if u == nil {
		return nil
	}
	return &UniversityResponse{
		ID:      u.ID,
		Name:    u.Name,
		Acronym: u.Acronym,
	}
}

func ToUniversityResponses(universities []*domain.University) []*UniversityResponse {
	responses := make([]*UniversityResponse, len(universities))
	for i, u := range universities {
		responses[i] = ToUniversityResponse(u)
	}
	return responses
}
