package validator

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/pkg/apperr"
	"github.com/victorive/klaro-backend/internal/pkg/response"
)

var globalValidator = validator.New()

func Validate(s interface{}) error {
	return globalValidator.Struct(s)
}

func ParseAndValidate(c *fiber.Ctx, dst interface{}) error {
	if err := c.BodyParser(dst); err != nil {
		return apperr.New("invalid request body", http.StatusBadRequest)
	}

	if err := globalValidator.Struct(dst); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errs := make([]response.ValidationError, 0, len(ve))
			for _, e := range ve {
				errs = append(errs, response.ValidationError{
					Field:   toCamelCase(e.Field()),
					Message: getValidationMessage(e),
				})
			}
			return apperr.New("validation failed", http.StatusUnprocessableEntity).WithData(errs)
		}
		return apperr.New("validation failed", http.StatusBadRequest)
	}

	return nil
}

func getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short (minimum: " + e.Param() + ")"
	case "max":
		return "Value is too long (maximum: " + e.Param() + ")"
	case "len":
		return "Value must be exactly " + e.Param() + " characters"
	case "password":
		return "Password must contain at least one uppercase letter, one number, and be at least 8 characters long"
	case "oneof":
		return "Value must be one of: " + e.Param()
	case "gte":
		return "Value must be greater than or equal to " + e.Param()
	case "lte":
		return "Value must be less than or equal to " + e.Param()
	case "numeric":
		return "Value must be numeric"
	case "alpha":
		return "Value must contain only letters"
	case "alphanum":
		return "Value must contain only letters and numbers"
	case "url":
		return "Invalid URL"
	case "uuid":
		return "Invalid UUID"
	default:
		return "Validation failed: " + e.Tag()
	}
}

func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 1 {
		return s
	}
	for i := 1; i < len(parts); i++ {
		parts[i] = strings.ToUpper(string(parts[i][0])) + parts[i][1:]
	}
	return strings.Join(parts, "")
}
