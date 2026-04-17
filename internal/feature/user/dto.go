package user

type UpdateProfileRequest struct {
	Phone string `json:"phone" validate:"required"`
}
