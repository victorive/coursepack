package department

import "github.com/victorive/klaro-backend/internal/domain"

type CreateDepartmentRequest struct {
	Name       string `json:"name"       validate:"required,min=2,max=200"`
	Code       string `json:"code"       validate:"required,min=1,max=20"`
	Faculty    string `json:"faculty"    validate:"required,min=2,max=200"`
	University string `json:"university" validate:"required,min=2,max=200"`
}

type UpdateDepartmentRequest struct {
	Name       string `json:"name"       validate:"required,min=2,max=200"`
	Code       string `json:"code"       validate:"required,min=1,max=20"`
	Faculty    string `json:"faculty"    validate:"required,min=2,max=200"`
	University string `json:"university" validate:"required,min=2,max=200"`
}

type DepartmentResponse struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	Faculty    string `json:"faculty"`
	University string `json:"university"`
}

func toDepartmentResponse(d *domain.Department) *DepartmentResponse {
	return &DepartmentResponse{
		ID:         d.ID,
		Name:       d.Name,
		Code:       d.Code,
		Faculty:    d.Faculty,
		University: d.University,
	}
}
