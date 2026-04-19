package class

import (
	"time"

	"github.com/victorive/klaro-backend/internal/domain"
)

type CreateClassRequest struct {
	DepartmentID     int64  `json:"department_id"    validate:"required,min=1"`
	AdmissionSession string `json:"admission_session" validate:"required,min=9,max=9"`
	CurrentLevelID   int64  `json:"current_level_id" validate:"required,min=1"`
	Name             string `json:"name"             validate:"required,min=1,max=100"`
}

type UpdateClassRequest struct {
	Name           string `json:"name"             validate:"omitempty,min=1,max=100"`
	CurrentLevelID int64  `json:"current_level_id" validate:"omitempty,min=1"`
}

type ClassResponse struct {
	ID               int64     `json:"id"`
	DepartmentID     int64     `json:"department_id"`
	AdmissionSession string    `json:"admission_session"`
	CurrentLevelID   int64     `json:"current_level_id"`
	Name             string    `json:"name"`
	InviteCode       string    `json:"invite_code"`
	CreatedAt        time.Time `json:"created_at"`
}

func toClassResponse(c *domain.Class) *ClassResponse {
	return &ClassResponse{
		ID:               c.ID,
		DepartmentID:     c.DepartmentID,
		AdmissionSession: c.AdmissionSession,
		CurrentLevelID:   c.CurrentLevelID,
		Name:             c.Name,
		InviteCode:       c.InviteCode,
		CreatedAt:        c.CreatedAt,
	}
}

type JoinRequest struct {
	InviteCode string `json:"invite_code" validate:"required,len=8"`
}

type ReviewJoinRequest struct {
	Action string `json:"action" validate:"required,oneof=approve reject"`
}

type JoinRequestResponse struct {
	ID          int64      `json:"id"`
	ClassID     int64      `json:"class_id"`
	UserID      int64      `json:"user_id"`
	Status      string     `json:"status"`
	RequestedAt time.Time  `json:"requested_at"`
	ResolvedAt  *time.Time `json:"resolved_at"`
}

func toJoinRequestResponse(jr *domain.JoinRequest) *JoinRequestResponse {
	return &JoinRequestResponse{
		ID:          jr.ID,
		ClassID:     jr.ClassID,
		UserID:      jr.UserID,
		Status:      string(jr.Status),
		RequestedAt: jr.RequestedAt,
		ResolvedAt:  jr.ResolvedAt,
	}
}

type MemberResponse struct {
	User     *SimpleUserResponse `json:"user"`
	ClassID  int64               `json:"class_id"`
	JoinedAt time.Time           `json:"joined_at"`
}

type SimpleUserResponse struct {
	ID        int64   `json:"id"`
	UUID      string  `json:"uuid"`
	Email     string  `json:"email"`
	FullName  string  `json:"full_name"`
	Phone     *string `json:"phone"`
	AvatarURL *string `json:"avatar_url"`
}

func toMemberResponse(m *domain.ClassMember) *MemberResponse {
	var user *SimpleUserResponse
	if m.User != nil {
		user = &SimpleUserResponse{
			ID:        m.User.ID,
			UUID:      m.User.UUID,
			Email:     m.User.Email,
			FullName:  m.User.FullName,
			Phone:     m.User.Phone,
			AvatarURL: m.User.AvatarURL,
		}
	}
	return &MemberResponse{
		User:     user,
		ClassID:  m.ClassID,
		JoinedAt: m.JoinedAt,
	}
}

type LevelResponse struct {
	ID    int64  `json:"id"`
	Value int16  `json:"value"`
	Label string `json:"label"`
}

func toLevelResponse(l *domain.Level) *LevelResponse {
	return &LevelResponse{ID: l.ID, Value: l.Value, Label: l.Label}
}
