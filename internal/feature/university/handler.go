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

// GetByUUID godoc
// @Summary      Get university by UUID
// @Description  Returns a specific university by its UUID.
// @Tags         Universities
// @Produce      json
// @Param        id path string true "University UUID"
// @Success      200 {object} response.Response{data=UniversityResponse}
// @Failure      404 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /universities/{id} [get]
func (h *Handler) GetByUUID(c *fiber.Ctx) error {
	uuid := c.Params("id")
	u, err := h.svc.GetByUUID(c.UserContext(), uuid)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "university retrieved", ToUniversityResponse(u))
}
