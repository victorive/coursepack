package department

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/coursepack/internal/permissions"
	"github.com/victorive/coursepack/internal/pkg/pagination"
	"github.com/victorive/coursepack/internal/pkg/response"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
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
	actor, err := permissions.ActorFrom(c.UserContext())
	if err != nil {
		return err
	}

	items, total, err := h.svc.ListByUniversity(c.UserContext(), actor, *q)
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

	actor, err := permissions.ActorFrom(c.UserContext())
	if err != nil {
		return err
	}

	dept, err := h.svc.GetByUUIDAndUniversity(c.UserContext(), actor, uuid)
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
