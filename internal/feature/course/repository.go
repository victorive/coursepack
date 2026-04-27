package course

import (
	"context"

	"github.com/victorive/klaro-backend/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	FindCoursesByUniversityCursor(ctx context.Context, universityID int64, departmentCode *string, after string, limit int) ([]*domain.Course, error)
	FindCourseByID(ctx context.Context, id int64) (*domain.Course, error)
	FindCourseByUUID(ctx context.Context, uuid string) (*domain.Course, error)
	CreateClassCourse(ctx context.Context, cc *domain.ClassCourse) error
	FindClassCoursesByClass(ctx context.Context, classID int64, semesterID *int64, kind *string) ([]*domain.ClassCourse, error)
	FindClassCourseByID(ctx context.Context, id int64) (*domain.ClassCourse, error)
	FindClassCourseByUUID(ctx context.Context, uuid string) (*domain.ClassCourse, error)
	FindClassCourseByUUIDAndUniversity(ctx context.Context, uuid string, universityID int64) (*domain.ClassCourse, error)
	FindClassCourseByConstraint(ctx context.Context, classID, semesterID, courseID int64) (*domain.ClassCourse, error)
	DeleteClassCourse(ctx context.Context, id int64) error
	CreateLecture(ctx context.Context, l *domain.Lecture) error
	FindLecturesByClassCourse(ctx context.Context, classCourseID int64) ([]*domain.Lecture, error)
	FindLectureByID(ctx context.Context, id int64) (*domain.Lecture, error)
	FindLectureByUUID(ctx context.Context, uuid string) (*domain.Lecture, error)
	FindLectureByUUIDAndUniversity(ctx context.Context, uuid string, universityID int64) (*domain.Lecture, error)
	UpdateLecture(ctx context.Context, l *domain.Lecture) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindCoursesByUniversityCursor(ctx context.Context, universityID int64, departmentCode *string, after string, limit int) ([]*domain.Course, error) {
	var courses []*domain.Course
	q := r.db.WithContext(ctx).
		Joins("JOIN departments d ON d.id = courses.department_id").
		Where("courses.university_id = ?", universityID)
	if departmentCode != nil {
		q = q.Where("d.code = ?", *departmentCode)
	}
	
	if after != "" {
		q = q.Where("courses.code > ?", after)
	}

	if err := q.Order("courses.code ASC").Limit(limit + 1).Find(&courses).Error; err != nil {
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

func (r *repository) FindCourseByUUID(ctx context.Context, uuid string) (*domain.Course, error) {
	var c domain.Course
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&c).Error; err != nil {
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

func (r *repository) FindClassCoursesByClass(ctx context.Context, classID int64, semesterID *int64, kind *string) ([]*domain.ClassCourse, error) {
	var classCourses []*domain.ClassCourse
	q := r.db.WithContext(ctx).Preload("Course").Preload("Semester").Where("class_courses.class_id = ?", classID)

	if semesterID != nil {
		q = q.Where("class_courses.semester_id = ?", *semesterID)
	}

	if kind != nil {
		q = q.Where("class_courses.semester_id IN (?)", 
			r.db.Select("id").Where("kind = ?", *kind).Table("semesters"))
	}

	if err := q.Order("class_courses.created_at ASC").Find(&classCourses).Error; err != nil {
		return nil, err
	}

	return classCourses, nil
}

func (r *repository) FindClassCourseByID(ctx context.Context, id int64) (*domain.ClassCourse, error) {
	var cc domain.ClassCourse
	if err := r.db.WithContext(ctx).Preload("Course").Preload("Semester").First(&cc, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &cc, nil
}

func (r *repository) FindClassCourseByUUID(ctx context.Context, uuid string) (*domain.ClassCourse, error) {
	var cc domain.ClassCourse
	if err := r.db.WithContext(ctx).Preload("Course").Preload("Semester").Where("uuid = ?", uuid).First(&cc).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &cc, nil
}

func (r *repository) FindClassCourseByUUIDAndUniversity(ctx context.Context, uuid string, universityID int64) (*domain.ClassCourse, error) {
	var cc domain.ClassCourse
	err := r.db.WithContext(ctx).
		Preload("Course").
		Preload("Semester").
		Joins("JOIN classes ON classes.id = class_courses.class_id").
		Joins("JOIN departments d ON d.id = classes.department_id").
		Where("class_courses.uuid = ? AND d.university_id = ?", uuid, universityID).
		First(&cc).Error
	if err != nil {
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

func (r *repository) FindLectureByUUID(ctx context.Context, uuid string) (*domain.Lecture, error) {
	var l domain.Lecture
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&l).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &l, nil
}

func (r *repository) FindLectureByUUIDAndUniversity(ctx context.Context, uuid string, universityID int64) (*domain.Lecture, error) {
	var l domain.Lecture
	err := r.db.WithContext(ctx).
		Joins("JOIN class_courses cc ON cc.id = lectures.class_course_id").
		Joins("JOIN classes ON classes.id = cc.class_id").
		Joins("JOIN departments d ON d.id = classes.department_id").
		Where("lectures.uuid = ? AND d.university_id = ?", uuid, universityID).
		First(&l).Error
	if err != nil {
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
