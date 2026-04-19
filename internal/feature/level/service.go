package level

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/victorive/klaro-backend/internal/domain"
	"github.com/victorive/klaro-backend/internal/feature/rbac"
	"github.com/victorive/klaro-backend/internal/pkg/apperr"
	"github.com/victorive/klaro-backend/internal/pkg/hash"
	"github.com/victorive/klaro-backend/internal/pkg/pagination"
)

type Service interface {
	Create(ctx context.Context, creatorUUID string, req *CreateLevelRequest) (*LevelResponse, error)
	List(ctx context.Context, q pagination.Query) ([]*LevelResponse, int64, error)
	Get(ctx context.Context, id int64) (*LevelResponse, error)
	Delete(ctx context.Context, id int64) error
	RequestJoin(ctx context.Context, userUUID string, req *JoinRequest) (*JoinRequestResponse, error)
	ListPendingJoinRequests(ctx context.Context, levelID int64) ([]*JoinRequestResponse, error)
	ReviewJoinRequest(ctx context.Context, reviewerUUID string, levelID, requestID int64, req *ReviewJoinRequest) (*JoinRequestResponse, error)
	ListMembers(ctx context.Context, levelID int64) ([]*MemberResponse, error)
	RemoveMember(ctx context.Context, levelID, userID int64) error
}

type UserRepository interface {
	FindByUUID(ctx context.Context, uuid string) (*domain.User, error)
}

type service struct {
	repo     Repository
	userRepo UserRepository
	rbacSvc  rbac.Service
	logger   *slog.Logger
}

func NewService(repo Repository, userRepo UserRepository, rbacSvc rbac.Service, logger *slog.Logger) Service {
	return &service{repo: repo, userRepo: userRepo, rbacSvc: rbacSvc, logger: logger}
}

func (s *service) resolveUserID(ctx context.Context, uuid string) (int64, error) {
	u, err := s.userRepo.FindByUUID(ctx, uuid)
	if err != nil {
		return 0, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if u == nil {
		return 0, apperr.New("user not found", http.StatusNotFound)
	}
	return u.ID, nil
}

func (s *service) Create(ctx context.Context, creatorUUID string, req *CreateLevelRequest) (*LevelResponse, error) {
	creatorID, err := s.resolveUserID(ctx, creatorUUID)
	if err != nil {
		return nil, err
	}

	code, err := hash.GenerateInviteCode()
	if err != nil {
		return nil, apperr.New("failed to generate invite code", http.StatusInternalServerError).Wrap(err)
	}

	c := &domain.Level{
		DepartmentID: req.DepartmentID,
		Year:         req.Year,
		InviteCode:   code,
		CreatedBy:    &creatorID,
	}

	if err := s.repo.Create(ctx, c); err != nil {
		return nil, apperr.New("failed to create level", http.StatusInternalServerError).Wrap(err)
	}

	if err := s.rbacSvc.AssignLevelRole(ctx, creatorID, c.ID, domain.RoleClassRep); err != nil {
		s.logger.Warn("failed to assign class_rep role", "error", err)
	}

	return toLevelResponse(c), nil
}

func (s *service) List(ctx context.Context, q pagination.Query) ([]*LevelResponse, int64, error) {
	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, apperr.New("failed to list levels", http.StatusInternalServerError).Wrap(err)
	}
	levels, err := s.repo.FindPaginated(ctx, q.Offset(), q.Limit())
	if err != nil {
		return nil, 0, apperr.New("failed to list levels", http.StatusInternalServerError).Wrap(err)
	}
	result := make([]*LevelResponse, 0, len(levels))
	for _, c := range levels {
		result = append(result, toLevelResponse(c))
	}
	return result, total, nil
}

func (s *service) Get(ctx context.Context, id int64) (*LevelResponse, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if c == nil {
		return nil, apperr.New("level not found", http.StatusNotFound)
	}
	return toLevelResponse(c), nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if c == nil {
		return apperr.New("level not found", http.StatusNotFound)
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return apperr.New("failed to delete level", http.StatusInternalServerError).Wrap(err)
	}
	return nil
}

