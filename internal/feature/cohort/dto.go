package cohort

import (
	"time"

	"github.com/victorive/klaro-backend/internal/domain"
)

type CreateCohortRequest struct {
	DepartmentID int64  `json:"department_id" validate:"required,min=1"`
	Level        int16  `json:"level"         validate:"required,oneof=100 200 300 400 500 600 700 800"`
}

type CohortResponse struct {
	ID           int64     `json:"id"`
	DepartmentID int64     `json:"department_id"`
	Level        int16     `json:"level"`
	InviteCode   string    `json:"invite_code"`
	CreatedAt    time.Time `json:"created_at"`
}

func toCohortResponse(c *domain.Cohort) *CohortResponse {
	return &CohortResponse{
		ID:           c.ID,
		DepartmentID: c.DepartmentID,
		Level:        c.Level,
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
	CohortID    int64      `json:"cohort_id"`
	UserID      int64      `json:"user_id"`
	Status      string     `json:"status"`
	RequestedAt time.Time  `json:"requested_at"`
	ResolvedAt  *time.Time `json:"resolved_at"`
}

func toJoinRequestResponse(jr *domain.JoinRequest) *JoinRequestResponse {
	return &JoinRequestResponse{
		ID:          jr.ID,
		CohortID:    jr.CohortID,
		UserID:      jr.UserID,
		Status:      string(jr.Status),
		RequestedAt: jr.RequestedAt,
		ResolvedAt:  jr.ResolvedAt,
	}
}

type MemberResponse struct {
	UserID   int64     `json:"user_id"`
	CohortID int64     `json:"cohort_id"`
	JoinedAt time.Time `json:"joined_at"`
}

func toMemberResponse(m *domain.CohortMember) *MemberResponse {
	return &MemberResponse{
		UserID:   m.UserID,
		CohortID: m.CohortID,
		JoinedAt: m.JoinedAt,
	}
}
