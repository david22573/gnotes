package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator implements echo.Validator interface
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator creates a new custom validator and registers custom validation functions
func NewCustomValidator() *CustomValidator {
	v := validator.New()

	// Register custom validation functions
	if err := v.RegisterValidation("username", validateUsername); err != nil {
		panic("failed to register username validator: " + err.Error())
	}

	return &CustomValidator{
		validator: v,
	}
}

// Validate implements echo.Validator interface
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(400, formatValidationError(err))
	}
	return nil
}

// validateUsername ensures username contains only allowed characters
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	// Allow alphanumeric characters and underscores
	for _, char := range username {
		if !isAlphanumeric(char) && char != '_' {
			return false
		}
	}
	return true
}

func isAlphanumeric(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9')
}

// formatValidationError converts validator errors to user-friendly messages
func formatValidationError(err error) string {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			switch e.Tag() {
			case "required":
				return field(e.Field()) + " is required"
			case "email":
				return field(e.Field()) + " must be a valid email address"
			case "min":
				return field(e.Field()) + " must be at least " + e.Param() + " characters long"
			case "max":
				return field(e.Field()) + " must not exceed " + e.Param() + " characters"
			case "username":
				return field(e.Field()) + " must contain only letters, numbers, and underscores"
			}
		}
	}
	return "Validation failed"
}

func field(s string) string {
	return s
}
