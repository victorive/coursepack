package file

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/victorive/coursepack/internal/permissions"
	"github.com/victorive/coursepack/internal/pkg/apperr"
	"github.com/victorive/coursepack/internal/pkg/response"
	"github.com/victorive/coursepack/internal/pkg/validator"
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
	actor, err := permissions.ActorFrom(c.UserContext())
	if err != nil {
		return err
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

	head := make([]byte, sniffLen)
	n, err := io.ReadFull(f, head)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return apperr.New("failed to read file", http.StatusInternalServerError).Wrap(err)
	}
	head = head[:n]

	mimeType := http.DetectContentType(head)
	if isUnsafeUploadType(mimeType) {
		return apperr.New("unsupported file type", http.StatusUnsupportedMediaType)
	}

	reader := io.MultiReader(bytes.NewReader(head), f)
	resp, err := h.svc.Upload(c.UserContext(), actor, req, reader, fh.Size, mimeType, fh.Filename)
	if err != nil {
		return err
	}
	return response.JSON(c, http.StatusCreated, true, "file uploaded", resp)
}

const sniffLen = 512

// isUnsafeUploadType rejects markup that browsers would execute when the
// signed URL is opened directly.
func isUnsafeUploadType(mimeType string) bool {
	switch strings.Split(mimeType, ";")[0] {
	case "text/html", "text/xml", "application/xml", "application/xhtml+xml",
		"image/svg+xml", "text/javascript", "application/javascript":
		return true
	}
	return false
}

// ListByClassCourse godoc
// @Summary      List files by class course
// @Description  Returns all non-archived files belonging to a class course.
// @Tags         Files
// @Produce      json
// @Param        id path string true "Class Course UUID"
// @Success      200 {object} response.Response{data=[]FileResponse}
// @Failure      400,401 {object} response.Response
// @Security     BearerAuth
// @Router       /courses/{id}/files [get]
func (h *Handler) ListByClassCourse(c *fiber.Ctx) error {
	classCourseUUID := c.Params("id")
	actor, err := permissions.ActorFrom(c.UserContext())
	if err != nil {
		return err
	}
	resp, err := h.svc.ListByClassCourse(c.UserContext(), actor, classCourseUUID)
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
	lectureUUID := c.Params("lectureId")
	actor, err := permissions.ActorFrom(c.UserContext())
	if err != nil {
		return err
	}
	resp, err := h.svc.ListByLecture(c.UserContext(), actor, lectureUUID)
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
	fileUUID := c.Params("id")
	actor, err := permissions.ActorFrom(c.UserContext())
	if err != nil {
		return err
	}
	resp, err := h.svc.Get(c.UserContext(), actor, fileUUID)
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
	fileUUID := c.Params("id")
	actor, err := permissions.ActorFrom(c.UserContext())
	if err != nil {
		return err
	}
	url, err := h.svc.GetDownloadURL(c.UserContext(), actor, fileUUID)
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
	actor, err := permissions.ActorFrom(c.UserContext())
	if err != nil {
		return err
	}
	fileUUID := c.Params("id")
	var req ReviewFileRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return err
	}
	resp, err := h.svc.Review(c.UserContext(), actor, fileUUID, req)
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
	actor, err := permissions.ActorFrom(c.UserContext())
	if err != nil {
		return err
	}

	fileUUID := c.Params("id")

	if err := h.svc.Archive(c.UserContext(), actor, fileUUID); err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, true, "file archived", nil)
}
