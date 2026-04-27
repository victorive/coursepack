package course

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/victorive/klaro-backend/internal/domain"
	"github.com/victorive/klaro-backend/internal/feature/rbac"
	"github.com/victorive/klaro-backend/internal/pkg/apperr"
	"github.com/victorive/klaro-backend/internal/pkg/pagination"
)

type Service interface {
	ListCourses(ctx context.Context, universityID int64, departmentCode *string, q pagination.CursorQuery) ([]*CourseRegistryResponse, *string, error)
	GetCourse(ctx context.Context, courseUUID string, universityID int64) (*CourseRegistryResponse, error)
	AddClassCourse(ctx context.Context, actorUUID string, classUUID string, universityID int64, req AddClassCourseRequest) (*ClassCourseResponse, error)
	ListClassCourses(ctx context.Context, classUUID string, universityID int64, semesterID *int64) ([]*ClassCourseResponse, error)
	GetClassCourse(ctx context.Context, uuid string, universityID int64) (*ClassCourseResponse, error)
	RemoveClassCourse(ctx context.Context, actorUUID string, classUUID string, classCourseUUID string) error
	CreateLecture(ctx context.Context, actorUUID string, classCourseUUID string, req CreateLectureRequest) (*LectureResponse, error)
	ListLectures(ctx context.Context, classCourseUUID string, universityID int64) ([]*LectureResponse, error)
	GetLecture(ctx context.Context, uuid string, universityID int64) (*LectureResponse, error)
	UpdateLecture(ctx context.Context, actorUUID string, uuid string, req UpdateLectureRequest) (*LectureResponse, error)
	ArchiveLecture(ctx context.Context, actorUUID string, uuid string) error
}

type ClassRepository interface {
	FindByUUID(ctx context.Context, uuid string) (*domain.Class, error)
	FindByUUIDAndUniversity(ctx context.Context, uuid string, universityID int64) (*domain.Class, error)
}

type SemesterRepository interface {
	FindSemesterByID(ctx context.Context, id int64) (*domain.Semester, error)
}

type service struct {
	repo      Repository
	classRepo ClassRepository
	semRepo   SemesterRepository
	rbacSvc   rbac.Service
	logger    *slog.Logger
}

func NewService(repo Repository, classRepo ClassRepository, semRepo SemesterRepository, rbacSvc rbac.Service, logger *slog.Logger) Service {
	return &service{repo: repo, classRepo: classRepo, semRepo: semRepo, rbacSvc: rbacSvc, logger: logger}
}

func (s *service) ListCourses(ctx context.Context, universityID int64, departmentCode *string, q pagination.CursorQuery) ([]*CourseRegistryResponse, *string, error) {
	courses, err := s.repo.FindCoursesByUniversityCursor(ctx, universityID, departmentCode, q.After, q.Limit)
	if err != nil {
		return nil, nil, apperr.New("failed to fetch courses", http.StatusInternalServerError).Wrap(err)
	}

	var nextCursor *string
	if len(courses) > q.Limit {
		courses = courses[:q.Limit]
		c := courses[q.Limit-1].Code
		nextCursor = &c
	}
	resp := make([]*CourseRegistryResponse, len(courses))
	for i, c := range courses {
		resp[i] = toCourseRegistryResponse(c)
	}
	return resp, nextCursor, nil
}

func (s *service) GetCourse(ctx context.Context, courseUUID string, universityID int64) (*CourseRegistryResponse, error) {
	c, err := s.repo.FindCourseByUUID(ctx, courseUUID)
	if err != nil {
		return nil, apperr.New("failed to fetch course", http.StatusInternalServerError).Wrap(err)
	}
	if c == nil || c.UniversityID != universityID {
		return nil, apperr.New("course not found", http.StatusNotFound)
	}
	return toCourseRegistryResponse(c), nil
}

