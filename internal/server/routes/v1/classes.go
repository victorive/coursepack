package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/coursepack/internal/app"
	"github.com/victorive/coursepack/internal/domain"
	"github.com/victorive/coursepack/internal/feature/class"
	"github.com/victorive/coursepack/internal/middleware"
)

func classRoutes(r fiber.Router, c *app.Container) {
	h := class.NewHandler(c.ClassService)
	auth := []fiber.Handler{middleware.Auth(c.AuthService), middleware.OTPVerified(c.AuthService)}

	global := func(perm domain.PermissionName) fiber.Handler {
		return middleware.RequireGlobalPermission(c.RBAC, perm)
	}
	scoped := func(perm domain.PermissionName, param string) fiber.Handler {
		return middleware.RequireClassPermission(c.RBAC, perm, param)
	}

	r.Group("/levels", auth...).Get("", h.ListLevels)

	classes := r.Group("/classes", auth...)
	classes.Post("", global(domain.PermManageClasses), h.Create)
	classes.Get("", h.List)
	classes.Get("/:id", h.Get)
	classes.Patch("/:id", scoped(domain.PermMembersManage, "id"), h.Update)
	classes.Delete("/:id", scoped(domain.PermMembersManage, "id"), h.Delete)

	classes.Post("/join", h.RequestJoin)
	classes.Get("/:id/join-requests", scoped(domain.PermMembersManage, "id"), h.ListPendingJoinRequests)
	classes.Patch("/:id/join-requests/:requestId", scoped(domain.PermMembersManage, "id"), h.ReviewJoinRequest)

	classes.Get("/:id/members", scoped(domain.PermMembersList, "id"), h.ListMembers)
	classes.Delete("/:id/members/:userUUID", scoped(domain.PermMembersManage, "id"), h.RemoveMember)
	classes.Post("/:id/invite-code/regenerate", scoped(domain.PermMembersManage, "id"), h.RegenerateInviteCode)
}
