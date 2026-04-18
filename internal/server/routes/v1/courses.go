package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/app"
	"github.com/victorive/klaro-backend/internal/feature/course"
	"github.com/victorive/klaro-backend/internal/feature/file"
	"github.com/victorive/klaro-backend/internal/middleware"
)

func courseRoutes(r fiber.Router, c *app.Container) {
	h := course.NewHandler(c.CourseService)
	fh := file.NewHandler(c.FileService)
	auth := []fiber.Handler{middleware.Auth(c.AuthService), middleware.OTPVerified(c.AuthService)}

	courses := r.Group("/courses", auth...)
	courses.Post("", h.CreateCourse)
	courses.Get("/:id", h.GetCourse)
	courses.Patch("/:id", h.UpdateCourse)
	courses.Delete("/:id", h.DeleteCourse)
	courses.Get("/:courseId/lectures", h.ListLectures)
	courses.Post("/:courseId/lectures", h.CreateLecture)
	courses.Get("/:courseId/files", fh.ListByCourse)

	r.Group("/levels", auth...).Get("/:levelId/courses", h.ListCoursesByLevel)
	r.Group("/semesters", auth...).Get("/:semesterId/courses", h.ListCoursesBySemester)

	lectures := r.Group("/lectures", auth...)
	lectures.Get("/:id", h.GetLecture)
	lectures.Patch("/:id", h.UpdateLecture)
	lectures.Patch("/:id/archive", h.ArchiveLecture)
	lectures.Get("/:lectureId/files", fh.ListByLecture)
}
