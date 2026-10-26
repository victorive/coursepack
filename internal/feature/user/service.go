package user

import (
	"context"
	"log/slog"

	"github.com/victorive/coursepack/internal/pkg/apperr"
)

type Service interface {
	UpdateProfile(ctx context.Context, userUUID string, req *UpdateProfileRequest) error
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{repo: repo, logger: logger}
}

func (s *service) UpdateProfile(ctx context.Context, userUUID string, req *UpdateProfileRequest) error {
	u, err := s.repo.FindByUUID(ctx, userUUID)
	if err != nil {
		return apperr.Internal(err)
	}
	if u == nil {
		return apperr.NotFound("user")
	}
	if req.Phone != nil {
		if err := s.repo.UpdatePhone(ctx, u.ID, *req.Phone); err != nil {
			return apperr.Internal(err)
		}
	}

	if req.UniversityID != nil {
		if err := s.repo.UpdateUniversity(ctx, u.ID, *req.UniversityID); err != nil {
			return apperr.Internal(err)
		}
	}

	return nil
}
