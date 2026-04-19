package academic

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/pkg/response"
	"github.com/victorive/klaro-backend/internal/pkg/validator"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
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

// ListSemestersByClass godoc
// @Summary      List semesters for a class
// @Tags         Academic
// @Produce      json
// @Param        classId path int true "Class ID"
// @Success      200 {object} response.Response{data=[]SemesterResponse}
// @Security     BearerAuth
// @Router       /classes/{classId}/semesters [get]
func (h *Handler) ListSemestersByClass(c *fiber.Ctx) error {
	classID, err := strconv.ParseInt(c.Params("classId"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid class id", nil)
	}
	semesters, err := h.svc.ListSemestersByClass(c.UserContext(), classID)
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
