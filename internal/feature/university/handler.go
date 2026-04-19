package university

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/pkg/response"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// List godoc
// @Summary      List universities
// @Description  Returns a list of all universities.
// @Tags         Universities
// @Produce      json
// @Success      200 {object} response.Response{data=[]UniversityResponse}
// @Failure      500 {object} response.Response
// @Router       /universities [get]
func (h *Handler) List(c *fiber.Ctx) error {
	universities, err := h.svc.ListUniversities(c.UserContext())
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "universities retrieved", ToUniversityResponses(universities))
}

// GetByID godoc
// @Summary      Get university by ID
// @Description  Returns a specific university by its ID.
// @Tags         Universities
// @Produce      json
// @Param        id path int true "University ID"
// @Success      200 {object} response.Response{data=UniversityResponse}
// @Failure      404 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /universities/{id} [get]
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid university id", nil)
	}
	universities, err := h.svc.GetByID(c.UserContext(), int64(id))
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "university retrieved", ToUniversityResponse(universities))
}
