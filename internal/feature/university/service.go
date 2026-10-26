package university

import (
	"context"
	"net/http"

	"github.com/victorive/coursepack/internal/domain"
	"github.com/victorive/coursepack/internal/pkg/apperr"
)

type Service interface {
	ListUniversities(ctx context.Context) ([]*domain.University, error)
	GetByID(ctx context.Context, id int64) (*domain.University, error)
	GetByUUID(ctx context.Context, uuid string) (*domain.University, error)
	GetByAcronym(ctx context.Context, acronym string) (*domain.University, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ListUniversities(ctx context.Context) ([]*domain.University, error) {
	universities, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, apperr.New("failed to fetch universities", http.StatusInternalServerError).Wrap(err)
	}
	return universities, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*domain.University, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperr.New("failed to fetch university", http.StatusInternalServerError).Wrap(err)
	}
	if u == nil {
		return nil, apperr.NotFound("university")
	}
	return u, nil
}

func (s *service) GetByUUID(ctx context.Context, uuid string) (*domain.University, error) {
	u, err := s.repo.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, apperr.New("failed to fetch university", http.StatusInternalServerError).Wrap(err)
	}
	if u == nil {
		return nil, apperr.NotFound("university")
	}

	return u, nil
}

func (s *service) GetByAcronym(ctx context.Context, acronym string) (*domain.University, error) {
	u, err := s.repo.FindByAcronym(ctx, acronym)
	if err != nil {
		return nil, apperr.New("failed to fetch university", http.StatusInternalServerError).Wrap(err)
	}
	if u == nil {
		return nil, apperr.NotFound("university")
	}
	return u, nil
}
