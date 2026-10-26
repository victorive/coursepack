package auth

import (
	"context"
	"errors"

	"github.com/victorive/coursepack/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id int64) (*domain.User, error)
	FindByUUID(ctx context.Context, uuid string) (*domain.User, error)
	Create(ctx context.Context, u *domain.User) error
	UpdateEmailVerified(ctx context.Context, id int64, verified bool) error
	UpdatePassword(ctx context.Context, id int64, hash string) error
	UpdateGoogleProfile(ctx context.Context, id int64, fullName string, avatarURL *string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.db.WithContext(ctx).Preload("University").Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}

func (r *repository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	var u domain.User
	err := r.db.WithContext(ctx).Preload("University").First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}

func (r *repository) FindByUUID(ctx context.Context, uuid string) (*domain.User, error) {
	var u domain.User
	err := r.db.WithContext(ctx).Preload("University").Where("uuid = ?", uuid).First(&u).Error
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

func (r *repository) UpdateGoogleProfile(ctx context.Context, id int64, fullName string, avatarURL *string) error {
	updates := map[string]any{"full_name": fullName, "avatar_url": avatarURL}
	return r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Updates(updates).Error
}
