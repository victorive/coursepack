package course

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/domain"
	"github.com/victorive/klaro-backend/internal/middleware"
	"github.com/victorive/klaro-backend/internal/pkg/apperr"
	"github.com/victorive/klaro-backend/internal/pkg/response"
	"github.com/victorive/klaro-backend/internal/pkg/validator"
)

type UserRepository interface {
	FindByUUID(ctx context.Context, uuid string) (*domain.User, error)
}

type Handler struct {
	svc      Service
	userRepo UserRepository
}

func NewHandler(svc Service, userRepo UserRepository) *Handler {
	return &Handler{svc: svc, userRepo: userRepo}
}

// ListCourses godoc
// @Summary      List courses for the current user's university
// @Description  Returns all courses belonging to the authenticated user's university. Optionally filter by department.
// @Tags         Courses
// @Produce      json
// @Param        department_id query int false "Filter by department ID"
// @Success      200 {object} response.Response{data=[]CourseRegistryResponse}
// @Failure      401,422 {object} response.Response
// @Security     BearerAuth
// @Router       /courses [get]
func (h *Handler) ListCourses(c *fiber.Ctx) error {
	actorUUID, ok := middleware.UserUUID(c)
	if !ok {
		return apperr.New("unauthorized", http.StatusUnauthorized)
	}
	u, err := h.userRepo.FindByUUID(c.UserContext(), actorUUID)
	if err != nil || u == nil {
		return apperr.New("user not found", http.StatusUnauthorized)
	}
	if u.UniversityID == nil {
		return apperr.New("user is not associated with a university", http.StatusUnprocessableEntity)
	}
	universityID := *u.UniversityID

	var departmentID *int64
	if s := c.Query("department_id"); s != "" {
		deptID, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fiber.ErrBadRequest
		}
		departmentID = &deptID
	}

	resp, err := h.svc.ListCourses(c.UserContext(), universityID, departmentID)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "courses fetched", resp)
}

// GetCourse godoc
// @Summary      Get a course
// @Description  Returns a single course by ID, scoped to the authenticated user's university.
// @Tags         Courses
// @Produce      json
// @Param        id path int true "Course ID"
// @Success      200 {object} response.Response{data=CourseRegistryResponse}
// @Failure      400,401,404 {object} response.Response
// @Security     BearerAuth
// @Router       /courses/{id} [get]
func (h *Handler) GetCourse(c *fiber.Ctx) error {
	actorUUID, ok := middleware.UserUUID(c)
	if !ok {
		return apperr.New("unauthorized", http.StatusUnauthorized)
	}

	u, err := h.userRepo.FindByUUID(c.UserContext(), actorUUID)
	if err != nil || u == nil {
		return apperr.New("user not found", http.StatusUnauthorized)
	}

	if u.UniversityID == nil {
		return apperr.New("user is not associated with a university", http.StatusUnprocessableEntity)
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}
	resp, err := h.svc.GetCourse(c.UserContext(), id, *u.UniversityID)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "course fetched", resp)
}

// AddClassCourse godoc
// @Summary      Add a course to a class
// @Description  Adds a course to a class for a given semester. Requires courses:manage permission.
// @Tags         Classes
// @Accept       json
// @Produce      json
// @Param        id   path int                  true "Class ID"
// @Param        body body AddClassCourseRequest true "Class course payload"
// @Success      201 {object} response.Response{data=ClassCourseResponse}
// @Failure      400,401,403,409 {object} response.Response
// @Security     BearerAuth
// @Router       /classes/{id}/class-courses [post]
func (h *Handler) AddClassCourse(c *fiber.Ctx) error {
	actorUUID, ok := middleware.UserUUID(c)
	if !ok {
		return apperr.New("unauthorized", http.StatusUnauthorized)
	}
	u, err := h.userRepo.FindByUUID(c.UserContext(), actorUUID)
	if err != nil || u == nil {
		return apperr.New("user not found", http.StatusUnauthorized)
	}

	if u.UniversityID == nil {
		return apperr.New("user is not associated with a university", http.StatusUnprocessableEntity)
	}

	classID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}
	var req AddClassCourseRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	resp, err := h.svc.AddClassCourse(c.UserContext(), actorUUID, classID, *u.UniversityID, req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusCreated, true, "course added to class", resp)
}

