package academic

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/victorive/klaro-backend/internal/domain"
	"github.com/victorive/klaro-backend/internal/pkg/apperr"
	"github.com/victorive/klaro-backend/internal/pkg/pagination"
)

type Service interface {
	CreateSession(ctx context.Context, req *CreateSessionRequest) (*SessionResponse, error)
	ListSessions(ctx context.Context, q pagination.Query) ([]*SessionResponse, int64, error)
	GetSession(ctx context.Context, id int64) (*SessionResponse, error)
	UpdateSession(ctx context.Context, id int64, req *UpdateSessionRequest) (*SessionResponse, error)
	CreateSemester(ctx context.Context, req *CreateSemesterRequest) (*SemesterResponse, error)
	ListSemestersByLevel(ctx context.Context, levelID int64) ([]*SemesterResponse, error)
	GetSemester(ctx context.Context, id int64) (*SemesterResponse, error)
	UpdateSemester(ctx context.Context, id int64, req *UpdateSemesterRequest) (*SemesterResponse, error)
	DeleteSemester(ctx context.Context, id int64) error
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{repo: repo, logger: logger}
}

func (s *service) CreateSession(ctx context.Context, req *CreateSessionRequest) (*SessionResponse, error) {
	if req.IsCurrent {
		if err := s.repo.UnsetCurrentSession(ctx); err != nil {
			return nil, apperr.New("failed to unset current session", http.StatusInternalServerError).Wrap(err)
		}
	}
	sess := &domain.AcademicSession{Label: req.Label, IsCurrent: req.IsCurrent}
	if err := s.repo.CreateSession(ctx, sess); err != nil {
		return nil, apperr.New("failed to create session", http.StatusInternalServerError).Wrap(err)
	}
	return toSessionResponse(sess), nil
}

func (s *service) ListSessions(ctx context.Context, q pagination.Query) ([]*SessionResponse, int64, error) {
	total, err := s.repo.CountSessions(ctx)
	if err != nil {
		return nil, 0, apperr.New("failed to list sessions", http.StatusInternalServerError).Wrap(err)
	}
	sessions, err := s.repo.FindSessionsPaginated(ctx, q.Offset(), q.Limit())
	if err != nil {
		return nil, 0, apperr.New("failed to list sessions", http.StatusInternalServerError).Wrap(err)
	}
	result := make([]*SessionResponse, 0, len(sessions))
	for _, sess := range sessions {
		result = append(result, toSessionResponse(sess))
	}
	return result, total, nil
}

func (s *service) GetSession(ctx context.Context, id int64) (*SessionResponse, error) {
	sess, err := s.repo.FindSessionByID(ctx, id)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if sess == nil {
		return nil, apperr.New("academic session not found", http.StatusNotFound)
	}
	return toSessionResponse(sess), nil
}

func (s *service) UpdateSession(ctx context.Context, id int64, req *UpdateSessionRequest) (*SessionResponse, error) {
	sess, err := s.repo.FindSessionByID(ctx, id)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if sess == nil {
		return nil, apperr.New("academic session not found", http.StatusNotFound)
	}
	if req.IsCurrent && !sess.IsCurrent {
		if err := s.repo.UnsetCurrentSession(ctx); err != nil {
			return nil, apperr.New("failed to unset current session", http.StatusInternalServerError).Wrap(err)
		}
	}
	sess.Label = req.Label
	sess.IsCurrent = req.IsCurrent
	if err := s.repo.UpdateSession(ctx, sess); err != nil {
		return nil, apperr.New("failed to update session", http.StatusInternalServerError).Wrap(err)
	}
	return toSessionResponse(sess), nil
}

func (s *service) CreateSemester(ctx context.Context, req *CreateSemesterRequest) (*SemesterResponse, error) {
	if req.IsCurrent {
		if err := s.repo.UnsetCurrentSemesterForLevel(ctx, req.LevelID); err != nil {
			return nil, apperr.New("failed to unset current semester", http.StatusInternalServerError).Wrap(err)
		}
	}
	sem := &domain.Semester{
		LevelID:   req.LevelID,
		SessionID: req.SessionID,
		Kind:      req.Kind,
		IsCurrent: req.IsCurrent,
		StartsAt:  req.StartsAt,
		EndsAt:    req.EndsAt,
	}
	if err := s.repo.CreateSemester(ctx, sem); err != nil {
		return nil, apperr.New("failed to create semester", http.StatusInternalServerError).Wrap(err)
	}
	return toSemesterResponse(sem), nil
}

func (s *service) ListSemestersByLevel(ctx context.Context, levelID int64) ([]*SemesterResponse, error) {
	semesters, err := s.repo.FindSemestersByLevelID(ctx, levelID)
	if err != nil {
		return nil, apperr.New("failed to list semesters", http.StatusInternalServerError).Wrap(err)
	}
	result := make([]*SemesterResponse, 0, len(semesters))
	for _, sem := range semesters {
		result = append(result, toSemesterResponse(sem))
	}
	return result, nil
}

func (s *service) GetSemester(ctx context.Context, id int64) (*SemesterResponse, error) {
	sem, err := s.repo.FindSemesterByID(ctx, id)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if sem == nil {
		return nil, apperr.New("semester not found", http.StatusNotFound)
	}
	return toSemesterResponse(sem), nil
}

func (s *service) UpdateSemester(ctx context.Context, id int64, req *UpdateSemesterRequest) (*SemesterResponse, error) {
	sem, err := s.repo.FindSemesterByID(ctx, id)
	if err != nil {
		return nil, apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if sem == nil {
		return nil, apperr.New("semester not found", http.StatusNotFound)
	}
	if req.IsCurrent && !sem.IsCurrent {
		if err := s.repo.UnsetCurrentSemesterForLevel(ctx, sem.LevelID); err != nil {
			return nil, apperr.New("failed to unset current semester", http.StatusInternalServerError).Wrap(err)
		}
	}
	sem.Kind = req.Kind
	sem.IsCurrent = req.IsCurrent
	sem.StartsAt = req.StartsAt
	sem.EndsAt = req.EndsAt
	if err := s.repo.UpdateSemester(ctx, sem); err != nil {
		return nil, apperr.New("failed to update semester", http.StatusInternalServerError).Wrap(err)
	}
	return toSemesterResponse(sem), nil
}

func (s *service) DeleteSemester(ctx context.Context, id int64) error {
	sem, err := s.repo.FindSemesterByID(ctx, id)
	if err != nil {
		return apperr.New("an unexpected error occurred", http.StatusInternalServerError).Wrap(err)
	}
	if sem == nil {
		return apperr.New("semester not found", http.StatusNotFound)
	}
	if err := s.repo.DeleteSemester(ctx, id); err != nil {
		return apperr.New("failed to delete semester", http.StatusInternalServerError).Wrap(err)
	}
	return nil
}
