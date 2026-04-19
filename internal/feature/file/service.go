package file

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/victorive/klaro-backend/internal/domain"
	"github.com/victorive/klaro-backend/internal/feature/rbac"
	"github.com/victorive/klaro-backend/internal/pkg/apperr"
	"github.com/victorive/klaro-backend/internal/storage"
)

type Service interface {
	Upload(ctx context.Context, uploaderUUID string, req UploadFileRequest, reader io.Reader, size int64, mimeType string, originalName string) (*FileResponse, error)
	ListByCourse(ctx context.Context, courseID int64) ([]*FileResponse, error)
	ListByLecture(ctx context.Context, lectureID int64) ([]*FileResponse, error)
	Get(ctx context.Context, id int64) (*FileResponse, error)
	GetDownloadURL(ctx context.Context, id int64) (string, error)
	Review(ctx context.Context, id int64, reviewerUUID string, req ReviewFileRequest) (*FileResponse, error)
	Archive(ctx context.Context, id int64, actorUUID string) error
}

type UserRepository interface {
	FindByUUID(ctx context.Context, uuid string) (*domain.User, error)
}

type CourseRepository interface {
	FindCourseByID(ctx context.Context, id int64) (*domain.Course, error)
}

type service struct {
	repo        Repository
	userRepo    UserRepository
	courseRepo  CourseRepository
	rbacSvc     rbac.Service
	storage     storage.Storage
	maxFileSize int64
	logger      *slog.Logger
}

func NewService(repo Repository, userRepo UserRepository, courseRepo CourseRepository, rbacSvc rbac.Service, stor storage.Storage, maxFileSize int64, logger *slog.Logger) Service {
	return &service{repo: repo, userRepo: userRepo, courseRepo: courseRepo, rbacSvc: rbacSvc, storage: stor, maxFileSize: maxFileSize, logger: logger}
}

func (s *service) resolveUserID(ctx context.Context, uuid string) (int64, error) {
	u, err := s.userRepo.FindByUUID(ctx, uuid)
	if err != nil {
		return 0, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if u == nil {
		return 0, apperr.New("user not found", http.StatusNotFound)
	}
	return u.ID, nil
}

func (s *service) Upload(ctx context.Context, uploaderUUID string, req UploadFileRequest, reader io.Reader, size int64, mimeType string, originalName string) (*FileResponse, error) {
	uploaderID, err := s.resolveUserID(ctx, uploaderUUID)
	if err != nil {
		return nil, err
	}
	course, err := s.courseRepo.FindCourseByID(ctx, req.CourseID)
	if err != nil {
		return nil, apperr.New("failed to fetch course", http.StatusInternalServerError).Wrap(err)
	}
	if course == nil {
		return nil, apperr.New("course not found", http.StatusNotFound)
	}
	ok, err := s.rbacSvc.HasClassPermission(ctx, uploaderUUID, course.ClassID, domain.PermFilesUpload)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if !ok {
		return nil, apperr.New("unauthorized", http.StatusForbidden)
	}
	if size > s.maxFileSize {
		return nil, apperr.New(fmt.Sprintf("file exceeds maximum allowed size of %d bytes", s.maxFileSize), http.StatusRequestEntityTooLarge)
	}

	key := fmt.Sprintf("files/%s/%s", uuid.NewString(), filepath.Base(originalName))
	if _, err := s.storage.Upload(ctx, key, reader, size, mimeType); err != nil {
		return nil, apperr.New("failed to upload file", http.StatusInternalServerError).Wrap(err)
	}

	f := &domain.File{
		UploadedBy: uploaderID,
		CourseID:   req.CourseID,
		LectureID:  req.LectureID,
		Category:   req.Category,
		Title:      req.Title,
		StorageKey: key,
		SizeBytes:  size,
		MimeType:   mimeType,
		YearLabel:  req.YearLabel,
		Status:     string(domain.FileStatusPending),
	}
	if err := s.repo.Create(ctx, f); err != nil {
		_ = s.storage.Delete(ctx, key)
		return nil, apperr.New("failed to save file record", http.StatusInternalServerError).Wrap(err)
	}
	return toFileResponse(f), nil
}

func (s *service) ListByCourse(ctx context.Context, courseID int64) ([]*FileResponse, error) {
	files, err := s.repo.FindByCourse(ctx, courseID)
	if err != nil {
		return nil, apperr.New("failed to fetch files", http.StatusInternalServerError).Wrap(err)
	}
	resp := make([]*FileResponse, len(files))
	for i, f := range files {
		resp[i] = toFileResponse(f)
	}
	return resp, nil
}

func (s *service) ListByLecture(ctx context.Context, lectureID int64) ([]*FileResponse, error) {
	files, err := s.repo.FindByLecture(ctx, lectureID)
	if err != nil {
		return nil, apperr.New("failed to fetch files", http.StatusInternalServerError).Wrap(err)
	}
	resp := make([]*FileResponse, len(files))
	for i, f := range files {
		resp[i] = toFileResponse(f)
	}
	return resp, nil
}

func (s *service) Get(ctx context.Context, id int64) (*FileResponse, error) {
	f, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperr.New("failed to fetch file", http.StatusInternalServerError).Wrap(err)
	}
	if f == nil {
		return nil, apperr.New("file not found", http.StatusNotFound)
	}
	return toFileResponse(f), nil
}