func (s *service) RequestJoin(ctx context.Context, userUUID string, req *JoinRequest) (*JoinRequestResponse, error) {
	userID, err := s.resolveUserID(ctx, userUUID)
	if err != nil {
		return nil, err
	}
	level, err := s.repo.FindByInviteCode(ctx, req.InviteCode)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if level == nil {
		return nil, apperr.New("invalid invite code", http.StatusNotFound)
	}

	isMember, err := s.repo.IsMember(ctx, level.ID, userID)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if isMember {
		return nil, apperr.New("already a member of this level", http.StatusConflict)
	}

	existing, err := s.repo.FindJoinRequestByLevelAndUser(ctx, level.ID, userID)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if existing != nil && existing.Status == domain.JoinStatusPending {
		return nil, apperr.New("join request already pending", http.StatusConflict)
	}

	jr := &domain.JoinRequest{
		LevelID: level.ID,
		UserID:  userID,
		Status:  domain.JoinStatusPending,
	}
	if err := s.repo.CreateJoinRequest(ctx, jr); err != nil {
		return nil, apperr.New("failed to create join request", http.StatusInternalServerError).Wrap(err)
	}
	return toJoinRequestResponse(jr), nil
}

func (s *service) ListPendingJoinRequests(ctx context.Context, levelID int64) ([]*JoinRequestResponse, error) {
	requests, err := s.repo.FindPendingJoinRequests(ctx, levelID)
	if err != nil {
		return nil, apperr.New("failed to list join requests", http.StatusInternalServerError).Wrap(err)
	}
	result := make([]*JoinRequestResponse, 0, len(requests))
	for _, jr := range requests {
		result = append(result, toJoinRequestResponse(jr))
	}
	return result, nil
}

func (s *service) ReviewJoinRequest(ctx context.Context, reviewerUUID string, levelID, requestID int64, req *ReviewJoinRequest) (*JoinRequestResponse, error) {
	reviewerID, err := s.resolveUserID(ctx, reviewerUUID)
	if err != nil {
		return nil, err
	}
	jr, err := s.repo.FindJoinRequestByID(ctx, requestID)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if jr == nil {
		return nil, apperr.New("join request not found", http.StatusNotFound)
	}
	if jr.LevelID != levelID {
		return nil, apperr.New("join request not found", http.StatusNotFound)
	}
	if jr.Status != domain.JoinStatusPending {
		return nil, apperr.New("join request already resolved", http.StatusConflict)
	}

	now := time.Now()
	jr.ReviewedBy = &reviewerID
	jr.ResolvedAt = &now

	if req.Action == "approve" {
		jr.Status = domain.JoinStatusApproved
		if err := s.repo.UpdateJoinRequest(ctx, jr); err != nil {
			return nil, apperr.New("failed to approve join request", http.StatusInternalServerError).Wrap(err)
		}
		m := &domain.LevelMember{LevelID: jr.LevelID, UserID: jr.UserID}
		if err := s.repo.AddMember(ctx, m); err != nil {
			return nil, apperr.New("failed to add member", http.StatusInternalServerError).Wrap(err)
		}
		if err := s.rbacSvc.AssignLevelRole(ctx, jr.UserID, jr.LevelID, domain.RoleStudent); err != nil {
			s.logger.Warn("failed to assign student role", "error", err)
		}
	} else {
		jr.Status = domain.JoinStatusRejected
		if err := s.repo.UpdateJoinRequest(ctx, jr); err != nil {
			return nil, apperr.New("failed to reject join request", http.StatusInternalServerError).Wrap(err)
		}
	}

	return toJoinRequestResponse(jr), nil
}

func (s *service) ListMembers(ctx context.Context, levelID int64) ([]*MemberResponse, error) {
	members, err := s.repo.FindMembers(ctx, levelID)
	if err != nil {
		return nil, apperr.New("failed to list members", http.StatusInternalServerError).Wrap(err)
	}
	result := make([]*MemberResponse, 0, len(members))
	for _, m := range members {
		result = append(result, toMemberResponse(m))
	}
	return result, nil
}

func (s *service) RemoveMember(ctx context.Context, levelID, userID int64) error {
	isMember, err := s.repo.IsMember(ctx, levelID, userID)
	if err != nil {
		return apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if !isMember {
		return apperr.New("user is not a member of this level", http.StatusNotFound)
	}
	if err := s.repo.RemoveMember(ctx, levelID, userID); err != nil {
		return apperr.New("failed to remove member", http.StatusInternalServerError).Wrap(err)
	}
	return nil
}
