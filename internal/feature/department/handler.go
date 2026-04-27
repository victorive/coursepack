package department

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/domain"
	"github.com/victorive/klaro-backend/internal/middleware"
	"github.com/victorive/klaro-backend/internal/pkg/pagination"
	"github.com/victorive/klaro-backend/internal/pkg/response"
)

type UserRepository interface {
	FindByUUID(ctx context.Context, uuid string) (*domain.User, error)
}

type Handler struct {
	svc      Service
	userRepo UserRepository
}

func NewHandler(svc Service, userRepo UserRepository) *Handler {
	return &Handler{svc: svc, userRepo: userRepo}
}

func (h *Handler) resolveUniversityID(c *fiber.Ctx) (int64, bool) {
	userUUID, ok := middleware.UserUUID(c)
	if !ok {
		return 0, false
	}

	u, err := h.userRepo.FindByUUID(c.UserContext(), userUUID)
	if err != nil || u == nil || u.UniversityID == nil {
		return 0, false
	}

	return *u.UniversityID, true
}

// List godoc
// @Summary      List departments
// @Description  Returns a paginated list of departments scoped to the authenticated user's university. Returns an empty list if the user has not yet set a university.
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
	universityID, hasUniversity := h.resolveUniversityID(c)
	if !hasUniversity {
		return response.JSON(c, http.StatusOK, true, "departments retrieved", pagination.NewPaginated(c, []*DepartmentResponse{}, 0, q))
	}

	items, total, err := h.svc.ListByUniversity(c.UserContext(), universityID, *q)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, true, "departments retrieved", pagination.NewPaginated(c, items, total, q))
}

// Get godoc
// @Summary      Get department
// @Description  Returns a single department by ID, scoped to the authenticated user's university.
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

	universityID, hasUniversity := h.resolveUniversityID(c)
	if !hasUniversity {
		return response.JSON(c, http.StatusNotFound, false, "department not found", nil)
	}

	dept, err := h.svc.GetByIDAndUniversity(c.UserContext(), id, universityID)
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
