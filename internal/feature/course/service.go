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
	CreateCourse(ctx context.Context, actorUUID string, req CreateCourseRequest) (*CourseResponse, error)
	ListCoursesByLevel(ctx context.Context, levelID int64) ([]*CourseResponse, error)
	ListCoursesBySemester(ctx context.Context, semesterID int64) ([]*CourseResponse, error)
	GetCourse(ctx context.Context, id int64) (*CourseResponse, error)
	UpdateCourse(ctx context.Context, actorUUID string, id int64, req UpdateCourseRequest) (*CourseResponse, error)
	DeleteCourse(ctx context.Context, actorUUID string, id int64) error

	CreateLecture(ctx context.Context, actorUUID string, req CreateLectureRequest) (*LectureResponse, error)
	ListLectures(ctx context.Context, courseID int64) ([]*LectureResponse, error)
	GetLecture(ctx context.Context, id int64) (*LectureResponse, error)
	UpdateLecture(ctx context.Context, actorUUID string, id int64, req UpdateLectureRequest) (*LectureResponse, error)
	ArchiveLecture(ctx context.Context, actorUUID string, id int64) error
}

type service struct {
	repo     Repository
	rbacSvc  rbac.Service
	logger   *slog.Logger
}

func NewService(repo Repository, rbacSvc rbac.Service, logger *slog.Logger) Service {
	return &service{repo: repo, rbacSvc: rbacSvc, logger: logger}
}

func (s *service) CreateCourse(ctx context.Context, actorUUID string, req CreateCourseRequest) (*CourseResponse, error) {
	ok, err := s.rbacSvc.HasLevelPermission(ctx, actorUUID, req.LevelID, domain.PermCoursesManage)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}

	if !ok {
		return nil, apperr.New("unauthorized", http.StatusForbidden)
	}

	c := &domain.Course{
		LevelID:    req.LevelID,
		SemesterID: req.SemesterID,
		Code:       req.Code,
		Title:      req.Title,
		Units:      req.Units,
	}
	if err := s.repo.CreateCourse(ctx, c); err != nil {
		return nil, apperr.New("failed to create course", http.StatusInternalServerError).Wrap(err)
	}
	return toCourseResponse(c), nil
}

func (s *service) ListCoursesByLevel(ctx context.Context, levelID int64) ([]*CourseResponse, error) {
	courses, err := s.repo.FindCoursesByLevel(ctx, levelID)
	if err != nil {
		return nil, apperr.New("failed to fetch courses", http.StatusInternalServerError).Wrap(err)
	}
	resp := make([]*CourseResponse, len(courses))
	for i, c := range courses {
		resp[i] = toCourseResponse(c)
	}
	return resp, nil
}

func (s *service) ListCoursesBySemester(ctx context.Context, semesterID int64) ([]*CourseResponse, error) {
	courses, err := s.repo.FindCoursesBySemester(ctx, semesterID)
	if err != nil {
		return nil, apperr.New("failed to fetch courses", http.StatusInternalServerError).Wrap(err)
	}
	resp := make([]*CourseResponse, len(courses))
	for i, c := range courses {
		resp[i] = toCourseResponse(c)
	}
	return resp, nil
}

func (s *service) GetCourse(ctx context.Context, id int64) (*CourseResponse, error) {
	c, err := s.repo.FindCourseByID(ctx, id)
	if err != nil {
		return nil, apperr.New("failed to fetch course", http.StatusInternalServerError).Wrap(err)
	}
	if c == nil {
		return nil, apperr.New("course not found", http.StatusNotFound)
	}
	return toCourseResponse(c), nil
}

func (s *service) UpdateCourse(ctx context.Context, actorUUID string, id int64, req UpdateCourseRequest) (*CourseResponse, error) {
	c, err := s.repo.FindCourseByID(ctx, id)
	if err != nil {
		return nil, apperr.New("failed to fetch course", http.StatusInternalServerError).Wrap(err)
	}

	if c == nil {
		return nil, apperr.New("course not found", http.StatusNotFound)
	}

	ok, err := s.rbacSvc.HasLevelPermission(ctx, actorUUID, c.LevelID, domain.PermCoursesManage)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}

	if !ok {
		return nil, apperr.New("unauthorized", http.StatusForbidden)
	}

	if req.Code != "" {
		c.Code = req.Code
	}

	if req.Title != "" {
		c.Title = req.Title
	}

	if req.Units != 0 {
		c.Units = req.Units
	}

	if err := s.repo.UpdateCourse(ctx, c); err != nil {
		return nil, apperr.New("failed to update course", http.StatusInternalServerError).Wrap(err)
	}

	return toCourseResponse(c), nil
}

