package department

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/victorive/klaro-backend/internal/domain"
	"github.com/victorive/klaro-backend/internal/pkg/apperr"
	"github.com/victorive/klaro-backend/internal/pkg/pagination"
)

type Service interface {
	Create(ctx context.Context, req *CreateDepartmentRequest) (*DepartmentResponse, error)
	List(ctx context.Context, q pagination.Query) ([]*DepartmentResponse, int64, error)
	Get(ctx context.Context, id int64) (*DepartmentResponse, error)
	Update(ctx context.Context, id int64, req *UpdateDepartmentRequest) (*DepartmentResponse, error)
	Delete(ctx context.Context, id int64) error
	ListByUniversity(ctx context.Context, universityID int64, q pagination.Query) ([]*DepartmentResponse, int64, error)
	GetByIDAndUniversity(ctx context.Context, id int64, universityID int64) (*DepartmentResponse, error)
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{repo: repo, logger: logger}
}

func (s *service) Create(ctx context.Context, req *CreateDepartmentRequest) (*DepartmentResponse, error) {
	d := &domain.Department{
		UniversityID: req.UniversityID,
		Name:         req.Name,
		Code:         req.Code,
		Faculty:      req.Faculty,
	}
	if err := s.repo.Create(ctx, d); err != nil {
		return nil, apperr.New("failed to create department", http.StatusInternalServerError).Wrap(err)
	}
	return toDepartmentResponse(d), nil
}

func (s *service) List(ctx context.Context, q pagination.Query) ([]*DepartmentResponse, int64, error) {
	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, apperr.New("failed to list departments", http.StatusInternalServerError).Wrap(err)
	}
	departments, err := s.repo.FindPaginated(ctx, q.Offset(), q.Limit())
	if err != nil {
		return nil, 0, apperr.New("failed to list departments", http.StatusInternalServerError).Wrap(err)
	}
	result := make([]*DepartmentResponse, 0, len(departments))
	for _, d := range departments {
		result = append(result, toDepartmentResponse(d))
	}
	return result, total, nil
}

func (s *service) Get(ctx context.Context, id int64) (*DepartmentResponse, error) {
	d, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if d == nil {
		return nil, apperr.New("department not found", http.StatusNotFound)
	}
	return toDepartmentResponse(d), nil
}

func (s *service) Update(ctx context.Context, id int64, req *UpdateDepartmentRequest) (*DepartmentResponse, error) {
	d, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if d == nil {
		return nil, apperr.New("department not found", http.StatusNotFound)
	}
	d.Name = req.Name
	d.Code = req.Code
	d.Faculty = req.Faculty
	if err := s.repo.Update(ctx, d); err != nil {
		return nil, apperr.New("failed to update department", http.StatusInternalServerError).Wrap(err)
	}
	return toDepartmentResponse(d), nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	d, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if d == nil {
		return apperr.New("department not found", http.StatusNotFound)
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return apperr.New("failed to delete department", http.StatusInternalServerError).Wrap(err)
	}
	return nil
}

// Scoped methods - use these for user-facing endpoints
func (s *service) ListByUniversity(ctx context.Context, universityID int64, q pagination.Query) ([]*DepartmentResponse, int64, error) {
	total, err := s.repo.CountByUniversity(ctx, universityID)
	if err != nil {
		return nil, 0, apperr.New("failed to list departments", http.StatusInternalServerError).Wrap(err)
	}

	departments, err := s.repo.FindPaginatedByUniversity(ctx, universityID, q.Offset(), q.Limit())
	if err != nil {
		return nil, 0, apperr.New("failed to list departments", http.StatusInternalServerError).Wrap(err)
	}

	result := make([]*DepartmentResponse, 0, len(departments))
	for _, d := range departments {
		result = append(result, toDepartmentResponse(d))
	}

	return result, total, nil
}

func (s *service) GetByIDAndUniversity(ctx context.Context, id int64, universityID int64) (*DepartmentResponse, error) {
	d, err := s.repo.FindByIDAndUniversity(ctx, id, universityID)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}

	if d == nil {
		return nil, apperr.New("department not found", http.StatusNotFound)
	}
	
	return toDepartmentResponse(d), nil
}
