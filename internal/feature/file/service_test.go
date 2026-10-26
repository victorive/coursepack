package file

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"testing"

	"github.com/victorive/coursepack/internal/domain"
	"github.com/victorive/coursepack/internal/feature/rbac"
	"github.com/victorive/coursepack/internal/pkg/apperr"
)

type fakeFileRepository struct {
	Repository
	file *domain.File
}

func (f *fakeFileRepository) FindByUUID(ctx context.Context, uuid string) (*domain.File, error) {
	return f.file, nil
}

type fakeUserRepository struct {
	UserRepository
	user *domain.User
}

func (f *fakeUserRepository) FindByUUID(ctx context.Context, uuid string) (*domain.User, error) {
	return f.user, nil
}

type fakeClassCourseRepository struct {
	ClassCourseRepository
	classCourse *domain.ClassCourse
}

func (f *fakeClassCourseRepository) FindClassCourseByID(ctx context.Context, id int64) (*domain.ClassCourse, error) {
	return f.classCourse, nil
}

type fakeClassRepository struct {
	ClassRepository
	member bool
}

func (f *fakeClassRepository) IsMember(ctx context.Context, classID, userID int64) (bool, error) {
	return f.member, nil
}

type fakeRBACService struct {
	rbac.Service
	canApprove bool
}

func (f *fakeRBACService) HasClassPermission(ctx context.Context, userUUID string, classID int64, perm domain.PermissionName) (bool, error) {
	return f.canApprove, nil
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func newServiceUnderTest(file *domain.File, user *domain.User, member, canApprove bool) Service {
	return NewService(
		&fakeFileRepository{file: file},
		&fakeUserRepository{user: user},
		&fakeClassCourseRepository{classCourse: &domain.ClassCourse{ID: 5, ClassID: 42}},
		nil,
		&fakeClassRepository{member: member},
		&fakeRBACService{canApprove: canApprove},
		nil,
		0,
		nil,
		discardLogger(),
	)
}

func approvedFile() *domain.File {
	return &domain.File{
		ID:            11,
		UUID:          "file-11",
		UploadedBy:    99,
		ClassCourseID: 5,
		Category:      string(domain.FileCategoryLectureNote),
		Title:         "Week 1",
		StorageKey:    "files/file-11.pdf",
		SizeBytes:     1024,
		MimeType:      "application/pdf",
		Status:        string(domain.FileStatusApproved),
	}
}

func assertNotFound(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected a not-found error")
	}
	appErr, ok := err.(*apperr.AppError)
	if !ok {
		t.Fatalf("got %T, want *apperr.AppError", err)
	}
	if appErr.HTTPStatus != http.StatusNotFound {
		t.Fatalf("got status %d, want %d", appErr.HTTPStatus, http.StatusNotFound)
	}
}

func TestGetHidesApprovedFileFromNonMember(t *testing.T) {
	user := &domain.User{ID: 1, UUID: "outsider-uuid"}

	_, err := newServiceUnderTest(approvedFile(), user, false, false).Get(context.Background(), user.UUID, "file-11")
	assertNotFound(t, err)
}

func TestGetReturnsApprovedFileToMember(t *testing.T) {
	user := &domain.User{ID: 1, UUID: "member-uuid"}

	got, err := newServiceUnderTest(approvedFile(), user, true, false).Get(context.Background(), user.UUID, "file-11")
	if err != nil {
		t.Fatalf("Get returned an error: %v", err)
	}
	if got.UUID != "file-11" {
		t.Fatalf("got file %q, want %q", got.UUID, "file-11")
	}
}

func TestGetHidesPendingFileFromMember(t *testing.T) {
	user := &domain.User{ID: 1, UUID: "member-uuid"}
	pending := approvedFile()
	pending.Status = string(domain.FileStatusPending)

	_, err := newServiceUnderTest(pending, user, true, false).Get(context.Background(), user.UUID, pending.UUID)
	assertNotFound(t, err)
}

func TestGetReturnsPendingFileToUploader(t *testing.T) {
	user := &domain.User{ID: 99, UUID: "uploader-uuid"}
	pending := approvedFile()
	pending.Status = string(domain.FileStatusPending)

	got, err := newServiceUnderTest(pending, user, false, false).Get(context.Background(), user.UUID, pending.UUID)
	if err != nil {
		t.Fatalf("Get returned an error: %v", err)
	}
	if got.UUID != pending.UUID {
		t.Fatalf("got file %q, want %q", got.UUID, pending.UUID)
	}
}

func TestGetReturnsPendingFileToApprover(t *testing.T) {
	user := &domain.User{ID: 1, UUID: "rep-uuid"}
	pending := approvedFile()
	pending.Status = string(domain.FileStatusPending)

	got, err := newServiceUnderTest(pending, user, false, true).Get(context.Background(), user.UUID, pending.UUID)
	if err != nil {
		t.Fatalf("Get returned an error: %v", err)
	}
	if got.UUID != pending.UUID {
		t.Fatalf("got file %q, want %q", got.UUID, pending.UUID)
	}
}
