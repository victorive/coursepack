package academic

import (
	"context"

	"github.com/victorive/klaro-backend/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	CreateSemester(ctx context.Context, s *domain.Semester) error
	FindSemestersByClassID(ctx context.Context, classID int64) ([]*domain.Semester, error)
	FindSemesterByID(ctx context.Context, id int64) (*domain.Semester, error)
	UpdateSemester(ctx context.Context, s *domain.Semester) error
	UnsetCurrentSemesterForClass(ctx context.Context, classID int64) error
	DeleteSemester(ctx context.Context, id int64) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateSemester(ctx context.Context, s *domain.Semester) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *repository) FindSemestersByClassID(ctx context.Context, classID int64) ([]*domain.Semester, error) {
	var semesters []*domain.Semester
	if err := r.db.WithContext(ctx).Where("class_id = ?", classID).Order("session DESC, kind ASC").Find(&semesters).Error; err != nil {
		return nil, err
	}
	return semesters, nil
}

func (r *repository) FindSemesterByID(ctx context.Context, id int64) (*domain.Semester, error) {
	var s domain.Semester
	if err := r.db.WithContext(ctx).First(&s, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *repository) UpdateSemester(ctx context.Context, s *domain.Semester) error {
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *repository) UnsetCurrentSemesterForClass(ctx context.Context, classID int64) error {
	return r.db.WithContext(ctx).
		Model(&domain.Semester{}).
		Where("class_id = ? AND is_current = TRUE", classID).
		Update("is_current", false).Error
}

func (r *repository) DeleteSemester(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Semester{}, id).Error
}
