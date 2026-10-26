package file

import (
	"context"

	"github.com/victorive/coursepack/internal/domain"
	"github.com/victorive/coursepack/internal/pkg/dbtx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type pointRepository struct {
	db *gorm.DB
}

func NewPointRepository(db *gorm.DB) PointRepository {
	return &pointRepository{db: db}
}

func (r *pointRepository) HasTransactionForSource(ctx context.Context, userID int64, sourceType domain.PointSource, sourceID int64) (bool, error) {
	var count int64
	err := dbtx.Conn(ctx, r.db).Model(&domain.PointTransaction{}).
		Where("user_id = ? AND source_type = ? AND source_id = ?", userID, sourceType, sourceID).
		Count(&count).Error
	return count > 0, err
}

func (r *pointRepository) AwardPoints(ctx context.Context, userID int64, delta int, sourceType domain.PointSource, sourceID int64, awardedBy int64) error {
	return dbtx.Conn(ctx, r.db).Transaction(func(tx *gorm.DB) error {
		up := domain.UserPoint{UserID: userID}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{}),
		}).FirstOrCreate(&up).Error; err != nil {
			return err
		}

		newBalance := up.Balance + delta
		newTotalEarned := up.TotalEarned + delta

		if err := tx.Model(&domain.UserPoint{}).Where("user_id = ?", userID).
			Updates(map[string]interface{}{
				"balance":      newBalance,
				"total_earned": newTotalEarned,
			}).Error; err != nil {
			return err
		}

		pt := &domain.PointTransaction{
			UserID:       userID,
			Delta:        delta,
			BalanceAfter: newBalance,
			SourceType:   sourceType,
			SourceID:     &sourceID,
			CreatedBy:    &awardedBy,
		}
		return tx.Create(pt).Error
	})
}