func (s *service) AddClassCourse(ctx context.Context, actorUUID string, classUUID string, universityID int64, req AddClassCourseRequest) (*ClassCourseResponse, error) {
	ok, err := s.rbacSvc.HasClassPermissionByUUID(ctx, actorUUID, classUUID, domain.PermCoursesManage)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if !ok {
		return nil, apperr.New("unauthorized", http.StatusForbidden)
	}

	cls, err := s.classRepo.FindByUUID(ctx, classUUID)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}

	if cls == nil {
		return nil, apperr.New("class not found", http.StatusNotFound)
	}

	sem, err := s.semRepo.FindSemesterByID(ctx, req.SemesterID)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if sem == nil || sem.ClassID != cls.ID {
		return nil, apperr.New("semester not found for this class", http.StatusUnprocessableEntity)
	}

	course, err := s.repo.FindCourseByUUID(ctx, req.CourseUUID)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if course == nil {
		return nil, apperr.New("course not found", http.StatusNotFound)
	}
	if course.UniversityID != universityID {
		return nil, apperr.New("course does not belong to your university", http.StatusUnprocessableEntity)
	}

	existing, err := s.repo.FindClassCourseByConstraint(ctx, cls.ID, req.SemesterID, course.ID)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if existing != nil {
		return nil, apperr.New("course already registered for this semester", http.StatusConflict)
	}

	cc := &domain.ClassCourse{
		UUID:       uuid.NewString(),
		ClassID:    cls.ID,
		SemesterID: req.SemesterID,
		CourseID:   course.ID,
	}
	if err := s.repo.CreateClassCourse(ctx, cc); err != nil {
		return nil, apperr.New("failed to add course to class", http.StatusInternalServerError).Wrap(err)
	}

	cc, err = s.repo.FindClassCourseByID(ctx, cc.ID)
	if err != nil || cc == nil {
		return nil, apperr.New("failed to fetch class course", http.StatusInternalServerError).Wrap(err)
	}

	return toClassCourseResponse(cc), nil
}

func (s *service) ListClassCourses(ctx context.Context, classUUID string, universityID int64, semesterID *int64) ([]*ClassCourseResponse, error) {
	cls, err := s.classRepo.FindByUUIDAndUniversity(ctx, classUUID, universityID)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}

	if cls == nil {
		return nil, apperr.New("class not found", http.StatusNotFound)
	}

	classCourses, err := s.repo.FindClassCoursesByClass(ctx, cls.ID, semesterID)
	if err != nil {
		return nil, apperr.New("failed to fetch class courses", http.StatusInternalServerError).Wrap(err)
	}

	resp := make([]*ClassCourseResponse, len(classCourses))

	for i, cc := range classCourses {
		resp[i] = toClassCourseResponse(cc)
	}

	return resp, nil
}

func (s *service) GetClassCourse(ctx context.Context, classCourseUUID string, universityID int64) (*ClassCourseResponse, error) {
	cc, err := s.repo.FindClassCourseByUUIDAndUniversity(ctx, classCourseUUID, universityID)
	if err != nil {
		return nil, apperr.New("failed to fetch class course", http.StatusInternalServerError).Wrap(err)
	}
	if cc == nil {
		return nil, apperr.New("class course not found", http.StatusNotFound)
	}
	return toClassCourseResponse(cc), nil
}

