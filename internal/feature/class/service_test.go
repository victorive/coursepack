package class

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"testing"

	"github.com/victorive/coursepack/internal/domain"
	"github.com/victorive/coursepack/internal/feature/rbac"
	"github.com/victorive/coursepack/internal/pkg/apperr"
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

func (f *fakeClassRepository) FindByUUIDAndUniversity(ctx context.Context, uuid string, universityID int64) (*domain.Class, error) {
	if f.class == nil || f.class.Department == nil {
		return nil, nil
	}
	if f.class.Department.UniversityID != universityID {
		return nil, nil
	}
	return f.class, nil
}

func (f *fakeClassRepository) Transact(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func (f *fakeClassRepository) Create(ctx context.Context, c *domain.Class) error {
	c.ID = 1
	return nil
}

func (f *fakeClassRepository) AddMember(ctx context.Context, m *domain.ClassMember) error {
	return nil
}

type fakeDepartmentRepository struct {
	DepartmentRepository
	department *domain.Department
}

func (f *fakeDepartmentRepository) FindByID(ctx context.Context, id int64) (*domain.Department, error) {
	return f.department, nil
}

type fakeRBACService struct {
	rbac.Service
	assignErr error
}

func (f *fakeRBACService) AssignClassRole(ctx context.Context, userID, classID int64, role domain.RoleName) error {
	return f.assignErr
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func ptrTo[T any](v T) *T { return &v }

func newServiceUnderTest(caller *domain.User, class *domain.Class) Service {
	return NewService(
		&fakeClassRepository{class: class},
		&fakeUserRepository{user: caller},
		nil,
		nil,
		discardLogger(),
	)
}

func TestGetHidesClassFromAnotherUniversity(t *testing.T) {
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

	_, err := newServiceUnderTest(caller, foreignClass).Get(context.Background(), caller.UUID, foreignClass.UUID)
	if err == nil {
		t.Fatal("expected a not-found error for a class in another university")
	}

	appErr, ok := err.(*apperr.AppError)
	if !ok {
		t.Fatalf("got %T, want *apperr.AppError", err)
	}
	if appErr.HTTPStatus != http.StatusNotFound {
		t.Fatalf("got status %d, want %d", appErr.HTTPStatus, http.StatusNotFound)
	}
}

func TestGetReturnsClassFromCallerUniversity(t *testing.T) {
	caller := &domain.User{ID: 1, UUID: "caller-uuid", UniversityID: ptrTo(int64(1))}
	ownClass := &domain.Class{
		ID:               8,
		UUID:             "class-8",
		DepartmentID:     4,
		Department:       &domain.Department{ID: 4, UniversityID: 1, Name: "Own", Code: "OWN"},
		AdmissionSession: "2021/2022",
		CurrentLevelID:   3,
		Name:             "CSC 301",
	}

	got, err := newServiceUnderTest(caller, ownClass).Get(context.Background(), caller.UUID, ownClass.UUID)
	if err != nil {
		t.Fatalf("Get returned an error: %v", err)
	}
	if got.UUID != ownClass.UUID {
		t.Fatalf("got class %q, want %q", got.UUID, ownClass.UUID)
	}
}

func TestCreateReturnsErrorWhenRoleAssignmentFails(t *testing.T) {
	caller := &domain.User{ID: 1, UUID: "caller-uuid", UniversityID: ptrTo(int64(1))}
	department := &domain.Department{ID: 3, UniversityID: 1}

	svc := NewService(
		&fakeClassRepository{},
		&fakeUserRepository{user: caller},
		&fakeDepartmentRepository{department: department},
		&fakeRBACService{assignErr: errors.New("role store down")},
		discardLogger(),
	)

	_, err := svc.Create(context.Background(), caller.UUID, &CreateClassRequest{
		DepartmentID:     3,
		AdmissionSession: "2021/2022",
		CurrentLevelID:   1,
		Name:             "CSC 401",
	})
	if err == nil {
		t.Fatal("expected Create to fail when role assignment fails")
	}
}
