package department

import (
	"github.com/victorive/klaro-backend/internal/domain"
	"github.com/victorive/klaro-backend/internal/feature/university"
)

type DepartmentResponse struct {
	ID         int64                          `json:"id"`
	University *university.UniversityResponse `json:"university"`
	Name       string                         `json:"name"`
	Code       string                         `json:"code"`
	Faculty    *string                        `json:"faculty"`
}

func toDepartmentResponse(d *domain.Department) *DepartmentResponse {
	resp := &DepartmentResponse{
		ID:         d.ID,
		University: university.ToUniversityResponse(d.University),
		Name:       d.Name,
		Code:       d.Code,
		Faculty:    d.Faculty,
	}
	return resp
}
