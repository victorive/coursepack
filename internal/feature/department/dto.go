package department

import "github.com/victorive/coursepack/internal/domain"

type DepartmentResponse struct {
	UUID       string              `json:"uuid"`
	University *UniversityResponse `json:"university"`
	Name       string              `json:"name"`
	Code       string              `json:"code"`
	Faculty    *string             `json:"faculty"`
}

// UniversityResponse is a local copy so department does not depend on the university feature.
type UniversityResponse struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	Acronym string `json:"acronym"`
}

func toUniversityResponse(u *domain.University) *UniversityResponse {
	return &UniversityResponse{
		UUID:    u.UUID,
		Name:    u.Name,
		Acronym: u.Acronym,
	}
}

func ToDepartmentResponse(d *domain.Department) *DepartmentResponse {
	resp := &DepartmentResponse{
		UUID:       d.UUID,
		University: toUniversityResponse(d.University),
		Name:       d.Name,
		Code:       d.Code,
		Faculty:    d.Faculty,
	}
	return resp
}
