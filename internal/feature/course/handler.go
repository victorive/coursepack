package course

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/middleware"
	"github.com/victorive/klaro-backend/internal/pkg/apperr"
	"github.com/victorive/klaro-backend/internal/pkg/response"
	"github.com/victorive/klaro-backend/internal/pkg/validator"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// CreateCourse godoc
// @Summary      Create course
// @Description  Creates a new course linked to a level and semester.
// @Tags         Courses
// @Accept       json
// @Produce      json
// @Param        body body CreateCourseRequest true "Create course payload"
// @Success      201 {object} response.Response{data=CourseResponse}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Security     BearerAuth
// @Router       /courses [post]
func (h *Handler) CreateCourse(c *fiber.Ctx) error {
	actorUUID, ok := middleware.UserUUID(c)
	if !ok {
		return apperr.New("unauthorized", http.StatusUnauthorized)
	}
	
	var req CreateCourseRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}

	resp, err := h.svc.CreateCourse(c.UserContext(), actorUUID, req)
	if err != nil {
		return err
	}
	
	return response.JSON(c, http.StatusCreated, true, "course created", resp)
}

// ListCoursesByLevel godoc
// @Summary      List courses by level
// @Description  Returns all courses belonging to a level.
// @Tags         Courses
// @Produce      json
// @Param        levelId path int true "Level ID"
// @Success      200 {object} response.Response{data=[]CourseResponse}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Security     BearerAuth
// @Router       /levels/{levelId}/courses [get]
func (h *Handler) ListCoursesByLevel(c *fiber.Ctx) error {
	levelID, err := strconv.ParseInt(c.Params("levelId"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	resp, err := h.svc.ListCoursesByLevel(c.UserContext(), levelID)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, true, "courses fetched", resp)
}

// ListCoursesBySemester godoc
// @Summary      List courses by semester
// @Description  Returns all courses belonging to a semester.
// @Tags         Courses
// @Produce      json
// @Param        semesterId path int true "Semester ID"
// @Success      200 {object} response.Response{data=[]CourseResponse}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Security     BearerAuth
// @Router       /semesters/{semesterId}/courses [get]
func (h *Handler) ListCoursesBySemester(c *fiber.Ctx) error {
	semesterID, err := strconv.ParseInt(c.Params("semesterId"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	resp, err := h.svc.ListCoursesBySemester(c.UserContext(), semesterID)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, true, "courses fetched", resp)
}

// GetCourse godoc
// @Summary      Get course
// @Description  Returns a single course by ID.
// @Tags         Courses
// @Produce      json
// @Param        id path int true "Course ID"
// @Success      200 {object} response.Response{data=CourseResponse}
// @Failure      404 {object} response.Response
// @Security     BearerAuth
// @Router       /courses/{id} [get]
func (h *Handler) GetCourse(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	resp, err := h.svc.GetCourse(c.UserContext(), id)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, true, "course fetched", resp)
}

// UpdateCourse godoc
// @Summary      Update course
// @Description  Updates an existing course.
// @Tags         Courses
// @Accept       json
// @Produce      json
// @Param        id   path int                true "Course ID"
// @Param        body body UpdateCourseRequest true "Update course payload"
// @Success      200 {object} response.Response{data=CourseResponse}
// @Failure      404 {object} response.Response
// @Security     BearerAuth
// @Router       /courses/{id} [patch]
func (h *Handler) UpdateCourse(c *fiber.Ctx) error {
	actorUUID, ok := middleware.UserUUID(c)
	if !ok {
		return apperr.New("unauthorized", http.StatusUnauthorized)
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	var req UpdateCourseRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}

	resp, err := h.svc.UpdateCourse(c.UserContext(), actorUUID, id, req)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, true, "course updated", resp)
}

