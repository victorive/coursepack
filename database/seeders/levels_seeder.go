package seeders

import (
	"github.com/victorive/klaro-backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type levelsSeeder struct{}

func (levelsSeeder) Name() string { return "levels" }

func (levelsSeeder) Run(db *gorm.DB) error {
	levels := []domain.Level{
		{Value: 100, Label: "100 Level"},
		{Value: 200, Label: "200 Level"},
		{Value: 300, Label: "300 Level"},
		{Value: 400, Label: "400 Level"},
		{Value: 500, Label: "500 Level"},
		{Value: 600, Label: "600 Level"},
		{Value: 700, Label: "700 Level"},
		{Value: 800, Label: "800 Level"},
	}
	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&levels).Error
}
