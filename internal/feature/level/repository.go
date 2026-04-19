package level

import (
	"context"

	"github.com/victorive/klaro-backend/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, c *domain.Level) error
	Count(ctx context.Context) (int64, error)
	FindPaginated(ctx context.Context, offset, limit int) ([]*domain.Level, error)
	FindByID(ctx context.Context, id int64) (*domain.Level, error)
	FindByInviteCode(ctx context.Context, code string) (*domain.Level, error)
	Delete(ctx context.Context, id int64) error
	AddMember(ctx context.Context, m *domain.LevelMember) error
	RemoveMember(ctx context.Context, levelID, userID int64) error
	FindMembers(ctx context.Context, levelID int64) ([]*domain.LevelMember, error)
	IsMember(ctx context.Context, levelID, userID int64) (bool, error)
	CreateJoinRequest(ctx context.Context, jr *domain.JoinRequest) error
	FindJoinRequestByID(ctx context.Context, id int64) (*domain.JoinRequest, error)
	FindJoinRequestByLevelAndUser(ctx context.Context, levelID, userID int64) (*domain.JoinRequest, error)
	FindPendingJoinRequests(ctx context.Context, levelID int64) ([]*domain.JoinRequest, error)
	UpdateJoinRequest(ctx context.Context, jr *domain.JoinRequest) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, c *domain.Level) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *repository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Level{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) FindPaginated(ctx context.Context, offset, limit int) ([]*domain.Level, error) {
	var levels []*domain.Level
	if err := r.db.WithContext(ctx).Order("id ASC").Offset(offset).Limit(limit).Find(&levels).Error; err != nil {
		return nil, err
	}
	return levels, nil
}

func (r *repository) FindByID(ctx context.Context, id int64) (*domain.Level, error) {
	var c domain.Level
	if err := r.db.WithContext(ctx).First(&c, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *repository) FindByInviteCode(ctx context.Context, code string) (*domain.Level, error) {
	var c domain.Level
	if err := r.db.WithContext(ctx).Where("invite_code = ?", code).First(&c).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Level{}, id).Error
}

func (r *repository) AddMember(ctx context.Context, m *domain.LevelMember) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *repository) RemoveMember(ctx context.Context, levelID, userID int64) error {
	return r.db.WithContext(ctx).
		Where("level_id = ? AND user_id = ?", levelID, userID).
		Delete(&domain.LevelMember{}).Error
}

func (r *repository) FindMembers(ctx context.Context, levelID int64) ([]*domain.LevelMember, error) {
	var members []*domain.LevelMember
	if err := r.db.WithContext(ctx).Where("level_id = ?", levelID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *repository) IsMember(ctx context.Context, levelID, userID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.LevelMember{}).
		Where("level_id = ? AND user_id = ?", levelID, userID).Count(&count).Error
	return count > 0, err
}

func (r *repository) CreateJoinRequest(ctx context.Context, jr *domain.JoinRequest) error {
	return r.db.WithContext(ctx).Create(jr).Error
}

func (r *repository) FindJoinRequestByID(ctx context.Context, id int64) (*domain.JoinRequest, error) {
	var jr domain.JoinRequest
	if err := r.db.WithContext(ctx).First(&jr, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &jr, nil
}

func (r *repository) FindJoinRequestByLevelAndUser(ctx context.Context, levelID, userID int64) (*domain.JoinRequest, error) {
	var jr domain.JoinRequest
	if err := r.db.WithContext(ctx).
		Where("level_id = ? AND user_id = ?", levelID, userID).
		First(&jr).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &jr, nil
}

func (r *repository) FindPendingJoinRequests(ctx context.Context, levelID int64) ([]*domain.JoinRequest, error) {
	var requests []*domain.JoinRequest
	if err := r.db.WithContext(ctx).
		Where("level_id = ? AND status = 'pending'", levelID).
		Order("requested_at ASC").
		Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

func (r *repository) UpdateJoinRequest(ctx context.Context, jr *domain.JoinRequest) error {
	return r.db.WithContext(ctx).Save(jr).Error
}
