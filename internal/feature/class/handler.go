package class

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/coursepack/internal/middleware"
	"github.com/victorive/coursepack/internal/pkg/pagination"
	"github.com/victorive/coursepack/internal/pkg/response"
	"github.com/victorive/coursepack/internal/pkg/validator"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// Create godoc
// @Summary      Create class
// @Description  Creates a new class. The authenticated user becomes the class rep.
// @Tags         Classes
// @Accept       json
// @Produce      json
// @Param        body body CreateClassRequest true "Create class payload"
// @Success      201 {object} response.Response{data=ClassResponse}
// @Failure      422 {object} response.ValidationErrorResponse "Validation error"
// @Security     BearerAuth
// @Router       /classes [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateClassRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	uuid, err := middleware.UserUUID(c)
	if err != nil {
		return err
	}
	cls, err := h.svc.Create(c.UserContext(), uuid, &req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusCreated, true, "class created", cls)
}

// List godoc
// @Summary      List classes
// @Description  Returns a paginated list of all classes.
// @Tags         Classes
// @Produce      json
// @Param        page  query int false "Page number"
// @Param        limit query int false "Items per page"
// @Success      200 {object} response.Response{data=object}
// @Security     BearerAuth
// @Router       /classes [get]
func (h *Handler) List(c *fiber.Ctx) error {
	q := pagination.Parse(c)
	uuid, err := middleware.UserUUID(c)
	if err != nil {
		return err
	}
	items, total, err := h.svc.List(c.UserContext(), uuid, *q)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "classes retrieved", pagination.NewPaginated(c, items, total, q))
}

// Get godoc
// @Summary      Get class
// @Description  Returns a single class by UUID.
// @Tags         Classes
// @Produce      json
// @Param        id path string true "Class UUID"
// @Success      200 {object} response.Response{data=ClassResponse}
// @Failure      401,404 {object} response.Response
// @Security     BearerAuth
// @Router       /classes/{id} [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	classUUID := c.Params("id")
	userUUID, err := middleware.UserUUID(c)
	if err != nil {
		return err
	}
	cls, err := h.svc.Get(c.UserContext(), userUUID, classUUID)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "class retrieved", cls)
}

// Delete godoc
// @Summary      Delete class
// @Description  Permanently deletes a class.
// @Tags         Classes
// @Produce      json
// @Param        id path string true "Class UUID"
// @Success      204
// @Failure      401,403,404 {object} response.Response
// @Security     BearerAuth
// @Router       /classes/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	classUUID := c.Params("id")
	if err := h.svc.Delete(c.UserContext(), classUUID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// Update godoc
// @Summary      Update class
// @Description  Updates mutable fields (name, current level) of a class.
// @Tags         Classes
// @Accept       json
// @Produce      json
// @Param        id   path string             true "Class UUID"
// @Param        body body UpdateClassRequest true "Update class payload"
// @Success      200 {object} response.Response{data=ClassResponse}
// @Failure      401,403,404,422 {object} response.Response
// @Security     BearerAuth
// @Router       /classes/{id} [patch]
func (h *Handler) Update(c *fiber.Ctx) error {
	classUUID := c.Params("id")
	var req UpdateClassRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	userUUID, err := middleware.UserUUID(c)
	if err != nil {
		return err
	}
	cls, err := h.svc.Update(c.UserContext(), userUUID, classUUID, &req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "class updated", cls)
}

// RequestJoin godoc
// @Summary      Request to join a class
// @Description  Submits a join request for the authenticated user using an invite code.
// @Tags         Classes
// @Accept       json
// @Produce      json
// @Param        body body JoinRequest true "Join request payload"
// @Success      201 {object} response.Response{data=JoinRequestResponse}
// @Failure      422 {object} response.ValidationErrorResponse "Validation error"
// @Security     BearerAuth
// @Router       /classes/join [post]
func (h *Handler) RequestJoin(c *fiber.Ctx) error {
	var req JoinRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	uuid, err := middleware.UserUUID(c)
	if err != nil {
		return err
	}
	jr, err := h.svc.RequestJoin(c.UserContext(), uuid, &req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusCreated, true, "join request submitted", jr)
}

// ListPendingJoinRequests godoc
// @Summary      List pending join requests
// @Description  Returns all pending join requests for a class.
// @Tags         Classes
// @Produce      json
// @Param        id path string true "Class UUID"
// @Success      200 {object} response.Response{data=[]JoinRequestResponse}
// @Failure      401,403,404 {object} response.Response
// @Security     BearerAuth
// @Router       /classes/{id}/join-requests [get]
func (h *Handler) ListPendingJoinRequests(c *fiber.Ctx) error {
	classUUID := c.Params("id")
	requests, err := h.svc.ListPendingJoinRequests(c.UserContext(), classUUID)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "join requests retrieved", requests)
}