// DeleteCourse godoc
// @Summary      Delete course
// @Description  Permanently deletes a course.
// @Tags         Courses
// @Produce      json
// @Param        id path int true "Course ID"
// @Success      204
// @Failure      404 {object} response.Response
// @Security     BearerAuth
// @Router       /courses/{id} [delete]
func (h *Handler) DeleteCourse(c *fiber.Ctx) error {
	actorUUID, ok := middleware.UserUUID(c)
	if !ok {
		return apperr.New("unauthorized", http.StatusUnauthorized)
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := h.svc.DeleteCourse(c.UserContext(), actorUUID, id); err != nil {
		return err
	}

	return c.SendStatus(http.StatusNoContent)
}

// ---- Lectures ----

// CreateLecture godoc
// @Summary      Create lecture
// @Description  Creates a new lecture under a course.
// @Tags         Lectures
// @Accept       json
// @Produce      json
// @Param        courseId path int                  true "Course ID"
// @Param        body     body CreateLectureRequest true "Create lecture payload"
// @Success      201 {object} response.Response{data=LectureResponse}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Security     BearerAuth
// @Router       /courses/{courseId}/lectures [post]
func (h *Handler) CreateLecture(c *fiber.Ctx) error {
	actorUUID, ok := middleware.UserUUID(c)
	if !ok {
		return apperr.New("unauthorized", http.StatusUnauthorized)
	}

	var req CreateLectureRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}

	resp, err := h.svc.CreateLecture(c.UserContext(), actorUUID, req)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusCreated, true, "lecture created", resp)
}

// ListLectures godoc
// @Summary      List lectures for a course
// @Description  Returns all published lectures for a course.
// @Tags         Lectures
// @Produce      json
// @Param        courseId path int true "Course ID"
// @Success      200 {object} response.Response{data=[]LectureResponse}
// @Failure      401 {object} response.Response
// @Security     BearerAuth
// @Router       /courses/{courseId}/lectures [get]
func (h *Handler) ListLectures(c *fiber.Ctx) error {
	courseID, err := strconv.ParseInt(c.Params("courseId"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	resp, err := h.svc.ListLectures(c.UserContext(), courseID)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, true, "lectures fetched", resp)
}

// GetLecture godoc
// @Summary      Get a lecture
// @Description  Returns a single lecture by ID.
// @Tags         Lectures
// @Produce      json
// @Param        id path int true "Lecture ID"
// @Success      200 {object} response.Response{data=LectureResponse}
// @Failure      400,401,404 {object} response.Response
// @Security     BearerAuth
// @Router       /lectures/{id} [get]
func (h *Handler) GetLecture(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	resp, err := h.svc.GetLecture(c.UserContext(), id)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, true, "lecture fetched", resp)
}

// UpdateLecture godoc
// @Summary      Update a lecture
// @Description  Updates a lecture's details.
// @Tags         Lectures
// @Accept       json
// @Produce      json
// @Param        id   path     int                 true "Lecture ID"
// @Param        body body     UpdateLectureRequest true "Update payload"
// @Success      200  {object} response.Response{data=LectureResponse}
// @Failure      400,401,403,404 {object} response.Response
// @Security     BearerAuth
// @Router       /lectures/{id} [patch]
func (h *Handler) UpdateLecture(c *fiber.Ctx) error {
	actorUUID, ok := middleware.UserUUID(c)
	if !ok {
		return apperr.New("unauthorized", http.StatusUnauthorized)
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}
	var req UpdateLectureRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}

	resp, err := h.svc.UpdateLecture(c.UserContext(), actorUUID, id, req)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, true, "lecture updated", resp)
}

// ArchiveLecture godoc
// @Summary      Archive a lecture
// @Description  Moves a lecture to the archived state.
// @Tags         Lectures
// @Produce      json
// @Param        id path int true "Lecture ID"
// @Success      200 {object} response.Response
// @Failure      400,401,403,404 {object} response.Response
// @Security     BearerAuth
// @Router       /lectures/{id}/archive [patch]
func (h *Handler) ArchiveLecture(c *fiber.Ctx) error {
	actorUUID, ok := middleware.UserUUID(c)
	if !ok {
		return apperr.New("unauthorized", http.StatusUnauthorized)
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err := h.svc.ArchiveLecture(c.UserContext(), actorUUID, id); err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "lecture archived", nil)
}
