package department

import (
	"context"

	"github.com/victorive/klaro-backend/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, d *domain.Department) error
	FindAll(ctx context.Context) ([]*domain.Department, error)
	Count(ctx context.Context) (int64, error)
	FindPaginated(ctx context.Context, offset, limit int) ([]*domain.Department, error)
	FindByID(ctx context.Context, id int64) (*domain.Department, error)
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

func (r *repository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Department{}).Count(&count).Error; err != nil {
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

func (r *repository) Update(ctx context.Context, d *domain.Department) error {
	return r.db.WithContext(ctx).Save(d).Error
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Department{}, id).Error
}