func (s *service) GetDownloadURL(ctx context.Context, id int64) (string, error) {
	f, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return "", apperr.New("failed to fetch file", http.StatusInternalServerError).Wrap(err)
	}
	if f == nil {
		return "", apperr.New("file not found", http.StatusNotFound)
	}
	url, err := s.storage.GetSignedURL(ctx, f.StorageKey, 15*time.Minute)
	if err != nil {
		return "", apperr.New("failed to generate download URL", http.StatusInternalServerError).Wrap(err)
	}
	return url, nil
}

func (s *service) Review(ctx context.Context, id int64, reviewerUUID string, req ReviewFileRequest) (*FileResponse, error) {
	reviewerID, err := s.resolveUserID(ctx, reviewerUUID)
	if err != nil {
		return nil, err
	}
	f, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperr.New("failed to fetch file", http.StatusInternalServerError).Wrap(err)
	}
	if f == nil {
		return nil, apperr.New("file not found", http.StatusNotFound)
	}
	course, err := s.courseRepo.FindCourseByID(ctx, f.CourseID)
	if err != nil {
		return nil, apperr.New("failed to fetch course", http.StatusInternalServerError).Wrap(err)
	}

	if course == nil {
		return nil, apperr.New("course not found", http.StatusNotFound)
	}
	ok, err := s.rbacSvc.HasClassPermission(ctx, reviewerUUID, course.ClassID, domain.PermFilesApprove)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}

	if !ok {
		return nil, apperr.New("unauthorized", http.StatusForbidden)
	}

	now := time.Now()
	f.ReviewedBy = &reviewerID
	f.ReviewedAt = &now
	f.ReviewNote = req.ReviewNote

	if req.Action == "approve" {
		f.Status = string(domain.FileStatusApproved)
	} else {
		f.Status = string(domain.FileStatusRejected)
	}
	if err := s.repo.Update(ctx, f); err != nil {
		return nil, apperr.New("failed to update file", http.StatusInternalServerError).Wrap(err)
	}

	return toFileResponse(f), nil
}

func (s *service) Archive(ctx context.Context, id int64, actorUUID string) error {
	f, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return apperr.New("failed to fetch file", http.StatusInternalServerError).Wrap(err)
	}

	if f == nil {
		return apperr.New("file not found", http.StatusNotFound)
	}

	course, err := s.courseRepo.FindCourseByID(ctx, f.CourseID)
	if err != nil {
		return apperr.New("failed to fetch course", http.StatusInternalServerError).Wrap(err)
	}

	if course == nil {
		return apperr.New("course not found", http.StatusNotFound)
	}

	ok, err := s.rbacSvc.HasClassPermission(ctx, actorUUID, course.ClassID, domain.PermFilesApprove)
	if err != nil {
		return apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}

	if !ok {
		return apperr.New("unauthorized", http.StatusForbidden)
	}

	f.Status = string(domain.FileStatusArchived)

	return s.repo.Update(ctx, f)
}
