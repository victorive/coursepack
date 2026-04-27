package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/app"
	"github.com/victorive/klaro-backend/internal/feature/course"
	"github.com/victorive/klaro-backend/internal/feature/file"
	"github.com/victorive/klaro-backend/internal/middleware"
)

func courseRoutes(r fiber.Router, c *app.Container) {
	h := course.NewHandler(c.CourseService, c.AuthRepository)
	fh := file.NewHandler(c.FileService)
	auth := []fiber.Handler{middleware.Auth(c.AuthService), middleware.OTPVerified(c.AuthService)}

	courses := r.Group("/courses", auth...)
	courses.Get("", h.ListCourses)
	courses.Get("/:id", h.GetCourse)

	classGroup := r.Group("/classes", auth...)
	classGroup.Post("/:id/courses", h.AddClassCourse)
	classGroup.Get("/:id/courses", h.ListClassCourses)
	classGroup.Delete("/:id/courses/:classCourseId", h.RemoveClassCourse)

	classCourseGroup := r.Group("/courses", auth...)
	classCourseGroup.Get("/:id", h.GetClassCourse)
	classCourseGroup.Post("/:id/lectures", h.CreateLecture)
	classCourseGroup.Get("/:id/lectures", h.ListLectures)
	classCourseGroup.Get("/:id/files", fh.ListByClassCourse)

	lectures := r.Group("/lectures", auth...)
	lectures.Get("/:id", h.GetLecture)
	lectures.Patch("/:id", h.UpdateLecture)
	lectures.Patch("/:id/archive", h.ArchiveLecture)
	lectures.Get("/:lectureId/files", fh.ListByLecture)
}
