package class

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/victorive/coursepack/internal/domain"
)

type fakeUserRepository struct {
	UserRepository
	user *domain.User
}

func (f *fakeUserRepository) FindByUUID(ctx context.Context, uuid string) (*domain.User, error) {
	return f.user, nil
}

type fakeClassRepository struct {
	Repository
	class *domain.Class
}

func (f *fakeClassRepository) FindByUUID(ctx context.Context, uuid string) (*domain.Class, error) {
	return f.class, nil
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func ptrTo[T any](v T) *T { return &v }

// Get currently takes no caller identity, so a class belonging to another
// university is returned as-is. Flip this test to expect a not-found error once
// the caller's university is checked.
func TestGetReturnsClassFromAnotherUniversity(t *testing.T) {
	caller := &domain.User{ID: 1, UUID: "caller-uuid", UniversityID: ptrTo(int64(1))}
	foreignClass := &domain.Class{
		ID:               7,
		UUID:             "class-7",
		DepartmentID:     3,
		Department:       &domain.Department{ID: 3, UniversityID: 2, Name: "Other", Code: "OTH"},
		AdmissionSession: "2021/2022",
		CurrentLevelID:   4,
		Name:             "CSC 401",
	}

	svc := NewService(
		&fakeClassRepository{class: foreignClass},
		&fakeUserRepository{user: caller},
		nil,
		nil,
		discardLogger(),
	)

	got, err := svc.Get(context.Background(), "class-7")
	if err != nil {
		t.Fatalf("Get returned an error: %v", err)
	}
	if got.UUID != foreignClass.UUID {
		t.Fatalf("got class %q, want %q", got.UUID, foreignClass.UUID)
	}
}
