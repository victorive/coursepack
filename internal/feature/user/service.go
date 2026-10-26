package user

import (
	"context"
	"log/slog"

	"github.com/victorive/coursepack/internal/permissions"
	"github.com/victorive/coursepack/internal/pkg/apperr"
)

type Service interface {
	UpdateProfile(ctx context.Context, actor permissions.Actor, req *UpdateProfileRequest) error
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{repo: repo, logger: logger}
}

func (s *service) UpdateProfile(ctx context.Context, actor permissions.Actor, req *UpdateProfileRequest) error {
	if req.Phone != nil {
		if err := s.repo.UpdatePhone(ctx, actor.UserID, *req.Phone); err != nil {
			return apperr.Internal(err)
		}
	}

	if req.UniversityID != nil {
		if err := s.repo.UpdateUniversity(ctx, actor.UserID, *req.UniversityID); err != nil {
			return apperr.Internal(err)
		}
	}

	return nil
}
