package university

import (
	"context"

	"github.com/victorive/coursepack/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(ctx context.Context) ([]*domain.University, error)
	FindByID(ctx context.Context, id int64) (*domain.University, error)
	FindByUUID(ctx context.Context, uuid string) (*domain.University, error)
	FindByAcronym(ctx context.Context, acronym string) (*domain.University, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(ctx context.Context) ([]*domain.University, error) {
	var universities []*domain.University
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&universities).Error; err != nil {
		return nil, err
	}
	return universities, nil
}

func (r *repository) FindByID(ctx context.Context, id int64) (*domain.University, error) {
	var u domain.University
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *repository) FindByUUID(ctx context.Context, uuid string) (*domain.University, error) {
	var u domain.University
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	
	return &u, nil
}

func (r *repository) FindByAcronym(ctx context.Context, acronym string) (*domain.University, error) {
	var u domain.University
	if err := r.db.WithContext(ctx).Where("acronym = ?", acronym).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
