package seeders

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/victorive/klaro-backend/database/seeders/json/courses"
	"github.com/victorive/klaro-backend/internal/domain"
	"gorm.io/gorm"
)

type coursesSeeder struct{}

func (s coursesSeeder) Name() string {
	return "courses"
}

type jsonCourseEntry struct {
	Code        string `json:"code"`
	Title       string `json:"title"`
	CreditUnits int16  `json:"credit_units"`
}

type jsonCoursesFile struct {
	Departments map[string][]jsonCourseEntry `json:"departments"`
}

func (s coursesSeeder) Run(db *gorm.DB) error {
	logger := slog.Default()

	var file jsonCoursesFile
	if err := json.Unmarshal(courses.UNNCoursesJSON, &file); err != nil {
		return fmt.Errorf("failed to parse unn.json: %w", err)
	}

	deptMap := map[string]domain.Department{}
	for deptCode := range file.Departments {
		var dept domain.Department
		if err := db.Where("code = ?", deptCode).First(&dept).Error; err != nil {
			continue
		}
		deptMap[deptCode] = dept
	}

	seeded, skipped := 0, 0
	for deptCode, courses := range file.Departments {
		dept, ok := deptMap[deptCode]
		if !ok {
			skipped += len(courses)
			continue
		}

		for _, e := range courses {
			title := strings.TrimSpace(e.Title)
			if title == "" {
				continue
			}

			course := domain.Course{
				UniversityID: dept.UniversityID,
				DepartmentID: dept.ID,
				Code:         e.Code,
				Title:        title,
				CreditUnits:  e.CreditUnits,
			}

			var existing domain.Course
			if db.Where("code = ? AND university_id = ?", course.Code, course.UniversityID).First(&existing).Error == nil {
				continue
			}

			course.UUID = uuid.NewString()
			if err := db.Create(&course).Error; err != nil {
				logger.Error("failed to seed course", "code", course.Code, "error", err)
				return err
			}
			seeded++
		}
	}

	logger.Info("courses seeded", "inserted", seeded, "skipped_no_dept", skipped)
	return nil
}
