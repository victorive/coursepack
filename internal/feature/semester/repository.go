package semester

import (
	"context"

	"github.com/victorive/coursepack/internal/domain"
	"github.com/victorive/coursepack/internal/pkg/dbtx"
	"gorm.io/gorm"
)

type Repository interface {
	CreateSemester(ctx context.Context, s *domain.Semester) error
	FindSemestersByClassID(ctx context.Context, classID int64) ([]*domain.Semester, error)
	FindSemesterByUUID(ctx context.Context, uuid string) (*domain.Semester, error)
	FindSemesterByClassSessionKind(ctx context.Context, classID int64, session, kind string) (*domain.Semester, error)
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
	return dbtx.Conn(ctx, r.db).Create(s).Error
}

func (r *repository) FindSemestersByClassID(ctx context.Context, classID int64) ([]*domain.Semester, error) {
	var semesters []*domain.Semester
	if err := dbtx.Conn(ctx, r.db).Where("class_id = ?", classID).Order("session DESC, kind ASC").Find(&semesters).Error; err != nil {
		return nil, err
	}
	return semesters, nil
}

func (r *repository) FindSemesterByUUID(ctx context.Context, uuid string) (*domain.Semester, error) {
	var s domain.Semester
	if err := dbtx.Conn(ctx, r.db).Where("uuid = ?", uuid).First(&s).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *repository) FindSemesterByClassSessionKind(ctx context.Context, classID int64, session, kind string) (*domain.Semester, error) {
	var s domain.Semester
	if err := dbtx.Conn(ctx, r.db).Where("class_id = ? AND session = ? AND kind = ?", classID, session, kind).First(&s).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *repository) UpdateSemester(ctx context.Context, s *domain.Semester) error {
	return dbtx.Conn(ctx, r.db).Save(s).Error
}

func (r *repository) UnsetCurrentSemesterForClass(ctx context.Context, classID int64) error {
	return dbtx.Conn(ctx, r.db).
		Model(&domain.Semester{}).
		Where("class_id = ? AND is_current = TRUE", classID).
		Update("is_current", false).Error
}

func (r *repository) DeleteSemester(ctx context.Context, id int64) error {
	return dbtx.Conn(ctx, r.db).Delete(&domain.Semester{}, id).Error
}
