package file

import (
	"context"

	"github.com/victorive/klaro-backend/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, f *domain.File) error
	FindByID(ctx context.Context, id int64) (*domain.File, error)
	FindByClassCourse(ctx context.Context, classCourseID int64) ([]*domain.File, error)
	FindByLecture(ctx context.Context, lectureID int64) ([]*domain.File, error)
	Update(ctx context.Context, f *domain.File) error
	Delete(ctx context.Context, id int64) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, f *domain.File) error {
	return r.db.WithContext(ctx).Create(f).Error
}

func (r *repository) FindByID(ctx context.Context, id int64) (*domain.File, error) {
	var f domain.File
	if err := r.db.WithContext(ctx).First(&f, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &f, nil
}

func (r *repository) FindByClassCourse(ctx context.Context, classCourseID int64) ([]*domain.File, error) {
	var files []*domain.File
	if err := r.db.WithContext(ctx).
		Where("class_course_id = ? AND status != 'archived'", classCourseID).
		Order("created_at DESC").
		Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (r *repository) FindByLecture(ctx context.Context, lectureID int64) ([]*domain.File, error) {
	var files []*domain.File
	if err := r.db.WithContext(ctx).
		Where("lecture_id = ? AND status != 'archived'", lectureID).
		Order("created_at DESC").
		Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (r *repository) Update(ctx context.Context, f *domain.File) error {
	return r.db.WithContext(ctx).Save(f).Error
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.File{}, id).Error
}