// ReviewJoinRequest godoc
// @Summary      Review a join request
// @Description  Approves or rejects a pending join request.
// @Tags         Classes
// @Accept       json
// @Produce      json
// @Param        id        path string            true "Class UUID"
// @Param        requestId path int               true "Join request ID"
// @Param        body      body ReviewJoinRequest true "Review payload"
// @Success      200 {object} response.Response{data=JoinRequestResponse}
// @Failure      400,401,403,404,422 {object} response.Response
// @Security     BearerAuth
// @Router       /classes/{id}/join-requests/{requestId} [patch]
func (h *Handler) ReviewJoinRequest(c *fiber.Ctx) error {
	classUUID := c.Params("id")
	requestID, err := strconv.ParseInt(c.Params("requestId"), 10, 64)
	if err != nil {
		return response.JSON(c, http.StatusBadRequest, false, "invalid request id", nil)
	}
	var req ReviewJoinRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	reviewerUUID, err := middleware.UserUUID(c)
	if err != nil {
		return err
	}
	jr, err := h.svc.ReviewJoinRequest(c.UserContext(), reviewerUUID, classUUID, requestID, &req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "join request reviewed", jr)
}

// ListMembers godoc
// @Summary      List class members
// @Description  Returns all members of a class.
// @Tags         Classes
// @Produce      json
// @Param        id path string true "Class UUID"
// @Success      200 {object} response.Response{data=[]MemberResponse}
// @Failure      401,403,404 {object} response.Response
// @Security     BearerAuth
// @Router       /classes/{id}/members [get]
func (h *Handler) ListMembers(c *fiber.Ctx) error {
	classUUID := c.Params("id")
	members, err := h.svc.ListMembers(c.UserContext(), classUUID)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "members retrieved", members)
}

// RemoveMember godoc
// @Summary      Remove a class member
// @Description  Removes a user from a class.
// @Tags         Classes
// @Produce      json
// @Param        id       path string true "Class UUID"
// @Param        userUUID path string true "User UUID"
// @Success      200 {object} response.Response
// @Failure      404 {object} response.Response "Not found"
// @Security     BearerAuth
// @Router       /classes/{id}/members/{userUUID} [delete]
func (h *Handler) RemoveMember(c *fiber.Ctx) error {
	classUUID := c.Params("id")
	targetUserUUID := c.Params("userUUID")
	if err := h.svc.RemoveMember(c.UserContext(), classUUID, targetUserUUID); err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "member removed", nil)
}

// ListLevels godoc
// @Summary      List levels
// @Description  Returns all available academic levels.
// @Tags         Classes
// @Produce      json
// @Success      200 {object} response.Response{data=[]LevelResponse}
// @Security     BearerAuth
// @Router       /levels [get]
func (h *Handler) ListLevels(c *fiber.Ctx) error {
	levels, err := h.svc.ListLevels(c.UserContext())
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "levels retrieved", levels)
}

// RegenerateInviteCode godoc
// @Summary      Regenerate invite code
// @Description  Generates a new invite code for the class, invalidating the old one. Requires members:manage permission.
// @Tags         Classes
// @Produce      json
// @Param        id path string true "Class UUID"
// @Success      200 {object} response.Response{data=ClassResponse}
// @Failure      400,401,403,404 {object} response.Response
// @Security     BearerAuth
// @Router       /classes/{id}/invite-code/regenerate [post]
func (h *Handler) RegenerateInviteCode(c *fiber.Ctx) error {
	classUUID := c.Params("id")

	userUUID, err := middleware.UserUUID(c)
	if err != nil {
		return err
	}

	cls, err := h.svc.RegenerateInviteCode(c.UserContext(), userUUID, classUUID)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, true, "invite code regenerated", cls)
}
