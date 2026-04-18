package level

import (
	"time"

	"github.com/victorive/klaro-backend/internal/domain"
)

type CreateLevelRequest struct {
	DepartmentID int64 `json:"department_id" validate:"required,min=1"`
	Year         int16 `json:"year"          validate:"required,oneof=100 200 300 400 500 600 700 800"`
}

type LevelResponse struct {
	ID           int64     `json:"id"`
	DepartmentID int64     `json:"department_id"`
	Year         int16     `json:"year"`
	InviteCode   string    `json:"invite_code"`
	CreatedAt    time.Time `json:"created_at"`
}

func toLevelResponse(c *domain.Level) *LevelResponse {
	return &LevelResponse{
		ID:           c.ID,
		DepartmentID: c.DepartmentID,
		Year:         c.Year,
		InviteCode:   c.InviteCode,
		CreatedAt:    c.CreatedAt,
	}
}

type JoinRequest struct {
	InviteCode string `json:"invite_code" validate:"required,len=6"`
}

type ReviewJoinRequest struct {
	Action string `json:"action" validate:"required,oneof=approve reject"`
}

type JoinRequestResponse struct {
	ID          int64      `json:"id"`
	LevelID     int64      `json:"level_id"`
	UserID      int64      `json:"user_id"`
	Status      string     `json:"status"`
	RequestedAt time.Time  `json:"requested_at"`
	ResolvedAt  *time.Time `json:"resolved_at"`
}

func toJoinRequestResponse(jr *domain.JoinRequest) *JoinRequestResponse {
	return &JoinRequestResponse{
		ID:          jr.ID,
		LevelID:     jr.LevelID,
		UserID:      jr.UserID,
		Status:      string(jr.Status),
		RequestedAt: jr.RequestedAt,
		ResolvedAt:  jr.ResolvedAt,
	}
}

type MemberResponse struct {
	UserID   int64     `json:"user_id"`
	LevelID  int64     `json:"level_id"`
	JoinedAt time.Time `json:"joined_at"`
}

func toMemberResponse(m *domain.LevelMember) *MemberResponse {
	return &MemberResponse{
		UserID:   m.UserID,
		LevelID:  m.LevelID,
		JoinedAt: m.JoinedAt,
	}
}
