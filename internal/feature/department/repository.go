package department

import (
	"context"

	"github.com/victorive/coursepack/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, d *domain.Department) error
	FindAll(ctx context.Context) ([]*domain.Department, error)
	FindAllByUniversity(ctx context.Context, universityID int64) ([]*domain.Department, error)
	Count(ctx context.Context) (int64, error)
	CountByUniversity(ctx context.Context, universityID int64) (int64, error)
	FindPaginated(ctx context.Context, offset, limit int) ([]*domain.Department, error)
	FindPaginatedByUniversity(ctx context.Context, universityID int64, offset, limit int) ([]*domain.Department, error)
	FindByID(ctx context.Context, id int64) (*domain.Department, error)
	FindByUUID(ctx context.Context, uuid string) (*domain.Department, error)
	FindByIDAndUniversity(ctx context.Context, id int64, universityID int64) (*domain.Department, error)
	FindByUUIDAndUniversity(ctx context.Context, uuid string, universityID int64) (*domain.Department, error)
	Update(ctx context.Context, d *domain.Department) error
	Delete(ctx context.Context, id int64) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, d *domain.Department) error {
	return r.db.WithContext(ctx).Create(d).Error
}

func (r *repository) FindAll(ctx context.Context) ([]*domain.Department, error) {
	var departments []*domain.Department
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&departments).Error; err != nil {
		return nil, err
	}
	return departments, nil
}

func (r *repository) FindAllByUniversity(ctx context.Context, universityID int64) ([]*domain.Department, error) {
	var departments []*domain.Department
	if err := r.db.WithContext(ctx).
		Where("university_id = ?", universityID).
		Order("name ASC").
		Find(&departments).Error; err != nil {
		return nil, err
	}

	return departments, nil
}

func (r *repository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Department{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) CountByUniversity(ctx context.Context, universityID int64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&domain.Department{}).
		Where("university_id = ?", universityID).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *repository) FindPaginated(ctx context.Context, offset, limit int) ([]*domain.Department, error) {
	var departments []*domain.Department
	if err := r.db.WithContext(ctx).Order("name ASC").Offset(offset).Limit(limit).Find(&departments).Error; err != nil {
		return nil, err
	}
	return departments, nil
}

func (r *repository) FindPaginatedByUniversity(ctx context.Context, universityID int64, offset, limit int) ([]*domain.Department, error) {
	var departments []*domain.Department
	if err := r.db.WithContext(ctx).
		Preload("University").
		Where("university_id = ?", universityID).
		Order("name ASC").
		Offset(offset).
		Limit(limit).
		Find(&departments).Error; err != nil {
		return nil, err
	}

	return departments, nil
}

func (r *repository) FindByID(ctx context.Context, id int64) (*domain.Department, error) {
	var d domain.Department
	if err := r.db.WithContext(ctx).First(&d, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &d, nil
}

func (r *repository) FindByUUID(ctx context.Context, uuid string) (*domain.Department, error) {
	var d domain.Department
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&d).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &d, nil
}

func (r *repository) FindByIDAndUniversity(ctx context.Context, id int64, universityID int64) (*domain.Department, error) {
	var d domain.Department
	if err := r.db.WithContext(ctx).
		Preload("University").
		Where("id = ? AND university_id = ?", id, universityID).
		First(&d).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &d, nil
}

func (r *repository) FindByUUIDAndUniversity(ctx context.Context, uuid string, universityID int64) (*domain.Department, error) {
	var d domain.Department
	if err := r.db.WithContext(ctx).
		Preload("University").
		Where("uuid = ? AND university_id = ?", uuid, universityID).
		First(&d).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &d, nil
}

func (r *repository) Update(ctx context.Context, d *domain.Department) error {
	return r.db.WithContext(ctx).Save(d).Error
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Department{}, id).Error
}
