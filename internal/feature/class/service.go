package class

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/victorive/coursepack/internal/domain"
	"github.com/victorive/coursepack/internal/feature/rbac"
	"github.com/victorive/coursepack/internal/pkg/apperr"
	"github.com/victorive/coursepack/internal/pkg/hash"
	"github.com/victorive/coursepack/internal/pkg/pagination"
)

type Service interface {
	Create(ctx context.Context, creatorUUID string, req *CreateClassRequest) (*ClassResponse, error)
	Update(ctx context.Context, updaterUUID string, classUUID string, req *UpdateClassRequest) (*ClassResponse, error)
	List(ctx context.Context, userUUID string, q pagination.Query) ([]*ClassResponse, int64, error)
	Get(ctx context.Context, userUUID string, classUUID string) (*ClassResponse, error)
	Delete(ctx context.Context, uuid string) error
	RequestJoin(ctx context.Context, userUUID string, req *JoinRequest) (*JoinRequestResponse, error)
	ListPendingJoinRequests(ctx context.Context, classUUID string) ([]*JoinRequestResponse, error)
	ReviewJoinRequest(ctx context.Context, reviewerUUID string, classUUID string, requestID int64, req *ReviewJoinRequest) (*JoinRequestResponse, error)
	ListMembers(ctx context.Context, classUUID string) ([]*MemberResponse, error)
	RemoveMember(ctx context.Context, classUUID string, userUUID string) error
	ListLevels(ctx context.Context) ([]*LevelResponse, error)
	RegenerateInviteCode(ctx context.Context, userUUID string, classUUID string) (*ClassResponse, error)
}

type UserRepository interface {
	FindByUUID(ctx context.Context, uuid string) (*domain.User, error)
}

type DepartmentRepository interface {
	FindByID(ctx context.Context, id int64) (*domain.Department, error)
}

type service struct {
	repo     Repository
	userRepo UserRepository
	deptRepo DepartmentRepository
	rbacSvc  rbac.Service
	logger   *slog.Logger
}

func NewService(repo Repository, userRepo UserRepository, deptRepo DepartmentRepository, rbacSvc rbac.Service, logger *slog.Logger) Service {
	return &service{repo: repo, userRepo: userRepo, deptRepo: deptRepo, rbacSvc: rbacSvc, logger: logger}
}

func (s *service) resolveUser(ctx context.Context, uuid string) (*domain.User, error) {
	u, err := s.userRepo.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if u == nil {
		return nil, apperr.NotFound("user")
	}
	return u, nil
}

func (s *service) resolveUserID(ctx context.Context, uuid string) (int64, error) {
	u, err := s.resolveUser(ctx, uuid)
	if err != nil {
		return 0, err
	}
	return u.ID, nil
}

