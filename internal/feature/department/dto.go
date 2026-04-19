package department

import "github.com/victorive/klaro-backend/internal/domain"

type CreateDepartmentRequest struct {
	Name         string  `json:"name"          validate:"required,min=2,max=200"`
	Code         string  `json:"code"          validate:"required,min=1,max=20"`
	Faculty      *string `json:"faculty"       validate:"omitempty,min=2,max=200"`
	UniversityID int64   `json:"university_id" validate:"required,min=1"`
}

type UpdateDepartmentRequest struct {
	Name    string  `json:"name"    validate:"required,min=2,max=200"`
	Code    string  `json:"code"    validate:"required,min=1,max=20"`
	Faculty *string `json:"faculty" validate:"omitempty,min=2,max=200"`
}

type DepartmentResponse struct {
	ID           int64   `json:"id"`
	UniversityID int64   `json:"university_id"`
	Name         string  `json:"name"`
	Code         string  `json:"code"`
	Faculty      *string `json:"faculty"`
}

func toDepartmentResponse(d *domain.Department) *DepartmentResponse {
	return &DepartmentResponse{
		ID:           d.ID,
		UniversityID: d.UniversityID,
		Name:         d.Name,
		Code:         d.Code,
		Faculty:      d.Faculty,
	}
}
