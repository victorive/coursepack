package cohort

import (
	"context"

	"github.com/victorive/klaro-backend/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, c *domain.Cohort) error
	Count(ctx context.Context) (int64, error)
	FindPaginated(ctx context.Context, offset, limit int) ([]*domain.Cohort, error)
	FindByID(ctx context.Context, id int64) (*domain.Cohort, error)
	FindByInviteCode(ctx context.Context, code string) (*domain.Cohort, error)
	Delete(ctx context.Context, id int64) error
	AddMember(ctx context.Context, m *domain.CohortMember) error
	RemoveMember(ctx context.Context, cohortID, userID int64) error
	FindMembers(ctx context.Context, cohortID int64) ([]*domain.CohortMember, error)
	IsMember(ctx context.Context, cohortID, userID int64) (bool, error)
	CreateJoinRequest(ctx context.Context, jr *domain.JoinRequest) error
	FindJoinRequestByID(ctx context.Context, id int64) (*domain.JoinRequest, error)
	FindJoinRequestByCohortAndUser(ctx context.Context, cohortID, userID int64) (*domain.JoinRequest, error)
	FindPendingJoinRequests(ctx context.Context, cohortID int64) ([]*domain.JoinRequest, error)
	UpdateJoinRequest(ctx context.Context, jr *domain.JoinRequest) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, c *domain.Cohort) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *repository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Cohort{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) FindPaginated(ctx context.Context, offset, limit int) ([]*domain.Cohort, error) {
	var cohorts []*domain.Cohort
	if err := r.db.WithContext(ctx).Order("level ASC").Offset(offset).Limit(limit).Find(&cohorts).Error; err != nil {
		return nil, err
	}
	return cohorts, nil
}

func (r *repository) FindByID(ctx context.Context, id int64) (*domain.Cohort, error) {
	var c domain.Cohort
	if err := r.db.WithContext(ctx).First(&c, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *repository) FindByInviteCode(ctx context.Context, code string) (*domain.Cohort, error) {
	var c domain.Cohort
	if err := r.db.WithContext(ctx).Where("invite_code = ?", code).First(&c).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Cohort{}, id).Error
}

func (r *repository) AddMember(ctx context.Context, m *domain.CohortMember) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *repository) RemoveMember(ctx context.Context, cohortID, userID int64) error {
	return r.db.WithContext(ctx).
		Where("cohort_id = ? AND user_id = ?", cohortID, userID).
		Delete(&domain.CohortMember{}).Error
}

func (r *repository) FindMembers(ctx context.Context, cohortID int64) ([]*domain.CohortMember, error) {
	var members []*domain.CohortMember
	if err := r.db.WithContext(ctx).Where("cohort_id = ?", cohortID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *repository) IsMember(ctx context.Context, cohortID, userID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.CohortMember{}).
		Where("cohort_id = ? AND user_id = ?", cohortID, userID).Count(&count).Error
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

func (r *repository) FindJoinRequestByCohortAndUser(ctx context.Context, cohortID, userID int64) (*domain.JoinRequest, error) {
	var jr domain.JoinRequest
	if err := r.db.WithContext(ctx).
		Where("cohort_id = ? AND user_id = ?", cohortID, userID).
		First(&jr).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &jr, nil
}

func (r *repository) FindPendingJoinRequests(ctx context.Context, cohortID int64) ([]*domain.JoinRequest, error) {
	var requests []*domain.JoinRequest
	if err := r.db.WithContext(ctx).
		Where("cohort_id = ? AND status = 'pending'", cohortID).
		Order("requested_at ASC").
		Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

func (r *repository) UpdateJoinRequest(ctx context.Context, jr *domain.JoinRequest) error {
	return r.db.WithContext(ctx).Save(jr).Error
}
