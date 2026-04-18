package file

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

// Upload godoc
// @Summary      Upload a file
// @Description  Uploads a file associated with a course or lecture.
// @Tags         Files
// @Accept       multipart/form-data
// @Produce      json
// @Param        file      formData file              true "File to upload"
// @Param        body      formData UploadFileRequest true "Upload metadata"
// @Success      201 {object} response.Response{data=FileResponse}
// @Failure      400,401 {object} response.Response
// @Security     BearerAuth
// @Router       /files [post]
func (h *Handler) Upload(c *fiber.Ctx) error {
	uuid, ok := c.Locals(middleware.LocalsUserId).(string)
	if !ok || uuid == "" {
		return apperr.New("unauthorized", http.StatusUnauthorized)
	}

	var req UploadFileRequest
	if err := c.BodyParser(&req); err != nil {
		return apperr.New("invalid request body", http.StatusBadRequest).Wrap(err)
	}
	if err := validator.Validate(&req); err != nil {
		return err
	}

	fh, err := c.FormFile("file")
	if err != nil {
		return apperr.New("file is required", http.StatusBadRequest)
	}

	f, err := fh.Open()
	if err != nil {
		return apperr.New("failed to read file", http.StatusInternalServerError)
	}
	defer f.Close()

	mimeType := fh.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	resp, err := h.svc.Upload(c.UserContext(), uuid, req, f, fh.Size, mimeType, fh.Filename)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusCreated, true, "file uploaded", resp)
}

// ListByCourse godoc
// @Summary      List files by course
// @Description  Returns all non-archived files belonging to a course.
// @Tags         Files
// @Produce      json
// @Param        courseId path int true "Course ID"
// @Success      200 {object} response.Response{data=[]FileResponse}
// @Failure      400,401 {object} response.Response
// @Security     BearerAuth
// @Router       /courses/{courseId}/files [get]
func (h *Handler) ListByCourse(c *fiber.Ctx) error {
	courseID, err := strconv.ParseInt(c.Params("courseId"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}
	resp, err := h.svc.ListByCourse(c.UserContext(), courseID)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "files fetched", resp)
}

// ListByLecture godoc
// @Summary      List files by lecture
// @Description  Returns all non-archived files belonging to a lecture.
// @Tags         Files
// @Produce      json
// @Param        lectureId path int true "Lecture ID"
// @Success      200 {object} response.Response{data=[]FileResponse}
// @Failure      400,401 {object} response.Response
// @Security     BearerAuth
// @Router       /lectures/{lectureId}/files [get]
func (h *Handler) ListByLecture(c *fiber.Ctx) error {
	lectureID, err := strconv.ParseInt(c.Params("lectureId"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}
	resp, err := h.svc.ListByLecture(c.UserContext(), lectureID)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "files fetched", resp)
}

// GetFile godoc
// @Summary      Get a file
// @Description  Returns a single file by ID.
// @Tags         Files
// @Produce      json
// @Param        id path int true "File ID"
// @Success      200 {object} response.Response{data=FileResponse}
// @Failure      400,401,404 {object} response.Response
// @Security     BearerAuth
// @Router       /files/{id} [get]
func (h *Handler) GetFile(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}
	resp, err := h.svc.Get(c.UserContext(), id)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "file fetched", resp)
}

// GetDownloadURL godoc
// @Summary      Get file download URL
// @Description  Returns a pre-signed or direct URL to download the file.
// @Tags         Files
// @Produce      json
// @Param        id path int true "File ID"
// @Success      200 {object} response.Response{data=map[string]string}
// @Failure      400,401,404 {object} response.Response
// @Security     BearerAuth
// @Router       /files/{id}/download [get]
func (h *Handler) GetDownloadURL(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}
	url, err := h.svc.GetDownloadURL(c.UserContext(), id)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "download URL generated", fiber.Map{"url": url})
}

// Review godoc
// @Summary      Review a file
// @Description  Approve or reject a file upload.
// @Tags         Files
// @Accept       json
// @Produce      json
// @Param        id   path     int               true "File ID"
// @Param        body body     ReviewFileRequest true "Review payload"
// @Success      200 {object} response.Response{data=FileResponse}
// @Failure      400,401,403,404 {object} response.Response
// @Security     BearerAuth
// @Router       /files/{id}/review [patch]
func (h *Handler) Review(c *fiber.Ctx) error {
	reviewerUUID, ok := c.Locals(middleware.LocalsUserId).(string)
	if !ok || reviewerUUID == "" {
		return apperr.New("unauthorized", http.StatusUnauthorized)
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}
	var req ReviewFileRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	resp, err := h.svc.Review(c.UserContext(), id, reviewerUUID, req)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "file reviewed", resp)
}

// Archive godoc
// @Summary      Archive a file
// @Description  Moves a file to the archived state.
// @Tags         Files
// @Produce      json
// @Param        id path int true "File ID"
// @Success      200 {object} response.Response
// @Failure      400,401,403,404 {object} response.Response
// @Security     BearerAuth
// @Router       /files/{id}/archive [patch]
func (h *Handler) Archive(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err := h.svc.Archive(c.UserContext(), id); err != nil {
		return err
	}
	return response.JSON(c, http.StatusOK, true, "file archived", nil)
}
