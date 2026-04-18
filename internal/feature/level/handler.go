package level

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/middleware"
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
// @Summary      Create level
// @Tags         Levels
// @Accept       json
// @Produce      json
// @Param        body body CreateLevelRequest true "Create level payload"
// @Success      201 {object} response.Response{data=LevelResponse}
// @Security     BearerAuth
// @Router       /levels [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateLevelRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	uuid, _ := c.Locals(middleware.LocalsUserId).(string)
	level, err := h.svc.Create(c.UserContext(), uuid, &req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusCreated, true, "level created", level)
}

// List godoc
// @Summary      List levels
// @Tags         Levels
// @Produce      json
// @Success      200 {object} response.Response{data=object}
// @Security     BearerAuth
// @Router       /levels [get]
func (h *Handler) List(c *fiber.Ctx) error {
	q := pagination.Parse(c)
	items, total, err := h.svc.List(c.UserContext(), *q)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "levels retrieved", pagination.NewPaginated(c, items, total, q))
}

// Get godoc
// @Summary      Get level
// @Tags         Levels
// @Produce      json
// @Param        id path int true "Level ID"
// @Success      200 {object} response.Response{data=LevelResponse}
// @Security     BearerAuth
// @Router       /levels/{id} [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid id", nil)
	}
	level, err := h.svc.Get(c.UserContext(), id)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "level retrieved", level)
}

// Delete godoc
// @Summary      Delete level
// @Tags         Levels
// @Produce      json
// @Param        id path int true "Level ID"
// @Success      204
// @Security     BearerAuth
// @Router       /levels/{id} [delete]
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

// RequestJoin godoc
// @Summary      Request to join a level
// @Tags         Levels
// @Accept       json
// @Produce      json
// @Param        body body JoinRequest true "Join request payload"
// @Success      201 {object} response.Response{data=JoinRequestResponse}
// @Security     BearerAuth
// @Router       /levels/join [post]
func (h *Handler) RequestJoin(c *fiber.Ctx) error {
	var req JoinRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	uuid, _ := c.Locals(middleware.LocalsUserId).(string)
	jr, err := h.svc.RequestJoin(c.UserContext(), uuid, &req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusCreated, true, "join request submitted", jr)
}

// ListPendingJoinRequests godoc
// @Summary      List pending join requests for a level
// @Tags         Levels
// @Produce      json
// @Param        id path int true "Level ID"
// @Success      200 {object} response.Response{data=[]JoinRequestResponse}
// @Security     BearerAuth
// @Router       /levels/{id}/join-requests [get]
func (h *Handler) ListPendingJoinRequests(c *fiber.Ctx) error {
	levelID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid level id", nil)
	}
	requests, err := h.svc.ListPendingJoinRequests(c.UserContext(), levelID)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "join requests retrieved", requests)
}

// ReviewJoinRequest godoc
// @Summary      Approve or reject a join request
// @Tags         Levels
// @Accept       json
// @Produce      json
// @Param        id        path int              true "Level ID"
// @Param        requestId path int              true "Join Request ID"
// @Param        body      body ReviewJoinRequest true "Review payload"
// @Success      200 {object} response.Response{data=JoinRequestResponse}
// @Security     BearerAuth
// @Router       /levels/{id}/join-requests/{requestId} [patch]
func (h *Handler) ReviewJoinRequest(c *fiber.Ctx) error {
	levelID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid level id", nil)
	}
	requestID, err := strconv.ParseInt(c.Params("requestId"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid request id", nil)
	}
	var req ReviewJoinRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	reviewerUUID, _ := c.Locals(middleware.LocalsUserId).(string)
	jr, err := h.svc.ReviewJoinRequest(c.UserContext(), reviewerUUID, levelID, requestID, &req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "join request reviewed", jr)
}

// ListMembers godoc
// @Summary      List members of a level
// @Tags         Levels
// @Produce      json
// @Param        id path int true "Level ID"
// @Success      200 {object} response.Response{data=[]MemberResponse}
// @Security     BearerAuth
// @Router       /levels/{id}/members [get]
func (h *Handler) ListMembers(c *fiber.Ctx) error {
	levelID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid level id", nil)
	}
	members, err := h.svc.ListMembers(c.UserContext(), levelID)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "members retrieved", members)
}

// RemoveMember godoc
// @Summary      Remove a member from a level
// @Tags         Levels
// @Produce      json
// @Param        id      path int true "Level ID"
// @Param        userId  path int true "User ID"
// @Success      200 {object} response.Response
// @Security     BearerAuth
// @Router       /levels/{id}/members/{userId} [delete]
func (h *Handler) RemoveMember(c *fiber.Ctx) error {
	levelID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid level id", nil)
	}
	userID, err := strconv.ParseInt(c.Params("userId"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid user id", nil)
	}
	if err := h.svc.RemoveMember(c.UserContext(), levelID, userID); err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "member removed", nil)
}
