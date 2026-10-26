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
	"github.com/victorive/coursepack/internal/domain"
	"github.com/victorive/coursepack/internal/permissions"
	"github.com/victorive/coursepack/internal/pkg/apperr"
	"github.com/victorive/coursepack/internal/storage"
)

type Service interface {
	Upload(ctx context.Context, uploaderUUID string, req UploadFileRequest, reader io.Reader, size int64, mimeType string, originalName string) (*FileResponse, error)
	ListByClassCourse(ctx context.Context, classCourseUUID string) ([]*FileResponse, error)
	ListByLecture(ctx context.Context, lectureUUID string) ([]*FileResponse, error)
	Get(ctx context.Context, actorUUID string, fileUUID string) (*FileResponse, error)
	GetDownloadURL(ctx context.Context, actorUUID string, fileUUID string) (string, error)
	Review(ctx context.Context, fileUUID string, reviewerUUID string, req ReviewFileRequest) (*FileResponse, error)
	Archive(ctx context.Context, fileUUID string, actorUUID string) error
}

type UserRepository interface {
	FindByUUID(ctx context.Context, uuid string) (*domain.User, error)
}

type ClassCourseRepository interface {
	FindClassCourseByID(ctx context.Context, id int64) (*domain.ClassCourse, error)
	FindClassCourseByUUID(ctx context.Context, uuid string) (*domain.ClassCourse, error)
}

type LectureRepository interface {
	FindLectureByUUID(ctx context.Context, uuid string) (*domain.Lecture, error)
}

type ClassRepository interface {
	IsMember(ctx context.Context, classID, userID int64) (bool, error)
}

type PointRepository interface {
	AwardPoints(ctx context.Context, userID int64, delta int, sourceType domain.PointSource, sourceID int64, awardedBy int64) error
}

type service struct {
	repo            Repository
	userRepo        UserRepository
	classCourseRepo ClassCourseRepository
	lectureRepo     LectureRepository
	classRepo       ClassRepository
	pointRepo       PointRepository
	permSvc         permissions.Service
	storage         storage.Storage
	maxFileSize     int64
	logger          *slog.Logger
}

func NewService(repo Repository, userRepo UserRepository, classCourseRepo ClassCourseRepository, lectureRepo LectureRepository, classRepo ClassRepository, permSvc permissions.Service, stor storage.Storage, maxFileSize int64, pointRepo PointRepository, logger *slog.Logger) Service {
	return &service{repo: repo, userRepo: userRepo, classCourseRepo: classCourseRepo, lectureRepo: lectureRepo, classRepo: classRepo, pointRepo: pointRepo, permSvc: permSvc, storage: stor, maxFileSize: maxFileSize, logger: logger}
}

func (s *service) resolveUserID(ctx context.Context, uuid string) (int64, error) {
	u, err := s.userRepo.FindByUUID(ctx, uuid)
	if err != nil {
		return 0, apperr.Internal(err)
	}
	if u == nil {
		return 0, apperr.NotFound("user")
	}
	return u.ID, nil
}

// authorizeRead allows the uploader, an approver of the file's class, or a
// member of that class when the file is approved. Everything else is reported
// as not found so file UUIDs cannot be probed.
func (s *service) authorizeRead(ctx context.Context, actorUUID string, f *domain.File) error {
	user, err := s.userRepo.FindByUUID(ctx, actorUUID)
	if err != nil {
		return apperr.Internal(err)
	}
	if user == nil {
		return apperr.NotFound("user")
	}
	if f.UploadedBy == user.ID {
		return nil
	}

	classCourse, err := s.classCourseRepo.FindClassCourseByID(ctx, f.ClassCourseID)
	if err != nil {
		return apperr.Internal(err)
	}
	if classCourse == nil {
		return apperr.NotFound("file")
	}

	if f.Status != string(domain.FileStatusApproved) {
		canReview, err := s.permSvc.HasClassPermission(ctx, actorUUID, classCourse.ClassID, domain.PermFilesApprove)
		if err != nil {
			return apperr.Internal(err)
		}
		if !canReview {
			return apperr.NotFound("file")
		}
		return nil
	}

	isMember, err := s.classRepo.IsMember(ctx, classCourse.ClassID, user.ID)
	if err != nil {
		return apperr.Internal(err)
	}
	if !isMember {
		return apperr.NotFound("file")
	}
	return nil
}

