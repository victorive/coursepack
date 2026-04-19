package auth

type SignUpRequest struct {
	FullName string `json:"full_name" validate:"required,min=2,max=100"`
	Email    string `json:"email"     validate:"required,email"`
	Password string `json:"password"  validate:"required,min=8"`
}

type SignInRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp"   validate:"required,len=6"`
}

type ResendOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ExchangeGoogleCodeRequest struct {
	Code string `json:"code" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token"    validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type TokensResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresAt    int           `json:"expires_at"`
	User         *UserResponse `json:"user"`
}

type UserResponse struct {
	ID            int64    `json:"id"`
	UUID          string   `json:"uuid"`
	Email         string   `json:"email"`
	FullName      string   `json:"full_name"`
	Phone         *string  `json:"phone"`
	AvatarURL     *string  `json:"avatar_url"`
	Provider      *string  `json:"provider"`
	EmailVerified bool     `json:"email_verified"`
	UniversityID  *int64   `json:"university_id"`
	Roles         []string `json:"roles"`
	Permissions   []string `json:"permissions"`
}
