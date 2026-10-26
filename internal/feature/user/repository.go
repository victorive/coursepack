package user

import (
	"context"
	"errors"

	"github.com/victorive/coursepack/internal/domain"
	"github.com/victorive/coursepack/internal/pkg/dbtx"
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
	err := dbtx.Conn(ctx, r.db).Where("uuid = ?", uuid).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}

func (r *repository) UpdatePhone(ctx context.Context, id int64, phone string) error {
	return dbtx.Conn(ctx, r.db).Model(&domain.User{}).Where("id = ?", id).
		Update("phone", phone).Error
}

func (r *repository) UpdateUniversity(ctx context.Context, id int64, universityID int64) error {
	return dbtx.Conn(ctx, r.db).Model(&domain.User{}).Where("id = ?", id).
		Update("university_id", universityID).Error
}