func (s *service) DeleteCourse(ctx context.Context, actorUUID string, id int64) error {
	c, err := s.repo.FindCourseByID(ctx, id)
	if err != nil {
		return apperr.New("failed to fetch course", http.StatusInternalServerError).Wrap(err)
	}

	if c == nil {
		return apperr.New("course not found", http.StatusNotFound)
	}

	ok, err := s.rbacSvc.HasLevelPermission(ctx, actorUUID, c.LevelID, domain.PermCoursesManage)
	if err != nil {
		return apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}

	if !ok {
		return apperr.New("unauthorized", http.StatusForbidden)
	}
	return s.repo.DeleteCourse(ctx, id)
}

func (s *service) CreateLecture(ctx context.Context, actorUUID string, req CreateLectureRequest) (*LectureResponse, error) {
	c, err := s.repo.FindCourseByID(ctx, req.CourseID)
	if err != nil {
		return nil, apperr.New("failed to fetch course", http.StatusInternalServerError).Wrap(err)
	}

	if c == nil {
		return nil, apperr.New("course not found", http.StatusNotFound)
	}

	ok, err := s.rbacSvc.HasLevelPermission(ctx, actorUUID, c.LevelID, domain.PermLecturesManage)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}

	if !ok {
		return nil, apperr.New("unauthorized", http.StatusForbidden)
	}

	l := &domain.Lecture{
		CourseID:     req.CourseID,
		Topic:        req.Topic,
		Description:  req.Description,
		LecturerName: req.LecturerName,
		WeekNumber:   req.WeekNumber,
		Kind:         req.Kind,
		Status:       string(domain.LectureStatusPublished),
	}

	if err := s.repo.CreateLecture(ctx, l); err != nil {
		return nil, apperr.New("failed to create lecture", http.StatusInternalServerError).Wrap(err)
	}

	return toLectureResponse(l), nil
}

func (s *service) ListLectures(ctx context.Context, courseID int64) ([]*LectureResponse, error) {
	lectures, err := s.repo.FindLecturesByCourse(ctx, courseID)
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

	c, err := s.repo.FindCourseByID(ctx, l.CourseID)
	if err != nil {
		return nil, apperr.New("failed to fetch course", http.StatusInternalServerError).Wrap(err)
	}
	
	if c == nil {
		return nil, apperr.New("course not found", http.StatusNotFound)
	}

	ok, err := s.rbacSvc.HasLevelPermission(ctx, actorUUID, c.LevelID, domain.PermLecturesManage)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}

	if !ok {
		return nil, apperr.New("unauthorized", http.StatusForbidden)
	}
	
	if req.Topic != "" {
		l.Topic = req.Topic
	}

	if req.Description != nil {
		l.Description = req.Description
	}

	if req.LecturerName != nil {
		l.LecturerName = req.LecturerName
	}

	if req.WeekNumber != nil {
		l.WeekNumber = req.WeekNumber
	}

	if req.Kind != "" {
		l.Kind = req.Kind
	}

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
	c, err := s.repo.FindCourseByID(ctx, l.CourseID)
	if err != nil {
		return apperr.New("failed to fetch course", http.StatusInternalServerError).Wrap(err)
	}

	if c == nil {
		return apperr.New("course not found", http.StatusNotFound)
	}

	ok, err := s.rbacSvc.HasLevelPermission(ctx, actorUUID, c.LevelID, domain.PermLecturesManage)
	if err != nil {
		return apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}

	if !ok {
		return apperr.New("unauthorized", http.StatusForbidden)
	}

	l.Status = string(domain.LectureStatusArchived)
	
	return s.repo.UpdateLecture(ctx, l)
}
