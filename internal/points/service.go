package points

import (
	"context"

	"github.com/victorive/coursepack/internal/domain"
)

// Service awards points for repeatable activities.
type Service interface {
	Award(ctx context.Context, userID int64, delta int, sourceType domain.PointSource, sourceID int64, awardedBy int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Award(ctx context.Context, userID int64, delta int, sourceType domain.PointSource, sourceID int64, awardedBy int64) error {
	return s.repo.AwardPoints(ctx, userID, delta, sourceType, sourceID, awardedBy)
}
