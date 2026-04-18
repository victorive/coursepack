package academic

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/pkg/pagination"
	"github.com/victorive/klaro-backend/internal/pkg/response"
	"github.com/victorive/klaro-backend/internal/pkg/validator"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// ---- Sessions ----

// CreateSession godoc
// @Summary      Create academic session
// @Tags         Academic
// @Accept       json
// @Produce      json
// @Param        body body CreateSessionRequest true "Create session payload"
// @Success      201 {object} response.Response{data=SessionResponse}
// @Security     BearerAuth
// @Router       /sessions [post]
func (h *Handler) CreateSession(c *fiber.Ctx) error {
	var req CreateSessionRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	sess, err := h.svc.CreateSession(c.UserContext(), &req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusCreated, true, "academic session created", sess)
}

// ListSessions godoc
// @Summary      List academic sessions
// @Tags         Academic
// @Produce      json
// @Success      200 {object} response.Response{data=object}
// @Security     BearerAuth
// @Router       /sessions [get]
func (h *Handler) ListSessions(c *fiber.Ctx) error {
	q := pagination.Parse(c)
	items, total, err := h.svc.ListSessions(c.UserContext(), *q)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "sessions retrieved", pagination.NewPaginated(c, items, total, q))
}

// GetSession godoc
// @Summary      Get academic session
// @Tags         Academic
// @Produce      json
// @Param        id path int true "Session ID"
// @Success      200 {object} response.Response{data=SessionResponse}
// @Security     BearerAuth
// @Router       /sessions/{id} [get]
func (h *Handler) GetSession(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid id", nil)
	}
	sess, err := h.svc.GetSession(c.UserContext(), id)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "session retrieved", sess)
}

// UpdateSession godoc
// @Summary      Update academic session
// @Tags         Academic
// @Accept       json
// @Produce      json
// @Param        id   path int                  true "Session ID"
// @Param        body body UpdateSessionRequest true "Update session payload"
// @Success      200 {object} response.Response{data=SessionResponse}
// @Security     BearerAuth
// @Router       /sessions/{id} [patch]
func (h *Handler) UpdateSession(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid id", nil)
	}
	var req UpdateSessionRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	sess, err := h.svc.UpdateSession(c.UserContext(), id, &req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "session updated", sess)
}

// ---- Semesters ----

// CreateSemester godoc
// @Summary      Create semester
// @Tags         Academic
// @Accept       json
// @Produce      json
// @Param        body body CreateSemesterRequest true "Create semester payload"
// @Success      201 {object} response.Response{data=SemesterResponse}
// @Security     BearerAuth
// @Router       /semesters [post]
func (h *Handler) CreateSemester(c *fiber.Ctx) error {
	var req CreateSemesterRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	sem, err := h.svc.CreateSemester(c.UserContext(), &req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusCreated, true, "semester created", sem)
}

// ListSemestersByCohort godoc
// @Summary      List semesters for a cohort
// @Tags         Academic
// @Produce      json
// @Param        cohortId path int true "Cohort ID"
// @Success      200 {object} response.Response{data=[]SemesterResponse}
// @Security     BearerAuth
// @Router       /cohorts/{cohortId}/semesters [get]
func (h *Handler) ListSemestersByCohort(c *fiber.Ctx) error {
	cohortID, err := strconv.ParseInt(c.Params("cohortId"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid cohort id", nil)
	}
	semesters, err := h.svc.ListSemestersByCohort(c.UserContext(), cohortID)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "semesters retrieved", semesters)
}

// GetSemester godoc
// @Summary      Get semester
// @Tags         Academic
// @Produce      json
// @Param        id path int true "Semester ID"
// @Success      200 {object} response.Response{data=SemesterResponse}
// @Security     BearerAuth
// @Router       /semesters/{id} [get]
func (h *Handler) GetSemester(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid id", nil)
	}
	sem, err := h.svc.GetSemester(c.UserContext(), id)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "semester retrieved", sem)
}

// UpdateSemester godoc
// @Summary      Update semester
// @Tags         Academic
// @Accept       json
// @Produce      json
// @Param        id   path int                    true "Semester ID"
// @Param        body body UpdateSemesterRequest  true "Update semester payload"
// @Success      200 {object} response.Response{data=SemesterResponse}
// @Security     BearerAuth
// @Router       /semesters/{id} [patch]
func (h *Handler) UpdateSemester(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid id", nil)
	}
	var req UpdateSemesterRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	sem, err := h.svc.UpdateSemester(c.UserContext(), id, &req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "semester updated", sem)
}

// DeleteSemester godoc
// @Summary      Delete semester
// @Tags         Academic
// @Produce      json
// @Param        id path int true "Semester ID"
// @Success      204
// @Security     BearerAuth
// @Router       /semesters/{id} [delete]
func (h *Handler) DeleteSemester(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid id", nil)
	}
	if err := h.svc.DeleteSemester(c.UserContext(), id); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}