func (s *service) RemoveClassCourse(ctx context.Context, actorUUID string, classUUID string, classCourseUUID string) error {
	ok, err := s.rbacSvc.HasClassPermissionByUUID(ctx, actorUUID, classUUID, domain.PermCoursesManage)
	if err != nil {
		return apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if !ok {
		return apperr.New("unauthorized", http.StatusForbidden)
	}
	cls, err := s.classRepo.FindByUUID(ctx, classUUID)
	if err != nil {
		return apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}

	if cls == nil {
		return apperr.New("class not found", http.StatusNotFound)
	}

	cc, err := s.repo.FindClassCourseByUUID(ctx, classCourseUUID)
	if err != nil {
		return apperr.New("failed to fetch class course", http.StatusInternalServerError).Wrap(err)
	}

	if cc == nil || cc.ClassID != cls.ID {
		return apperr.New("class course not found", http.StatusNotFound)
	}

	return s.repo.DeleteClassCourse(ctx, cc.ID)
}

func (s *service) CreateLecture(ctx context.Context, actorUUID string, classCourseUUID string, req CreateLectureRequest) (*LectureResponse, error) {
	cc, err := s.repo.FindClassCourseByUUID(ctx, classCourseUUID)
	if err != nil {
		return nil, apperr.New("failed to fetch class course", http.StatusInternalServerError).Wrap(err)
	}
	if cc == nil {
		return nil, apperr.New("class course not found", http.StatusNotFound)
	}

	ok, err := s.rbacSvc.HasClassPermission(ctx, actorUUID, cc.ClassID, domain.PermLecturesManage)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if !ok {
		return nil, apperr.New("unauthorized", http.StatusForbidden)
	}

	l := &domain.Lecture{
		UUID:          uuid.NewString(),
		ClassCourseID: cc.ID,
		Topic:         req.Topic,
		Description:   req.Description,
		LecturerName:  req.LecturerName,
		WeekNumber:    req.WeekNumber,
		Kind:          req.Kind,
		Status:        string(domain.LectureStatusPublished),
	}

	if err := s.repo.CreateLecture(ctx, l); err != nil {
		return nil, apperr.New("failed to create lecture", http.StatusInternalServerError).Wrap(err)
	}

	return toLectureResponse(l), nil
}

func (s *service) ListLectures(ctx context.Context, classCourseUUID string, universityID int64) ([]*LectureResponse, error) {
	cc, err := s.repo.FindClassCourseByUUIDAndUniversity(ctx, classCourseUUID, universityID)
	if err != nil {
		return nil, apperr.New("failed to fetch class course", http.StatusInternalServerError).Wrap(err)
	}

	if cc == nil {
		return nil, apperr.New("class course not found", http.StatusNotFound)
	}

	lectures, err := s.repo.FindLecturesByClassCourse(ctx, cc.ID)
	if err != nil {
		return nil, apperr.New("failed to fetch lectures", http.StatusInternalServerError).Wrap(err)
	}

	resp := make([]*LectureResponse, len(lectures))
	for i, l := range lectures {
		resp[i] = toLectureResponse(l)
	}
	return resp, nil
}

func (s *service) GetLecture(ctx context.Context, lectureUUID string, universityID int64) (*LectureResponse, error) {
	l, err := s.repo.FindLectureByUUIDAndUniversity(ctx, lectureUUID, universityID)
	if err != nil {
		return nil, apperr.New("failed to fetch lecture", http.StatusInternalServerError).Wrap(err)
	}
	if l == nil {
		return nil, apperr.New("lecture not found", http.StatusNotFound)
	}
	return toLectureResponse(l), nil
}

func (s *service) UpdateLecture(ctx context.Context, actorUUID string, lectureUUID string, req UpdateLectureRequest) (*LectureResponse, error) {
	l, err := s.repo.FindLectureByUUID(ctx, lectureUUID)
	if err != nil {
		return nil, apperr.New("failed to fetch lecture", http.StatusInternalServerError).Wrap(err)
	}
	if l == nil {
		return nil, apperr.New("lecture not found", http.StatusNotFound)
	}

	cc, err := s.repo.FindClassCourseByID(ctx, l.ClassCourseID)
	if err != nil || cc == nil {
		return nil, apperr.New("failed to fetch class course", http.StatusInternalServerError).Wrap(err)
	}

	ok, err := s.rbacSvc.HasClassPermission(ctx, actorUUID, cc.ClassID, domain.PermLecturesManage)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if !ok {
		return nil, apperr.New("unauthorized", http.StatusForbidden)
	}

	if req.Topic != "" {
		l.Topic = req.Topic
	}

	if req.Kind != "" {
		l.Kind = req.Kind
	}
	l.Description = req.Description
	l.LecturerName = req.LecturerName
	l.WeekNumber = req.WeekNumber

	if err := s.repo.UpdateLecture(ctx, l); err != nil {
		return nil, apperr.New("failed to update lecture", http.StatusInternalServerError).Wrap(err)
	}
	return toLectureResponse(l), nil
}

func (s *service) ArchiveLecture(ctx context.Context, actorUUID string, lectureUUID string) error {
	l, err := s.repo.FindLectureByUUID(ctx, lectureUUID)
	if err != nil {
		return apperr.New("failed to fetch lecture", http.StatusInternalServerError).Wrap(err)
	}
	if l == nil {
		return apperr.New("lecture not found", http.StatusNotFound)
	}

	cc, err := s.repo.FindClassCourseByID(ctx, l.ClassCourseID)
	if err != nil || cc == nil {
		return apperr.New("failed to fetch class course", http.StatusInternalServerError).Wrap(err)
	}

	ok, err := s.rbacSvc.HasClassPermission(ctx, actorUUID, cc.ClassID, domain.PermLecturesManage)
	if err != nil {
		return apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if !ok {
		return apperr.New("unauthorized", http.StatusForbidden)
	}

	l.Status = string(domain.LectureStatusArchived)
	return s.repo.UpdateLecture(ctx, l)
}