func (s *service) Create(ctx context.Context, creatorUUID string, req *CreateClassRequest) (*ClassResponse, error) {
	creator, err := s.resolveUser(ctx, creatorUUID)
	if err != nil {
		return nil, err
	}

	dept, err := s.deptRepo.FindByID(ctx, req.DepartmentID)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if dept == nil {
		return nil, apperr.NotFound("department")
	}
	if creator.UniversityID == nil || dept.UniversityID != *creator.UniversityID {
		return nil, apperr.New("invalid department", http.StatusUnprocessableEntity)
	}

	code, err := hash.GenerateInviteCode()
	if err != nil {
		return nil, apperr.New("failed to generate invite code", http.StatusInternalServerError).Wrap(err)
	}

	c := &domain.Class{
		UUID:             uuid.NewString(),
		DepartmentID:     req.DepartmentID,
		AdmissionSession: req.AdmissionSession,
		CurrentLevelID:   req.CurrentLevelID,
		Name:             req.Name,
		InviteCode:       code,
		CreatedBy:        &creator.ID,
	}

	err = s.repo.Transact(ctx, func(ctx context.Context) error {
		if err := s.repo.Create(ctx, c); err != nil {
			return apperr.New("failed to create class", http.StatusInternalServerError).Wrap(err)
		}

		if err := s.rbacSvc.AssignClassRole(ctx, creator.ID, c.ID, domain.RoleClassRep); err != nil {
			return apperr.New("failed to assign class rep role", http.StatusInternalServerError).Wrap(err)
		}

		member := &domain.ClassMember{ClassID: c.ID, UserID: creator.ID}
		if err := s.repo.AddMember(ctx, member); err != nil {
			return apperr.New("failed to add creator as member", http.StatusInternalServerError).Wrap(err)
		}

		if err := s.rbacSvc.AssignClassRole(ctx, creator.ID, c.ID, domain.RoleStudent); err != nil {
			return apperr.New("failed to assign student role", http.StatusInternalServerError).Wrap(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	c.Department = dept
	return toClassResponse(c), nil
}

func (s *service) List(ctx context.Context, userUUID string, q pagination.Query) ([]*ClassResponse, int64, error) {
	user, err := s.resolveUser(ctx, userUUID)
	if err != nil {
		return nil, 0, err
	}
	if user.UniversityID == nil {
		return []*ClassResponse{}, 0, nil
	}
	universityID := *user.UniversityID
	total, err := s.repo.CountByUniversity(ctx, universityID)
	if err != nil {
		return nil, 0, apperr.New("failed to list classes", http.StatusInternalServerError).Wrap(err)
	}
	classes, err := s.repo.FindPaginatedByUniversity(ctx, universityID, q.Offset(), q.Limit())
	if err != nil {
		return nil, 0, apperr.New("failed to list classes", http.StatusInternalServerError).Wrap(err)
	}
	result := make([]*ClassResponse, 0, len(classes))
	for _, c := range classes {
		result = append(result, toClassResponse(c))
	}
	return result, total, nil
}

func (s *service) resolveClass(ctx context.Context, classUUID string) (*domain.Class, error) {
	c, err := s.repo.FindByUUID(ctx, classUUID)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if c == nil {
		return nil, apperr.NotFound("class")
	}

	return c, nil
}

func (s *service) Get(ctx context.Context, userUUID string, classUUID string) (*ClassResponse, error) {
	user, err := s.resolveUser(ctx, userUUID)
	if err != nil {
		return nil, err
	}
	if user.UniversityID == nil {
		return nil, apperr.NotFound("class")
	}

	c, err := s.repo.FindByUUIDAndUniversity(ctx, classUUID, *user.UniversityID)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if c == nil {
		return nil, apperr.NotFound("class")
	}

	return toClassResponse(c), nil
}

func (s *service) Delete(ctx context.Context, classUUID string) error {
	c, err := s.resolveClass(ctx, classUUID)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, c.ID); err != nil {
		return apperr.New("failed to delete class", http.StatusInternalServerError).Wrap(err)
	}
	return nil
}

func (s *service) Update(ctx context.Context, updaterUUID string, classUUID string, req *UpdateClassRequest) (*ClassResponse, error) {
	_, err := s.resolveUserID(ctx, updaterUUID)
	if err != nil {
		return nil, err
	}
	c, err := s.resolveClass(ctx, classUUID)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		c.Name = req.Name
	}
	if req.CurrentLevelID != 0 {
		c.CurrentLevelID = req.CurrentLevelID
	}
	if err := s.repo.Update(ctx, c); err != nil {
		return nil, apperr.New("failed to update class", http.StatusInternalServerError).Wrap(err)
	}
	return toClassResponse(c), nil
}

func (s *service) RequestJoin(ctx context.Context, userUUID string, req *JoinRequest) (*JoinRequestResponse, error) {
	user, err := s.resolveUser(ctx, userUUID)
	if err != nil {
		return nil, err
	}
	userID := user.ID

	cls, err := s.repo.FindByInviteCode(ctx, req.InviteCode)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if cls == nil {
		return nil, apperr.New("invalid invite code", http.StatusNotFound)
	}

	dept, err := s.deptRepo.FindByID(ctx, cls.DepartmentID)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if dept == nil || user.UniversityID == nil || dept.UniversityID != *user.UniversityID {
		return nil, apperr.New("invalid invite code", http.StatusNotFound)
	}

	isMember, err := s.repo.IsMember(ctx, cls.ID, userID)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if isMember {
		return nil, apperr.New("already a member of this class", http.StatusConflict)
	}

	if user.UniversityID != nil {
		isActive, err := s.repo.IsActiveMemberAnywhere(ctx, userID, *user.UniversityID)
		if err != nil {
			return nil, apperr.Internal(err)
		}
		if isActive {
			return nil, apperr.New("you must leave your current class before joining another", http.StatusConflict)
		}
	}

	existing, err := s.repo.FindJoinRequestByClassAndUser(ctx, cls.ID, userID)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if existing != nil && existing.Status == domain.JoinStatusPending {
		return nil, apperr.New("join request already pending", http.StatusConflict)
	}

	jr := &domain.JoinRequest{
		ClassID: cls.ID,
		UserID:  userID,
		Status:  domain.JoinStatusPending,
	}
	if err := s.repo.CreateJoinRequest(ctx, jr); err != nil {
		return nil, apperr.New("failed to create join request", http.StatusInternalServerError).Wrap(err)
	}
	return toJoinRequestResponse(jr), nil
}

func (s *service) ListPendingJoinRequests(ctx context.Context, classUUID string) ([]*JoinRequestResponse, error) {
	c, err := s.resolveClass(ctx, classUUID)
	if err != nil {
		return nil, err
	}

	requests, err := s.repo.FindPendingJoinRequests(ctx, c.ID)
	if err != nil {
		return nil, apperr.New("failed to list join requests", http.StatusInternalServerError).Wrap(err)
	}
	result := make([]*JoinRequestResponse, 0, len(requests))
	for _, jr := range requests {
		result = append(result, toJoinRequestResponse(jr))
	}
	return result, nil
}

func (s *service) ReviewJoinRequest(ctx context.Context, reviewerUUID string, classUUID string, requestID int64, req *ReviewJoinRequest) (*JoinRequestResponse, error) {
	reviewerID, err := s.resolveUserID(ctx, reviewerUUID)
	if err != nil {
		return nil, err
	}
	cls, err := s.resolveClass(ctx, classUUID)
	if err != nil {
		return nil, err
	}
	jr, err := s.repo.FindJoinRequestByID(ctx, requestID)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if jr == nil {
		return nil, apperr.NotFound("join request")
	}
	if jr.ClassID != cls.ID {
		return nil, apperr.NotFound("join request")
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
		m := &domain.ClassMember{ClassID: jr.ClassID, UserID: jr.UserID}
		if err := s.repo.AddMember(ctx, m); err != nil {
			return nil, apperr.New("failed to add member", http.StatusInternalServerError).Wrap(err)
		}
		if err := s.rbacSvc.AssignClassRole(ctx, jr.UserID, jr.ClassID, domain.RoleStudent); err != nil {
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

func (s *service) ListMembers(ctx context.Context, classUUID string) ([]*MemberResponse, error) {
	c, err := s.resolveClass(ctx, classUUID)
	if err != nil {
		return nil, err
	}
	members, err := s.repo.FindMembers(ctx, c.ID)
	if err != nil {
		return nil, apperr.New("failed to list members", http.StatusInternalServerError).Wrap(err)
	}
	result := make([]*MemberResponse, 0, len(members))
	for _, m := range members {
		result = append(result, toMemberResponse(m))
	}
	return result, nil
}

func (s *service) RemoveMember(ctx context.Context, classUUID string, targetUserUUID string) error {
	cls, err := s.resolveClass(ctx, classUUID)
	if err != nil {
		return err
	}

	targetUser, err := s.userRepo.FindByUUID(ctx, targetUserUUID)
	if err != nil {
		return apperr.Internal(err)
	}

	if targetUser == nil {
		return apperr.NotFound("user")
	}

	isMember, err := s.repo.IsMember(ctx, cls.ID, targetUser.ID)
	if err != nil {
		return apperr.Internal(err)
	}
	if !isMember {
		return apperr.NotFound("member")
	}
	if err := s.repo.RemoveMember(ctx, cls.ID, targetUser.ID); err != nil {
		return apperr.New("failed to remove member", http.StatusInternalServerError).Wrap(err)
	}
	return nil
}

func (s *service) ListLevels(ctx context.Context) ([]*LevelResponse, error) {
	levels, err := s.repo.FindAllLevels(ctx)
	if err != nil {
		return nil, apperr.New("failed to fetch levels", http.StatusInternalServerError).Wrap(err)
	}

	resp := make([]*LevelResponse, len(levels))
	for i, l := range levels {
		resp[i] = toLevelResponse(l)
	}

	return resp, nil
}

func (s *service) RegenerateInviteCode(ctx context.Context, userUUID string, classUUID string) (*ClassResponse, error) {
	_, err := s.resolveUserID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	c, err := s.resolveClass(ctx, classUUID)
	if err != nil {
		return nil, err
	}

	code, err := hash.GenerateInviteCode()
	if err != nil {
		return nil, apperr.New("failed to generate invite code", http.StatusInternalServerError).Wrap(err)
	}

	if err := s.repo.UpdateInviteCode(ctx, c.ID, code); err != nil {
		return nil, apperr.New("failed to update invite code", http.StatusInternalServerError).Wrap(err)
	}

	c.InviteCode = code
	return toClassResponse(c), nil
}