func (s *service) findReadableFile(ctx context.Context, actorUUID string, fileUUID string) (*domain.File, error) {
	f, err := s.repo.FindByUUID(ctx, fileUUID)
	if err != nil {
		return nil, apperr.New("failed to fetch file", http.StatusInternalServerError).Wrap(err)
	}
	if f == nil {
		return nil, apperr.NotFound("file")
	}
	if err := s.authorizeRead(ctx, actorUUID, f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *service) Upload(ctx context.Context, uploaderUUID string, req UploadFileRequest, reader io.Reader, size int64, mimeType string, originalName string) (*FileResponse, error) {
	uploaderID, err := s.resolveUserID(ctx, uploaderUUID)
	if err != nil {
		return nil, err
	}
	classCourse, err := s.classCourseRepo.FindClassCourseByUUID(ctx, req.ClassCourseUUID)
	if err != nil {
		return nil, apperr.New("failed to fetch class course", http.StatusInternalServerError).Wrap(err)
	}
	if classCourse == nil {
		return nil, apperr.NotFound("class course")
	}

	var lectureID *int64
	if req.LectureUUID != nil {
		l, err := s.lectureRepo.FindLectureByUUID(ctx, *req.LectureUUID)
		if err != nil {
			return nil, apperr.New("failed to fetch lecture", http.StatusInternalServerError).Wrap(err)
		}

		if l == nil {
			return nil, apperr.NotFound("lecture")
		}

		lectureID = &l.ID
	}

	ok, err := s.permSvc.HasClassPermission(ctx, uploaderUUID, classCourse.ClassID, domain.PermFilesUpload)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if !ok {
		return nil, apperr.Forbidden()
	}
	if size > s.maxFileSize {
		return nil, apperr.New(fmt.Sprintf("file exceeds maximum allowed size of %d bytes", s.maxFileSize), http.StatusRequestEntityTooLarge)
	}

	key := fmt.Sprintf("files/%s/%s", uuid.NewString(), filepath.Base(originalName))
	if _, err := s.storage.Upload(ctx, key, reader, size, mimeType); err != nil {
		return nil, apperr.New("failed to upload file", http.StatusInternalServerError).Wrap(err)
	}

	f := &domain.File{
		UUID:          uuid.NewString(),
		UploadedBy:    uploaderID,
		ClassCourseID: classCourse.ID,
		LectureID:     lectureID,
		Category:      string(req.Category),
		Title:         req.Title,
		StorageKey:    key,
		SizeBytes:     size,
		MimeType:      mimeType,
		YearLabel:     req.YearLabel,
		Status:        string(domain.FileStatusPending),
	}
	if err := s.repo.Create(ctx, f); err != nil {
		_ = s.storage.Delete(ctx, key)
		return nil, apperr.New("failed to save file record", http.StatusInternalServerError).Wrap(err)
	}
	return toFileResponse(f), nil
}

func (s *service) ListByClassCourse(ctx context.Context, classCourseUUID string) ([]*FileResponse, error) {
	cc, err := s.classCourseRepo.FindClassCourseByUUID(ctx, classCourseUUID)
	if err != nil {
		return nil, apperr.New("failed to fetch class course", http.StatusInternalServerError).Wrap(err)
	}

	if cc == nil {
		return nil, apperr.NotFound("class course")
	}

	files, err := s.repo.FindByClassCourse(ctx, cc.ID)
	if err != nil {
		return nil, apperr.New("failed to fetch files", http.StatusInternalServerError).Wrap(err)
	}

	resp := make([]*FileResponse, len(files))
	for i, f := range files {
		resp[i] = toFileResponse(f)
	}

	return resp, nil
}

func (s *service) ListByLecture(ctx context.Context, lectureUUID string) ([]*FileResponse, error) {
	l, err := s.lectureRepo.FindLectureByUUID(ctx, lectureUUID)
	if err != nil {
		return nil, apperr.New("failed to fetch lecture", http.StatusInternalServerError).Wrap(err)
	}

	if l == nil {
		return nil, apperr.NotFound("lecture")
	}

	files, err := s.repo.FindByLecture(ctx, l.ID)
	if err != nil {
		return nil, apperr.New("failed to fetch files", http.StatusInternalServerError).Wrap(err)
	}

	resp := make([]*FileResponse, len(files))
	for i, f := range files {
		resp[i] = toFileResponse(f)
	}

	return resp, nil
}

func (s *service) Get(ctx context.Context, actorUUID string, fileUUID string) (*FileResponse, error) {
	f, err := s.findReadableFile(ctx, actorUUID, fileUUID)
	if err != nil {
		return nil, err
	}
	return toFileResponse(f), nil
}

func (s *service) GetDownloadURL(ctx context.Context, actorUUID string, fileUUID string) (string, error) {
	f, err := s.findReadableFile(ctx, actorUUID, fileUUID)
	if err != nil {
		return "", err
	}
	url, err := s.storage.GetSignedURL(ctx, f.StorageKey, 15*time.Minute)
	if err != nil {
		return "", apperr.New("failed to generate download URL", http.StatusInternalServerError).Wrap(err)
	}
	return url, nil
}

func (s *service) Review(ctx context.Context, fileUUID string, reviewerUUID string, req ReviewFileRequest) (*FileResponse, error) {
	reviewerID, err := s.resolveUserID(ctx, reviewerUUID)
	if err != nil {
		return nil, err
	}
	f, err := s.repo.FindByUUID(ctx, fileUUID)
	if err != nil {
		return nil, apperr.New("failed to fetch file", http.StatusInternalServerError).Wrap(err)
	}
	if f == nil {
		return nil, apperr.NotFound("file")
	}
	classCourse, err := s.classCourseRepo.FindClassCourseByID(ctx, f.ClassCourseID)
	if err != nil {
		return nil, apperr.New("failed to fetch class course", http.StatusInternalServerError).Wrap(err)
	}
	if classCourse == nil {
		return nil, apperr.NotFound("class course")
	}
	ok, err := s.permSvc.HasClassPermission(ctx, reviewerUUID, classCourse.ClassID, domain.PermFilesApprove)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if !ok {
		return nil, apperr.Forbidden()
	}

	now := time.Now()
	f.ReviewedBy = &reviewerID
	f.ReviewedAt = &now
	f.ReviewNote = req.ReviewNote

	if req.Action == domain.ReviewActionApprove {
		f.Status = string(domain.FileStatusApproved)
		if err := s.repo.Update(ctx, f); err != nil {
			return nil, apperr.New("failed to update file", http.StatusInternalServerError).Wrap(err)
		}

		if err := s.pointRepo.AwardPoints(ctx, f.UploadedBy, domain.PointsFileApproval, domain.PointSourceFileApproval, f.ID, reviewerID); err != nil {
			s.logger.Warn("failed to award points for file approval", "file_id", f.ID, "error", err)
		}
	} else {
		f.Status = string(domain.FileStatusRejected)
		if err := s.repo.Update(ctx, f); err != nil {
			return nil, apperr.New("failed to update file", http.StatusInternalServerError).Wrap(err)
		}
	}

	return toFileResponse(f), nil
}

func (s *service) Archive(ctx context.Context, fileUUID string, actorUUID string) error {
	f, err := s.repo.FindByUUID(ctx, fileUUID)
	if err != nil {
		return apperr.New("failed to fetch file", http.StatusInternalServerError).Wrap(err)
	}
	if f == nil {
		return apperr.NotFound("file")
	}

	classCourse, err := s.classCourseRepo.FindClassCourseByID(ctx, f.ClassCourseID)
	if err != nil {
		return apperr.New("failed to fetch class course", http.StatusInternalServerError).Wrap(err)
	}
	if classCourse == nil {
		return apperr.NotFound("class course")
	}

	ok, err := s.permSvc.HasClassPermission(ctx, actorUUID, classCourse.ClassID, domain.PermFilesApprove)
	if err != nil {
		return apperr.Internal(err)
	}
	if !ok {
		return apperr.Forbidden()
	}

	f.Status = string(domain.FileStatusArchived)
	return s.repo.Update(ctx, f)
}
