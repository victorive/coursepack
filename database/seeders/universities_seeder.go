package seeders

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/victorive/coursepack/internal/domain"
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
		var existing domain.University
		err := db.Where(domain.University{Acronym: universities[i].Acronym}).First(&existing).Error
		if err == nil {
			continue
		}
		
		universities[i].UUID = uuid.NewString()
		if err := db.Create(&universities[i]).Error; err != nil {
			logger.Error("failed to seed university", "error", err, "acronym", universities[i].Acronym)
			return err
		}
	}

	logger.Info("universities seeded successfully", "count", len(universities))
	return nil
}

