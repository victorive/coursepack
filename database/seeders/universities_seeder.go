package seeders

import (
	"log/slog"

	"github.com/victorive/klaro-backend/internal/domain"
	"gorm.io/gorm"
)

type universitiesSeeder struct{}

func (s universitiesSeeder) Name() string {
	return "universities"
}

func (s universitiesSeeder) Run(db *gorm.DB) error {
	logger := slog.Default()
	universities := []domain.University{
		{Name: "University of Nigeria, Nsukka", Acronym: "UNN"},
		{Name: "University of Abuja", Acronym: "UNIABUJA"},
	}

	for i := range universities {
		if err := db.FirstOrCreate(&universities[i], domain.University{Acronym: universities[i].Acronym}).Error; err != nil {
			logger.Error("failed to seed university", "error", err, "acronym", universities[i].Acronym)
			return err
		}
	}

	logger.Info("universities seeded successfully", "count", len(universities))
	return nil
}

