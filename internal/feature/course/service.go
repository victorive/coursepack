package course

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/victorive/klaro-backend/internal/domain"
	"github.com/victorive/klaro-backend/internal/feature/rbac"
	"github.com/victorive/klaro-backend/internal/pkg/apperr"
)

type Service interface {
	ListCourses(ctx context.Context, universityID int64, departmentID *int64) ([]*CourseRegistryResponse, error)
	GetCourse(ctx context.Context, id, universityID int64) (*CourseRegistryResponse, error)
	AddClassCourse(ctx context.Context, actorUUID string, classID, universityID int64, req AddClassCourseRequest) (*ClassCourseResponse, error)
	ListClassCourses(ctx context.Context, classID int64, semesterID *int64) ([]*ClassCourseResponse, error)
	GetClassCourse(ctx context.Context, id int64) (*ClassCourseResponse, error)
	RemoveClassCourse(ctx context.Context, actorUUID string, classID, classCourseID int64) error
	CreateLecture(ctx context.Context, actorUUID string, classCourseID int64, req CreateLectureRequest) (*LectureResponse, error)
	ListLectures(ctx context.Context, classCourseID int64) ([]*LectureResponse, error)
	GetLecture(ctx context.Context, id int64) (*LectureResponse, error)
	UpdateLecture(ctx context.Context, actorUUID string, id int64, req UpdateLectureRequest) (*LectureResponse, error)
	ArchiveLecture(ctx context.Context, actorUUID string, id int64) error
}

type service struct {
	repo    Repository
	rbacSvc rbac.Service
	logger  *slog.Logger
}

func NewService(repo Repository, rbacSvc rbac.Service, logger *slog.Logger) Service {
	return &service{repo: repo, rbacSvc: rbacSvc, logger: logger}
}

func (s *service) ListCourses(ctx context.Context, universityID int64, departmentID *int64) ([]*CourseRegistryResponse, error) {
	courses, err := s.repo.FindCoursesByUniversity(ctx, universityID, departmentID)
	if err != nil {
		return nil, apperr.New("failed to fetch courses", http.StatusInternalServerError).Wrap(err)
	}
	resp := make([]*CourseRegistryResponse, len(courses))
	for i, c := range courses {
		resp[i] = toCourseRegistryResponse(c)
	}
	return resp, nil
}

func (s *service) GetCourse(ctx context.Context, id, universityID int64) (*CourseRegistryResponse, error) {
	c, err := s.repo.FindCourseByID(ctx, id)
	if err != nil {
		return nil, apperr.New("failed to fetch course", http.StatusInternalServerError).Wrap(err)
	}
	if c == nil || c.UniversityID != universityID {
		return nil, apperr.New("course not found", http.StatusNotFound)
	}
	return toCourseRegistryResponse(c), nil
}

func (s *service) AddClassCourse(ctx context.Context, actorUUID string, classID, universityID int64, req AddClassCourseRequest) (*ClassCourseResponse, error) {
	ok, err := s.rbacSvc.HasClassPermission(ctx, actorUUID, classID, domain.PermCoursesManage)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if !ok {
		return nil, apperr.New("unauthorized", http.StatusForbidden)
	}

	course, err := s.repo.FindCourseByID(ctx, req.CourseID)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if course == nil {
		return nil, apperr.New("course not found", http.StatusNotFound)
	}
	if course.UniversityID != universityID {
		return nil, apperr.New("course does not belong to your university", http.StatusUnprocessableEntity)
	}

	existing, err := s.repo.FindClassCourseByConstraint(ctx, classID, req.SemesterID, req.CourseID)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if existing != nil {
		return nil, apperr.New("course already registered for this semester", http.StatusConflict)
	}

	cc := &domain.ClassCourse{
		ClassID:    classID,
		SemesterID: req.SemesterID,
		CourseID:   req.CourseID,
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

func (s *service) ListClassCourses(ctx context.Context, classID int64, semesterID *int64) ([]*ClassCourseResponse, error) {
	classCourses, err := s.repo.FindClassCoursesByClass(ctx, classID, semesterID)
	if err != nil {
		return nil, apperr.New("failed to fetch class courses", http.StatusInternalServerError).Wrap(err)
	}
	resp := make([]*ClassCourseResponse, len(classCourses))

	for i, cc := range classCourses {
		resp[i] = toClassCourseResponse(cc)
	}

	return resp, nil
}

func (s *service) GetClassCourse(ctx context.Context, id int64) (*ClassCourseResponse, error) {
	cc, err := s.repo.FindClassCourseByID(ctx, id)
	if err != nil {
		return nil, apperr.New("failed to fetch class course", http.StatusInternalServerError).Wrap(err)
	}
	if cc == nil {
		return nil, apperr.New("class course not found", http.StatusNotFound)
	}
	return toClassCourseResponse(cc), nil
}

func (s *service) RemoveClassCourse(ctx context.Context, actorUUID string, classID, classCourseID int64) error {
	ok, err := s.rbacSvc.HasClassPermission(ctx, actorUUID, classID, domain.PermCoursesManage)
	if err != nil {
		return apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if !ok {
		return apperr.New("unauthorized", http.StatusForbidden)
	}
	cc, err := s.repo.FindClassCourseByID(ctx, classCourseID)
	if err != nil {
		return apperr.New("failed to fetch class course", http.StatusInternalServerError).Wrap(err)
	}

	if cc == nil || cc.ClassID != classID {
		return apperr.New("class course not found", http.StatusNotFound)
	}

	return s.repo.DeleteClassCourse(ctx, classCourseID)
}

func (s *service) CreateLecture(ctx context.Context, actorUUID string, classCourseID int64, req CreateLectureRequest) (*LectureResponse, error) {
	cc, err := s.repo.FindClassCourseByID(ctx, classCourseID)
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
		ClassCourseID: classCourseID,
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

func (s *service) ListLectures(ctx context.Context, classCourseID int64) ([]*LectureResponse, error) {
	lectures, err := s.repo.FindLecturesByClassCourse(ctx, classCourseID)
	if err != nil {
		return nil, apperr.New("failed to fetch lectures", http.StatusInternalServerError).Wrap(err)
	}

	resp := make([]*LectureResponse, len(lectures))
	for i, l := range lectures {
		resp[i] = toLectureResponse(l)
	}
	return resp, nil
}

func (s *service) GetLecture(ctx context.Context, id int64) (*LectureResponse, error) {
	l, err := s.repo.FindLectureByID(ctx, id)
	if err != nil {
		return nil, apperr.New("failed to fetch lecture", http.StatusInternalServerError).Wrap(err)
	}
	if l == nil {
		return nil, apperr.New("lecture not found", http.StatusNotFound)
	}
	return toLectureResponse(l), nil
}

func (s *service) UpdateLecture(ctx context.Context, actorUUID string, id int64, req UpdateLectureRequest) (*LectureResponse, error) {
	l, err := s.repo.FindLectureByID(ctx, id)
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

func (s *service) ArchiveLecture(ctx context.Context, actorUUID string, id int64) error {
	l, err := s.repo.FindLectureByID(ctx, id)
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
