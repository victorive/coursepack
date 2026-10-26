package department

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/victorive/coursepack/internal/pkg/apperr"
	"github.com/victorive/coursepack/internal/pkg/pagination"
)

type Service interface {
	ListByUniversity(ctx context.Context, universityID int64, q pagination.Query) ([]*DepartmentResponse, int64, error)
	GetByIDAndUniversity(ctx context.Context, id int64, universityID int64) (*DepartmentResponse, error)
	GetByUUIDAndUniversity(ctx context.Context, uuid string, universityID int64) (*DepartmentResponse, error)
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{repo: repo, logger: logger}
}

func (s *service) ListByUniversity(ctx context.Context, universityID int64, q pagination.Query) ([]*DepartmentResponse, int64, error) {
	total, err := s.repo.CountByUniversity(ctx, universityID)
	if err != nil {
		return nil, 0, apperr.New("failed to list departments", http.StatusInternalServerError).Wrap(err)
	}

	if total == 0 {
		return []*DepartmentResponse{}, 0, nil
	}

	departments, err := s.repo.FindPaginatedByUniversity(ctx, universityID, q.Offset(), q.Limit())
	if err != nil {
		return nil, 0, apperr.New("failed to list departments", http.StatusInternalServerError).Wrap(err)
	}

	result := make([]*DepartmentResponse, 0, len(departments))
	for _, d := range departments {
		result = append(result, ToDepartmentResponse(d))
	}

	return result, total, nil
}

func (s *service) GetByIDAndUniversity(ctx context.Context, id int64, universityID int64) (*DepartmentResponse, error) {
	d, err := s.repo.FindByIDAndUniversity(ctx, id, universityID)
	if err != nil {
		return nil, apperr.Internal(err)
	}

	if d == nil {
		return nil, apperr.NotFound("department")
	}

	return ToDepartmentResponse(d), nil
}

func (s *service) GetByUUIDAndUniversity(ctx context.Context, uuid string, universityID int64) (*DepartmentResponse, error) {
	d, err := s.repo.FindByUUIDAndUniversity(ctx, uuid, universityID)
	if err != nil {
		return nil, apperr.Internal(err)
	}

	if d == nil {
		return nil, apperr.NotFound("department")
	}

	return ToDepartmentResponse(d), nil
}