// ListClassCourses godoc
// @Summary      List class courses
// @Description  Returns all courses assigned to a class. Optionally filter by semester.
// @Tags         Classes
// @Produce      json
// @Param        id          path int false "Class ID"
// @Param        semester_id query int false "Filter by semester ID"
// @Success      200 {object} response.Response{data=[]ClassCourseResponse}
// @Failure      400,401 {object} response.Response
// @Security     BearerAuth
// @Router       /classes/{id}/class-courses [get]
func (h *Handler) ListClassCourses(c *fiber.Ctx) error {
	classID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}
	var semesterID *int64
	if s := c.Query("semester_id"); s != "" {
		sid, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fiber.ErrBadRequest
		}
		semesterID = &sid
	}

	resp, err := h.svc.ListClassCourses(c.UserContext(), classID, semesterID)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "class courses fetched", resp)
}

// GetClassCourse godoc
// @Summary      Get a class course
// @Description  Returns a single class course by ID.
// @Tags         Class Courses
// @Produce      json
// @Param        id path int true "Class Course ID"
// @Success      200 {object} response.Response{data=ClassCourseResponse}
// @Failure      400,401,404 {object} response.Response
// @Security     BearerAuth
// @Router       /class-courses/{id} [get]
func (h *Handler) GetClassCourse(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	resp, err := h.svc.GetClassCourse(c.UserContext(), id)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, true, "class course fetched", resp)
}

// RemoveClassCourse godoc
// @Summary      Remove a course from a class
// @Description  Removes a class course. Requires courses:manage permission.
// @Tags         Classes
// @Produce      json
// @Param        id            path int true "Class ID"
// @Param        classCourseId path int true "Class Course ID"
// @Success      204
// @Failure      400,401,403,404 {object} response.Response
// @Security     BearerAuth
// @Router       /classes/{id}/class-courses/{classCourseId} [delete]
func (h *Handler) RemoveClassCourse(c *fiber.Ctx) error {
	actorUUID, ok := middleware.UserUUID(c)
	if !ok {
		return apperr.New("unauthorized", http.StatusUnauthorized)
	}
	classID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}
	classCourseID, err := strconv.ParseInt(c.Params("classCourseId"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err := h.svc.RemoveClassCourse(c.UserContext(), actorUUID, classID, classCourseID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// CreateLecture godoc
// @Summary      Create a lecture
// @Description  Creates a new lecture under a class course.
// @Tags         Class Courses
// @Accept       json
// @Produce      json
// @Param        id   path int                   true "Class Course ID"
// @Param        body body CreateLectureRequest  true "Lecture payload"
// @Success      201 {object} response.Response{data=LectureResponse}
// @Failure      400,401,403 {object} response.Response
// @Security     BearerAuth
// @Router       /class-courses/{id}/lectures [post]
func (h *Handler) CreateLecture(c *fiber.Ctx) error {
	actorUUID, ok := middleware.UserUUID(c)
	if !ok {
		return apperr.New("unauthorized", http.StatusUnauthorized)
	}
	classCourseID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	var req CreateLectureRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}

	resp, err := h.svc.CreateLecture(c.UserContext(), actorUUID, classCourseID, req)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusCreated, true, "lecture created", resp)
}

// ListLectures godoc
// @Summary      List lectures for a class course
// @Description  Returns all published lectures for a given class course.
// @Tags         Class Courses
// @Produce      json
// @Param        id path int true "Class Course ID"
// @Success      200 {object} response.Response{data=[]LectureResponse}
// @Failure      400,401 {object} response.Response
// @Security     BearerAuth
// @Router       /class-courses/{id}/lectures [get]
func (h *Handler) ListLectures(c *fiber.Ctx) error {
	classCourseID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}
	resp, err := h.svc.ListLectures(c.UserContext(), classCourseID)
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
// @Param        id   path int                 true "Lecture ID"
// @Param        body body UpdateLectureRequest true "Lecture update payload"
// @Success      200 {object} response.Response{data=LectureResponse}
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
