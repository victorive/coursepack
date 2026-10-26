package semester

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/coursepack/internal/pkg/response"
	"github.com/victorive/coursepack/internal/pkg/validator"
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
// @Router       /classes/{classId}/semesters [post]
func (h *Handler) CreateSemester(c *fiber.Ctx) error {
	classUUID := c.Params("classId")
	var req CreateSemesterRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	sem, err := h.svc.CreateSemester(c.UserContext(), classUUID, &req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusCreated, true, "semester created", sem)
}

// ListSemestersByClass godoc
// @Summary      List semesters for a class
// @Tags         Academic
// @Produce      json
// @Param        id path string true "Class UUID"
// @Success      200 {object} response.Response{data=[]SemesterResponse}
// @Security     BearerAuth
// @Router       /classes/{classId}/semesters [get]
func (h *Handler) ListSemestersByClass(c *fiber.Ctx) error {
	classUUID := c.Params("classId")
	semesters, err := h.svc.ListSemestersByClass(c.UserContext(), classUUID)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "semesters retrieved", semesters)
}

// GetSemester godoc
// @Summary      Get semester
// @Tags         Academic
// @Produce      json
// @Param        classId path string true "Class UUID"
// @Param        semId   path string true "Semester UUID"
// @Success      200 {object} response.Response{data=SemesterResponse}
// @Security     BearerAuth
// @Router       /classes/{classId}/semesters/{semId} [get]
func (h *Handler) GetSemester(c *fiber.Ctx) error {
	semUUID := c.Params("semId")
	sem, err := h.svc.GetSemester(c.UserContext(), semUUID)
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
// @Param        classId path string               true "Class UUID"
// @Param        semId   path string                true "Semester UUID"
// @Param        body    body UpdateSemesterRequest true "Update semester payload"
// @Success      200 {object} response.Response{data=SemesterResponse}
// @Security     BearerAuth
// @Router       /classes/{classId}/semesters/{semId} [patch]
func (h *Handler) UpdateSemester(c *fiber.Ctx) error {
	semUUID := c.Params("semId")
	var req UpdateSemesterRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	sem, err := h.svc.UpdateSemester(c.UserContext(), semUUID, &req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "semester updated", sem)
}

// DeleteSemester godoc
// @Summary      Delete semester
// @Tags         Academic
// @Produce      json
// @Param        classId path string true "Class UUID"
// @Param        semId   path string true "Semester UUID"
// @Success      204
// @Security     BearerAuth
// @Router       /classes/{classId}/semesters/{semId} [delete]
func (h *Handler) DeleteSemester(c *fiber.Ctx) error {
	semUUID := c.Params("semId")
	if err := h.svc.DeleteSemester(c.UserContext(), semUUID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}
