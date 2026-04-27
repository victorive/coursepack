package class

import (
	"context"

	"github.com/victorive/klaro-backend/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, c *domain.Class) error
	Update(ctx context.Context, c *domain.Class) error
	UpdateInviteCode(ctx context.Context, classID int64, code string) error
	FindByID(ctx context.Context, id int64) (*domain.Class, error)
	FindByUUID(ctx context.Context, uuid string) (*domain.Class, error)
	FindByUUIDAndUniversity(ctx context.Context, uuid string, universityID int64) (*domain.Class, error)
	FindByInviteCode(ctx context.Context, code string) (*domain.Class, error)
	Delete(ctx context.Context, id int64) error
	AddMember(ctx context.Context, m *domain.ClassMember) error
	RemoveMember(ctx context.Context, classID, userID int64) error
	FindMembers(ctx context.Context, classID int64) ([]*domain.ClassMember, error)
	IsMember(ctx context.Context, classID, userID int64) (bool, error)
	IsActiveMemberAnywhere(ctx context.Context, userID, universityID int64) (bool, error)
	CreateJoinRequest(ctx context.Context, jr *domain.JoinRequest) error
	FindJoinRequestByID(ctx context.Context, id int64) (*domain.JoinRequest, error)
	FindJoinRequestByClassAndUser(ctx context.Context, classID, userID int64) (*domain.JoinRequest, error)
	FindPendingJoinRequests(ctx context.Context, classID int64) ([]*domain.JoinRequest, error)
	UpdateJoinRequest(ctx context.Context, jr *domain.JoinRequest) error
	FindAllLevels(ctx context.Context) ([]*domain.Level, error)
	CountByUniversity(ctx context.Context, universityID int64) (int64, error)
	FindPaginatedByUniversity(ctx context.Context, universityID int64, offset, limit int) ([]*domain.Class, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, c *domain.Class) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *repository) Update(ctx context.Context, c *domain.Class) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *repository) FindByID(ctx context.Context, id int64) (*domain.Class, error) {
	var c domain.Class
	if err := r.db.WithContext(ctx).First(&c, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *repository) FindByUUID(ctx context.Context, uuid string) (*domain.Class, error) {
	var c domain.Class
	if err := r.db.WithContext(ctx).Preload("Department").Where("uuid = ?", uuid).First(&c).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}

func (r *repository) FindByUUIDAndUniversity(ctx context.Context, uuid string, universityID int64) (*domain.Class, error) {
	var c domain.Class
	err := r.db.WithContext(ctx).
		Joins("JOIN departments d ON d.id = classes.department_id").
		Where("classes.uuid = ? AND d.university_id = ?", uuid, universityID).
		First(&c).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}

func (r *repository) FindByInviteCode(ctx context.Context, code string) (*domain.Class, error) {
	var c domain.Class
	if err := r.db.WithContext(ctx).Where("invite_code = ?", code).First(&c).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Class{}, id).Error
}

func (r *repository) AddMember(ctx context.Context, m *domain.ClassMember) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *repository) RemoveMember(ctx context.Context, classID, userID int64) error {
	return r.db.WithContext(ctx).
		Where("class_id = ? AND user_id = ?", classID, userID).
		Delete(&domain.ClassMember{}).Error
}

func (r *repository) FindMembers(ctx context.Context, classID int64) ([]*domain.ClassMember, error) {
	var members []*domain.ClassMember
	if err := r.db.WithContext(ctx).Preload("User").Where("class_id = ?", classID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *repository) IsMember(ctx context.Context, classID, userID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.ClassMember{}).
		Where("class_id = ? AND user_id = ?", classID, userID).Count(&count).Error
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

func (r *repository) FindJoinRequestByClassAndUser(ctx context.Context, classID, userID int64) (*domain.JoinRequest, error) {
	var jr domain.JoinRequest
	if err := r.db.WithContext(ctx).
		Where("class_id = ? AND user_id = ?", classID, userID).
		First(&jr).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &jr, nil
}

func (r *repository) FindPendingJoinRequests(ctx context.Context, classID int64) ([]*domain.JoinRequest, error) {
	var requests []*domain.JoinRequest
	if err := r.db.WithContext(ctx).
		Where("class_id = ? AND status = 'pending'", classID).
		Order("requested_at ASC").
		Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

func (r *repository) UpdateJoinRequest(ctx context.Context, jr *domain.JoinRequest) error {
	return r.db.WithContext(ctx).Save(jr).Error
}

func (r *repository) FindAllLevels(ctx context.Context) ([]*domain.Level, error) {
	var levels []*domain.Level
	if err := r.db.WithContext(ctx).Order("value ASC").Find(&levels).Error; err != nil {
		return nil, err
	}
	return levels, nil
}

func (r *repository) CountByUniversity(ctx context.Context, universityID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Class{}).
		Joins("JOIN departments ON departments.id = classes.department_id").
		Where("departments.university_id = ?", universityID).
		Count(&count).Error

	return count, err
}

func (r *repository) FindPaginatedByUniversity(ctx context.Context, universityID int64, offset, limit int) ([]*domain.Class, error) {
	var classes []*domain.Class
	err := r.db.WithContext(ctx).
		Preload("Department").
		Joins("JOIN departments ON departments.id = classes.department_id").
		Where("departments.university_id = ?", universityID).
		Order("classes.id ASC").
		Offset(offset).Limit(limit).
		Find(&classes).Error

	return classes, err
}

func (r *repository) UpdateInviteCode(ctx context.Context, classID int64, code string) error {
	return r.db.WithContext(ctx).
		Model(&domain.Class{}).
		Where("id = ?", classID).
		Update("invite_code", code).Error
}

func (r *repository) IsActiveMemberAnywhere(ctx context.Context, userID, universityID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.ClassMember{}).
		Joins("JOIN classes ON classes.id = class_members.class_id").
		Joins("JOIN departments ON departments.id = classes.department_id").
		Where("class_members.user_id = ? AND departments.university_id = ?", userID, universityID).
		Count(&count).Error
	return count > 0, err
}
