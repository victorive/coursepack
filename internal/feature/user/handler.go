package user

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/coursepack/internal/middleware"
	"github.com/victorive/coursepack/internal/pkg/response"
	"github.com/victorive/coursepack/internal/pkg/validator"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// UpdateProfile godoc
// @Summary      Update profile
// @Description  Updates the authenticated user's profile. All fields are optional — include only what you want to change. Set university_id to link the user to a university.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        body body UpdateProfileRequest true "Update profile payload"
// @Success      200 {object} response.Response
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      422 {object} response.ValidationErrorResponse "Validation error"
// @Security     BearerAuth
// @Router       /users/me [patch]
func (h *Handler) UpdateProfile(c *fiber.Ctx) error {
	var req UpdateProfileRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	uuid, _ := c.Locals(middleware.LocalsUserId).(string)
	if err := h.svc.UpdateProfile(c.UserContext(), uuid, &req); err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "profile updated", nil)
}
