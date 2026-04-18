package cohort

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
// @Summary      Create cohort
// @Tags         Cohorts
// @Accept       json
// @Produce      json
// @Param        body body CreateCohortRequest true "Create cohort payload"
// @Success      201 {object} response.Response{data=CohortResponse}
// @Security     BearerAuth
// @Router       /cohorts [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateCohortRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	uuid, _ := c.Locals(middleware.LocalsUserId).(string)
	cohort, err := h.svc.Create(c.UserContext(), uuid, &req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusCreated, true, "cohort created", cohort)
}

// List godoc
// @Summary      List cohorts
// @Tags         Cohorts
// @Produce      json
// @Success      200 {object} response.Response{data=object}
// @Security     BearerAuth
// @Router       /cohorts [get]
func (h *Handler) List(c *fiber.Ctx) error {
	q := pagination.Parse(c)
	items, total, err := h.svc.List(c.UserContext(), *q)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "cohorts retrieved", pagination.NewPaginated(c, items, total, q))
}

// Get godoc
// @Summary      Get cohort
// @Tags         Cohorts
// @Produce      json
// @Param        id path int true "Cohort ID"
// @Success      200 {object} response.Response{data=CohortResponse}
// @Security     BearerAuth
// @Router       /cohorts/{id} [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid id", nil)
	}
	cohort, err := h.svc.Get(c.UserContext(), id)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "cohort retrieved", cohort)
}

// Delete godoc
// @Summary      Delete cohort
// @Tags         Cohorts
// @Produce      json
// @Param        id path int true "Cohort ID"
// @Success      204
// @Security     BearerAuth
// @Router       /cohorts/{id} [delete]
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
// @Summary      Request to join a cohort
// @Tags         Cohorts
// @Accept       json
// @Produce      json
// @Param        body body JoinRequest true "Join request payload"
// @Success      201 {object} response.Response{data=JoinRequestResponse}
// @Security     BearerAuth
// @Router       /cohorts/join [post]
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
// @Summary      List pending join requests for a cohort
// @Tags         Cohorts
// @Produce      json
// @Param        id path int true "Cohort ID"
// @Success      200 {object} response.Response{data=[]JoinRequestResponse}
// @Security     BearerAuth
// @Router       /cohorts/{id}/join-requests [get]
func (h *Handler) ListPendingJoinRequests(c *fiber.Ctx) error {
	cohortID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid cohort id", nil)
	}
	requests, err := h.svc.ListPendingJoinRequests(c.UserContext(), cohortID)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "join requests retrieved", requests)
}

// ReviewJoinRequest godoc
// @Summary      Approve or reject a join request
// @Tags         Cohorts
// @Accept       json
// @Produce      json
// @Param        id        path int              true "Cohort ID (unused, validation only)"
// @Param        requestId path int              true "Join Request ID"
// @Param        body      body ReviewJoinRequest true "Review payload"
// @Success      200 {object} response.Response{data=JoinRequestResponse}
// @Security     BearerAuth
// @Router       /cohorts/{id}/join-requests/{requestId} [patch]
func (h *Handler) ReviewJoinRequest(c *fiber.Ctx) error {
	cohortID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid cohort id", nil)
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
	jr, err := h.svc.ReviewJoinRequest(c.UserContext(), reviewerUUID, cohortID, requestID, &req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "join request reviewed", jr)
}

// ListMembers godoc
// @Summary      List members of a cohort
// @Tags         Cohorts
// @Produce      json
// @Param        id path int true "Cohort ID"
// @Success      200 {object} response.Response{data=[]MemberResponse}
// @Security     BearerAuth
// @Router       /cohorts/{id}/members [get]
func (h *Handler) ListMembers(c *fiber.Ctx) error {
	cohortID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid cohort id", nil)
	}
	members, err := h.svc.ListMembers(c.UserContext(), cohortID)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "members retrieved", members)
}

// RemoveMember godoc
// @Summary      Remove a member from a cohort
// @Tags         Cohorts
// @Produce      json
// @Param        id      path int true "Cohort ID"
// @Param        userId  path int true "User ID"
// @Success      200 {object} response.Response
// @Security     BearerAuth
// @Router       /cohorts/{id}/members/{userId} [delete]
func (h *Handler) RemoveMember(c *fiber.Ctx) error {
	cohortID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid cohort id", nil)
	}
	userID, err := strconv.ParseInt(c.Params("userId"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid user id", nil)
	}
	if err := h.svc.RemoveMember(c.UserContext(), cohortID, userID); err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "member removed", nil)
}
