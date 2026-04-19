package user

import (
	"context"
	"errors"

	"github.com/victorive/klaro-backend/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	FindByUUID(ctx context.Context, uuid string) (*domain.User, error)
	UpdatePhone(ctx context.Context, id int64, phone string) error
	UpdateUniversity(ctx context.Context, id int64, universityID int64) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByUUID(ctx context.Context, uuid string) (*domain.User, error) {
	var u domain.User
	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}

func (r *repository) UpdatePhone(ctx context.Context, id int64, phone string) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).
		Update("phone", phone).Error
}

func (r *repository) UpdateUniversity(ctx context.Context, id int64, universityID int64) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).
		Update("university_id", universityID).Error
}
