package course

import (
	"context"

	"github.com/victorive/klaro-backend/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	CreateCourse(ctx context.Context, c *domain.Course) error
	FindCoursesByClass(ctx context.Context, classID int64) ([]*domain.Course, error)
	FindCoursesBySemester(ctx context.Context, semesterID int64) ([]*domain.Course, error)
	FindCourseByID(ctx context.Context, id int64) (*domain.Course, error)
	UpdateCourse(ctx context.Context, c *domain.Course) error
	DeleteCourse(ctx context.Context, id int64) error
	CreateLecture(ctx context.Context, l *domain.Lecture) error
	FindLecturesByCourse(ctx context.Context, courseID int64) ([]*domain.Lecture, error)
	FindLectureByID(ctx context.Context, id int64) (*domain.Lecture, error)
	UpdateLecture(ctx context.Context, l *domain.Lecture) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateCourse(ctx context.Context, c *domain.Course) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *repository) FindCoursesByClass(ctx context.Context, classID int64) ([]*domain.Course, error) {
	var courses []*domain.Course
	if err := r.db.WithContext(ctx).Where("class_id = ?", classID).Order("code ASC").Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *repository) FindCoursesBySemester(ctx context.Context, semesterID int64) ([]*domain.Course, error) {
	var courses []*domain.Course
	if err := r.db.WithContext(ctx).Where("semester_id = ?", semesterID).Order("code ASC").Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *repository) FindCourseByID(ctx context.Context, id int64) (*domain.Course, error) {
	var c domain.Course
	if err := r.db.WithContext(ctx).First(&c, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *repository) UpdateCourse(ctx context.Context, c *domain.Course) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *repository) DeleteCourse(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Course{}, id).Error
}

func (r *repository) CreateLecture(ctx context.Context, l *domain.Lecture) error {
	return r.db.WithContext(ctx).Create(l).Error
}

func (r *repository) FindLecturesByCourse(ctx context.Context, courseID int64) ([]*domain.Lecture, error) {
	var lectures []*domain.Lecture
	if err := r.db.WithContext(ctx).
		Where("course_id = ? AND status = 'published'", courseID).
		Order("week_number ASC, created_at ASC").
		Find(&lectures).Error; err != nil {
		return nil, err
	}
	return lectures, nil
}

func (r *repository) FindLectureByID(ctx context.Context, id int64) (*domain.Lecture, error) {
	var l domain.Lecture
	if err := r.db.WithContext(ctx).First(&l, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &l, nil
}

func (r *repository) UpdateLecture(ctx context.Context, l *domain.Lecture) error {
	return r.db.WithContext(ctx).Save(l).Error
}
