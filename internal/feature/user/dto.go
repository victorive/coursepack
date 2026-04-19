package user

type UpdateProfileRequest struct {
	Phone        *string `json:"phone"         validate:"omitempty"`
	UniversityID *int64  `json:"university_id" validate:"omitempty,min=1"`
}
