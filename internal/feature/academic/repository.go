package academic

import (
	"context"

	"github.com/victorive/klaro-backend/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	CreateSession(ctx context.Context, s *domain.AcademicSession) error
	CountSessions(ctx context.Context) (int64, error)
	FindSessionsPaginated(ctx context.Context, offset, limit int) ([]*domain.AcademicSession, error)
	FindSessionByID(ctx context.Context, id int64) (*domain.AcademicSession, error)
	UpdateSession(ctx context.Context, s *domain.AcademicSession) error
	UnsetCurrentSession(ctx context.Context) error
	CreateSemester(ctx context.Context, s *domain.Semester) error
	FindSemestersByCohortID(ctx context.Context, cohortID int64) ([]*domain.Semester, error)
	FindSemesterByID(ctx context.Context, id int64) (*domain.Semester, error)
	UpdateSemester(ctx context.Context, s *domain.Semester) error
	UnsetCurrentSemesterForCohort(ctx context.Context, cohortID int64) error
	DeleteSemester(ctx context.Context, id int64) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateSession(ctx context.Context, s *domain.AcademicSession) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *repository) CountSessions(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.AcademicSession{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) FindSessionsPaginated(ctx context.Context, offset, limit int) ([]*domain.AcademicSession, error) {
	var sessions []*domain.AcademicSession
	if err := r.db.WithContext(ctx).Order("label DESC").Offset(offset).Limit(limit).Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *repository) FindSessionByID(ctx context.Context, id int64) (*domain.AcademicSession, error) {
	var s domain.AcademicSession
	if err := r.db.WithContext(ctx).First(&s, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *repository) UpdateSession(ctx context.Context, s *domain.AcademicSession) error {
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *repository) UnsetCurrentSession(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Model(&domain.AcademicSession{}).
		Where("is_current = TRUE").
		Update("is_current", false).Error
}

func (r *repository) CreateSemester(ctx context.Context, s *domain.Semester) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *repository) FindSemestersByCohortID(ctx context.Context, cohortID int64) ([]*domain.Semester, error) {
	var semesters []*domain.Semester
	if err := r.db.WithContext(ctx).Where("cohort_id = ?", cohortID).Order("session_id DESC, kind ASC").Find(&semesters).Error; err != nil {
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

func (r *repository) UnsetCurrentSemesterForCohort(ctx context.Context, cohortID int64) error {
	return r.db.WithContext(ctx).
		Model(&domain.Semester{}).
		Where("cohort_id = ? AND is_current = TRUE", cohortID).
		Update("is_current", false).Error
}

func (r *repository) DeleteSemester(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Semester{}, id).Error
}
