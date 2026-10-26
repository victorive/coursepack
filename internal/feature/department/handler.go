package department

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/coursepack/internal/domain"
	"github.com/victorive/coursepack/internal/middleware"
	"github.com/victorive/coursepack/internal/pkg/pagination"
	"github.com/victorive/coursepack/internal/pkg/response"
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
	userUUID, err := middleware.UserUUID(c)
	if err != nil {
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
// @Description  Returns a single department by UUID, scoped to the authenticated user's university.
// @Tags         Departments
// @Produce      json
// @Param        id path string true "Department UUID"
// @Success      200 {object} response.Response{data=DepartmentResponse}
// @Failure      404 {object} response.Response
// @Security     BearerAuth
// @Router       /departments/{id} [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	uuid := c.Params("id")

	universityID, hasUniversity := h.resolveUniversityID(c)
	if !hasUniversity {
		return response.JSON(c, http.StatusNotFound, false, "department not found", nil)
	}

	dept, err := h.svc.GetByUUIDAndUniversity(c.UserContext(), uuid, universityID)
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
