package academic

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/victorive/coursepack/internal/domain"
	"github.com/victorive/coursepack/internal/pkg/apperr"
)

type ClassRepository interface {
	FindByUUID(ctx context.Context, uuid string) (*domain.Class, error)
}

type Service interface {
	CreateSemester(ctx context.Context, classUUID string, req *CreateSemesterRequest) (*SemesterResponse, error)
	ListSemestersByClass(ctx context.Context, classUUID string) ([]*SemesterResponse, error)
	GetSemester(ctx context.Context, semUUID string) (*SemesterResponse, error)
	UpdateSemester(ctx context.Context, semUUID string, req *UpdateSemesterRequest) (*SemesterResponse, error)
	DeleteSemester(ctx context.Context, semUUID string) error
}

type service struct {
	repo      Repository
	classRepo ClassRepository
	logger    *slog.Logger
}

func NewService(repo Repository, classRepo ClassRepository, logger *slog.Logger) Service {
	return &service{repo: repo, classRepo: classRepo, logger: logger}
}

func (s *service) CreateSemester(ctx context.Context, classUUID string, req *CreateSemesterRequest) (*SemesterResponse, error) {
	cls, err := s.classRepo.FindByUUID(ctx, classUUID)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if cls == nil {
		return nil, apperr.NotFound("class")
	}

	existing, err := s.repo.FindSemesterByClassSessionKind(ctx, cls.ID, req.Session, req.Kind)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if existing != nil {
		return nil, apperr.New("a "+req.Kind+" semester for session "+req.Session+" already exists for this class", http.StatusConflict)
	}

	if req.IsCurrent {
		if err := s.repo.UnsetCurrentSemesterForClass(ctx, cls.ID); err != nil {
			return nil, apperr.New("failed to unset current semester", http.StatusInternalServerError).Wrap(err)
		}
	}
	var startsAt, endsAt *time.Time
	if req.StartsAt != "" {
		t, _ := time.Parse("2006-01-02", req.StartsAt)
		startsAt = &t
	}

	if req.EndsAt != "" {
		t, _ := time.Parse("2006-01-02", req.EndsAt)
		endsAt = &t
	}
	sem := &domain.Semester{
		UUID:      uuid.NewString(),
		ClassID:   cls.ID,
		Session:   req.Session,
		Kind:      req.Kind,
		IsCurrent: req.IsCurrent,
		StartsAt:  startsAt,
		EndsAt:    endsAt,
	}
	if err := s.repo.CreateSemester(ctx, sem); err != nil {
		return nil, apperr.New("failed to create semester", http.StatusInternalServerError).Wrap(err)
	}
	return ToSemesterResponse(sem), nil
}

func (s *service) ListSemestersByClass(ctx context.Context, classUUID string) ([]*SemesterResponse, error) {
	cls, err := s.classRepo.FindByUUID(ctx, classUUID)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if cls == nil {
		return nil, apperr.NotFound("class")
	}
	semesters, err := s.repo.FindSemestersByClassID(ctx, cls.ID)
	if err != nil {
		return nil, apperr.New("failed to list semesters", http.StatusInternalServerError).Wrap(err)
	}
	result := make([]*SemesterResponse, 0, len(semesters))
	for _, sem := range semesters {
		result = append(result, ToSemesterResponse(sem))
	}
	return result, nil
}

func (s *service) GetSemester(ctx context.Context, semUUID string) (*SemesterResponse, error) {
	sem, err := s.repo.FindSemesterByUUID(ctx, semUUID)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if sem == nil {
		return nil, apperr.NotFound("semester")
	}
	return ToSemesterResponse(sem), nil
}

func (s *service) UpdateSemester(ctx context.Context, semUUID string, req *UpdateSemesterRequest) (*SemesterResponse, error) {
	sem, err := s.repo.FindSemesterByUUID(ctx, semUUID)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if sem == nil {
		return nil, apperr.NotFound("semester")
	}
	if req.IsCurrent && !sem.IsCurrent {
		if err := s.repo.UnsetCurrentSemesterForClass(ctx, sem.ClassID); err != nil {
			return nil, apperr.New("failed to unset current semester", http.StatusInternalServerError).Wrap(err)
		}
	}
	sem.Kind = req.Kind
	sem.IsCurrent = req.IsCurrent
	if req.StartsAt != "" {
		t, _ := time.Parse("2006-01-02", req.StartsAt)
		sem.StartsAt = &t
	} else {
		sem.StartsAt = nil
	}

	if req.EndsAt != "" {
		t, _ := time.Parse("2006-01-02", req.EndsAt)
		sem.EndsAt = &t
	} else {
		sem.EndsAt = nil
	}

	if err := s.repo.UpdateSemester(ctx, sem); err != nil {
		return nil, apperr.New("failed to update semester", http.StatusInternalServerError).Wrap(err)
	}
	return ToSemesterResponse(sem), nil
}

func (s *service) DeleteSemester(ctx context.Context, semUUID string) error {
	sem, err := s.repo.FindSemesterByUUID(ctx, semUUID)
	if err != nil {
		return apperr.Internal(err)
	}
	if sem == nil {
		return apperr.NotFound("semester")
	}
	if err := s.repo.DeleteSemester(ctx, sem.ID); err != nil {
		return apperr.New("failed to delete semester", http.StatusInternalServerError).Wrap(err)
	}
	return nil
}
