package course

import (
	"context"

	"github.com/victorive/klaro-backend/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	FindCoursesByUniversity(ctx context.Context, universityID int64, departmentID *int64) ([]*domain.Course, error)
	FindCourseByID(ctx context.Context, id int64) (*domain.Course, error)
	CreateClassCourse(ctx context.Context, cc *domain.ClassCourse) error
	FindClassCoursesByClass(ctx context.Context, classID int64, semesterID *int64) ([]*domain.ClassCourse, error)
	FindClassCourseByID(ctx context.Context, id int64) (*domain.ClassCourse, error)
	FindClassCourseByConstraint(ctx context.Context, classID, semesterID, courseID int64) (*domain.ClassCourse, error)
	DeleteClassCourse(ctx context.Context, id int64) error
	CreateLecture(ctx context.Context, l *domain.Lecture) error
	FindLecturesByClassCourse(ctx context.Context, classCourseID int64) ([]*domain.Lecture, error)
	FindLectureByID(ctx context.Context, id int64) (*domain.Lecture, error)
	UpdateLecture(ctx context.Context, l *domain.Lecture) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindCoursesByUniversity(ctx context.Context, universityID int64, departmentID *int64) ([]*domain.Course, error) {
	var courses []*domain.Course
	q := r.db.WithContext(ctx).Where("university_id = ?", universityID)
	if departmentID != nil {
		q = q.Where("department_id = ?", *departmentID)
	}
	if err := q.Order("code ASC").Find(&courses).Error; err != nil {
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

func (r *repository) CreateClassCourse(ctx context.Context, cc *domain.ClassCourse) error {
	return r.db.WithContext(ctx).Create(cc).Error
}

func (r *repository) FindClassCoursesByClass(ctx context.Context, classID int64, semesterID *int64) ([]*domain.ClassCourse, error) {
	var classCourses []*domain.ClassCourse
	q := r.db.WithContext(ctx).Preload("Course").Where("class_id = ?", classID)
	if semesterID != nil {
		q = q.Where("semester_id = ?", *semesterID)
	}

	if err := q.Order("created_at ASC").Find(&classCourses).Error; err != nil {
		return nil, err
	}

	return classCourses, nil
}

func (r *repository) FindClassCourseByID(ctx context.Context, id int64) (*domain.ClassCourse, error) {
	var cc domain.ClassCourse
	if err := r.db.WithContext(ctx).Preload("Course").First(&cc, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &cc, nil
}

func (r *repository) FindClassCourseByConstraint(ctx context.Context, classID, semesterID, courseID int64) (*domain.ClassCourse, error) {
	var cc domain.ClassCourse
	err := r.db.WithContext(ctx).
		Where("class_id = ? AND semester_id = ? AND course_id = ?", classID, semesterID, courseID).
		First(&cc).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &cc, nil
}

func (r *repository) DeleteClassCourse(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.ClassCourse{}, id).Error
}

func (r *repository) CreateLecture(ctx context.Context, l *domain.Lecture) error {
	return r.db.WithContext(ctx).Create(l).Error
}

func (r *repository) FindLecturesByClassCourse(ctx context.Context, classCourseID int64) ([]*domain.Lecture, error) {
	var lectures []*domain.Lecture
	if err := r.db.WithContext(ctx).
		Where("class_course_id = ? AND status = 'published'", classCourseID).
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
