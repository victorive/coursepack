package academic

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/victorive/klaro-backend/internal/domain"
	"github.com/victorive/klaro-backend/internal/pkg/apperr"
)

type Service interface {
	CreateSemester(ctx context.Context, req *CreateSemesterRequest) (*SemesterResponse, error)
	ListSemestersByClass(ctx context.Context, classID int64) ([]*SemesterResponse, error)
	GetSemester(ctx context.Context, id int64) (*SemesterResponse, error)
	UpdateSemester(ctx context.Context, id int64, req *UpdateSemesterRequest) (*SemesterResponse, error)
	DeleteSemester(ctx context.Context, id int64) error
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{repo: repo, logger: logger}
}

func (s *service) CreateSemester(ctx context.Context, req *CreateSemesterRequest) (*SemesterResponse, error) {
	if req.IsCurrent {
		if err := s.repo.UnsetCurrentSemesterForClass(ctx, req.ClassID); err != nil {
			return nil, apperr.New("failed to unset current semester", http.StatusInternalServerError).Wrap(err)
		}
	}
	sem := &domain.Semester{
		ClassID:   req.ClassID,
		Session:   req.Session,
		Kind:      req.Kind,
		IsCurrent: req.IsCurrent,
		StartsAt:  req.StartsAt,
		EndsAt:    req.EndsAt,
	}
	if err := s.repo.CreateSemester(ctx, sem); err != nil {
		return nil, apperr.New("failed to create semester", http.StatusInternalServerError).Wrap(err)
	}
	return toSemesterResponse(sem), nil
}

func (s *service) ListSemestersByClass(ctx context.Context, classID int64) ([]*SemesterResponse, error) {
	semesters, err := s.repo.FindSemestersByClassID(ctx, classID)
	if err != nil {
		return nil, apperr.New("failed to list semesters", http.StatusInternalServerError).Wrap(err)
	}
	result := make([]*SemesterResponse, 0, len(semesters))
	for _, sem := range semesters {
		result = append(result, toSemesterResponse(sem))
	}
	return result, nil
}

func (s *service) GetSemester(ctx context.Context, id int64) (*SemesterResponse, error) {
	sem, err := s.repo.FindSemesterByID(ctx, id)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if sem == nil {
		return nil, apperr.New("semester not found", http.StatusNotFound)
	}
	return toSemesterResponse(sem), nil
}

func (s *service) UpdateSemester(ctx context.Context, id int64, req *UpdateSemesterRequest) (*SemesterResponse, error) {
	sem, err := s.repo.FindSemesterByID(ctx, id)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if sem == nil {
		return nil, apperr.New("semester not found", http.StatusNotFound)
	}
	if req.IsCurrent && !sem.IsCurrent {
			if err := s.repo.UnsetCurrentSemesterForClass(ctx, sem.ClassID); err != nil {
			return nil, apperr.New("failed to unset current semester", http.StatusInternalServerError).Wrap(err)
		}
	}
	sem.Kind = req.Kind
	sem.IsCurrent = req.IsCurrent
	sem.StartsAt = req.StartsAt
	sem.EndsAt = req.EndsAt
	if err := s.repo.UpdateSemester(ctx, sem); err != nil {
		return nil, apperr.New("failed to update semester", http.StatusInternalServerError).Wrap(err)
	}
	return toSemesterResponse(sem), nil
}

func (s *service) DeleteSemester(ctx context.Context, id int64) error {
	sem, err := s.repo.FindSemesterByID(ctx, id)
	if err != nil {
		return apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if sem == nil {
		return apperr.New("semester not found", http.StatusNotFound)
	}
	if err := s.repo.DeleteSemester(ctx, id); err != nil {
		return apperr.New("failed to delete semester", http.StatusInternalServerError).Wrap(err)
	}
	return nil
}
