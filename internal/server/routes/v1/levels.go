package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/app"
	"github.com/victorive/klaro-backend/internal/domain"
	"github.com/victorive/klaro-backend/internal/feature/level"
	"github.com/victorive/klaro-backend/internal/middleware"
)

func levelRoutes(r fiber.Router, c *app.Container) {
	h := level.NewHandler(c.LevelService)
	auth := []fiber.Handler{middleware.Auth(c.AuthService), middleware.OTPVerified(c.AuthService)}

	global := func(perm domain.PermissionName) fiber.Handler {
		return middleware.RequireGlobalPermission(c.RBAC, perm)
	}
	scoped := func(perm domain.PermissionName, param string) fiber.Handler {
		return middleware.RequireLevelPermission(c.RBAC, perm, param)
	}

	levels := r.Group("/levels", auth...)
	levels.Post("", global(domain.PermManageLevels), h.Create)
	levels.Get("", h.List)
	levels.Get("/:id", h.Get)
	levels.Delete("/:id", scoped(domain.PermMembersManage, "id"), h.Delete)

	levels.Post("/join", h.RequestJoin)
	levels.Get("/:id/join-requests", scoped(domain.PermMembersManage, "id"), h.ListPendingJoinRequests)
	levels.Patch("/:id/join-requests/:requestId", scoped(domain.PermMembersManage, "id"), h.ReviewJoinRequest)

	levels.Get("/:id/members", scoped(domain.PermMembersList, "id"), h.ListMembers)
	levels.Delete("/:id/members/:userId", scoped(domain.PermMembersManage, "id"), h.RemoveMember)
}
