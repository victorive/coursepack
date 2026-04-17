package auth

import (
	"context"
	"errors"

	"github.com/victorive/klaro-backend/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id int64) (*domain.User, error)
	FindByUUID(ctx context.Context, uuid string) (*domain.User, error)
	Create(ctx context.Context, u *domain.User) error
	UpdateEmailVerified(ctx context.Context, id int64, verified bool) error
	UpdatePassword(ctx context.Context, id int64, hash string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}

func (r *repository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	var u domain.User
	err := r.db.WithContext(ctx).First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}

func (r *repository) FindByUUID(ctx context.Context, uuid string) (*domain.User, error) {
	var u domain.User
	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}

func (r *repository) Create(ctx context.Context, u *domain.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *repository) UpdateEmailVerified(ctx context.Context, id int64, verified bool) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).
		Update("email_verified", verified).Error
}

func (r *repository) UpdatePassword(ctx context.Context, id int64, h string) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).
		Update("password", h).Error
}
