package department

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

// Create godoc
// @Summary      Create department
// @Description  Creates a new department. Super-admin only.
// @Tags         Departments
// @Accept       json
// @Produce      json
// @Param        body body CreateDepartmentRequest true "Create department payload"
// @Success      201 {object} response.Response{data=DepartmentResponse}
// @Failure      401 {object} response.Response
// @Failure      422 {object} response.ValidationErrorResponse
// @Security     BearerAuth
// @Router       /departments [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateDepartmentRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	dept, err := h.svc.Create(c.UserContext(), &req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusCreated, true, "department created", dept)
}

// List godoc
// @Summary      List departments
// @Description  Returns a paginated list of all departments.
// @Tags         Departments
// @Produce      json
// @Param        page    query int false "Page number"
// @Param        perPage query int false "Items per page"
// @Success      200 {object} response.Response{data=object}
// @Failure      401 {object} response.Response
// @Security     BearerAuth
// @Router       /departments [get]
func (h *Handler) List(c *fiber.Ctx) error {
	q := pagination.Parse(c)
	items, total, err := h.svc.List(c.UserContext(), *q)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "departments retrieved", pagination.NewPaginated(c, items, total, q))
}

// Get godoc
// @Summary      Get department
// @Description  Returns a single department by ID.
// @Tags         Departments
// @Produce      json
// @Param        id path int true "Department ID"
// @Success      200 {object} response.Response{data=DepartmentResponse}
// @Failure      404 {object} response.Response
// @Security     BearerAuth
// @Router       /departments/{id} [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid id", nil)
	}
	dept, err := h.svc.Get(c.UserContext(), id)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "department retrieved", dept)
}

// Update godoc
// @Summary      Update department
// @Description  Updates an existing department. Super-admin only.
// @Tags         Departments
// @Accept       json
// @Produce      json
// @Param        id   path int                    true "Department ID"
// @Param        body body UpdateDepartmentRequest true "Update department payload"
// @Success      200 {object} response.Response{data=DepartmentResponse}
// @Failure      404 {object} response.Response
// @Security     BearerAuth
// @Router       /departments/{id} [patch]
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid id", nil)
	}
	var req UpdateDepartmentRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	dept, err := h.svc.Update(c.UserContext(), id, &req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "department updated", dept)
}

// Delete godoc
// @Summary      Delete department
// @Description  Deletes a department. Super-admin only.
// @Tags         Departments
// @Produce      json
// @Param        id path int true "Department ID"
// @Success      204
// @Failure      404 {object} response.Response
// @Security     BearerAuth
// @Router       /departments/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid id", nil)
	}
	if err := h.svc.Delete(c.UserContext(), id); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}
