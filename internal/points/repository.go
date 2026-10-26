package points

import (
	"context"

	"github.com/victorive/coursepack/internal/domain"
	"github.com/victorive/coursepack/internal/pkg/dbtx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	AwardPoints(ctx context.Context, userID int64, delta int, sourceType domain.PointSource, sourceID int64, awardedBy int64) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// AwardPoints records one point transaction per source. A duplicate source is
// silently ignored, so the caller can retry without checking first.
func (r *repository) AwardPoints(ctx context.Context, userID int64, delta int, sourceType domain.PointSource, sourceID int64, awardedBy int64) error {
	return dbtx.Conn(ctx, r.db).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoNothing: true,
		}).Create(&domain.UserPoint{UserID: userID}).Error; err != nil {
			return err
		}

		var points domain.UserPoint
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", userID).
			First(&points).Error; err != nil {
			return err
		}

		newBalance := points.Balance + delta
		entry := &domain.PointTransaction{
			UserID:       userID,
			Delta:        delta,
			BalanceAfter: newBalance,
			SourceType:   sourceType,
			SourceID:     &sourceID,
			CreatedBy:    &awardedBy,
		}

		result := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "source_type"}, {Name: "source_id"}},
			DoNothing: true,
		}).Create(entry)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return nil
		}

		return tx.Model(&domain.UserPoint{}).Where("user_id = ?", userID).
			Updates(map[string]interface{}{
				"balance":      newBalance,
				"total_earned": points.TotalEarned + delta,
			}).Error
	})
}
