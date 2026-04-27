package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/victorive/klaro-backend/internal/pkg/apperr"
	"github.com/victorive/klaro-backend/internal/pkg/response"
)

var globalValidator = func() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
		}
		
		if name == "-" || name == "" {
			return fld.Name
		}
		return name
	})
	
	return v
}()

func Validate(s interface{}) error {
	if err := globalValidator.Struct(s); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			return buildValidationError(ve)
		}
		return apperr.New("invalid input", http.StatusUnprocessableEntity)
	}

	return nil
}

func ParseAndValidate(c *fiber.Ctx, dst interface{}) error {
	if err := c.BodyParser(dst); err != nil {
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &typeErr) {
			fields := response.FieldErrors{
				typeErr.Field: fmt.Sprintf("expected %s, got %s", typeErr.Type.String(), typeErr.Value),
			}
			return apperr.New("invalid input", http.StatusUnprocessableEntity).WithData(fields)
		}
		return apperr.New("invalid request body", http.StatusBadRequest)
	}

	if err := globalValidator.Struct(dst); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			return buildValidationError(ve)
		}
		return apperr.New("invalid input", http.StatusBadRequest)
	}

	return nil
}

func buildValidationError(ve validator.ValidationErrors) error {
	fields := make(response.FieldErrors, len(ve))
	for _, e := range ve {
		fields[e.Field()] = getValidationMessage(e)
	}

	return apperr.New("invalid input", http.StatusUnprocessableEntity).WithData(fields)
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
	case "ltefield":
		return "Must be before " + e.Param()
	case "gtefield":
		return "Must be after " + e.Param()
	case "datetime":
		return "Invalid date format, expected: " + e.Param()
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
