package response

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"success"`
	Data    any         `json:"data,omitempty"`
	Errors  FieldErrors `json:"errors,omitempty"`
}

type FieldErrors map[string]string

func JSON(c *fiber.Ctx, httpStatus int, success bool, message string, data any) error {
	return c.Status(httpStatus).JSON(&Response{
		Success: success,
		Message: message,
		Data:    data,
	})
}
