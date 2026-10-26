package file

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/victorive/coursepack/internal/domain"
)

type fakeFileRepository struct {
	Repository
	file *domain.File
}

func (f *fakeFileRepository) FindByUUID(ctx context.Context, uuid string) (*domain.File, error) {
	return f.file, nil
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// Get currently ignores the caller, so any authenticated user can read a file
// they have no access to. Flip this test to expect a not-found error once the
// caller's class membership is checked.
func TestGetReturnsFileToNonMember(t *testing.T) {
	stored := &domain.File{
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

	svc := NewService(
		&fakeFileRepository{file: stored},
		nil,
		nil,
		nil,
		nil,
		nil,
		0,
		nil,
		discardLogger(),
	)

	got, err := svc.Get(context.Background(), "file-11")
	if err != nil {
		t.Fatalf("Get returned an error: %v", err)
	}
	if got.UUID != stored.UUID {
		t.Fatalf("got file %q, want %q", got.UUID, stored.UUID)
	}
}
