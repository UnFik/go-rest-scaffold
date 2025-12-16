package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func NewValidator(viper *viper.Viper) *validator.Validate {
	return validator.New()
}

// FormatValidationErrors formats validator errors into readable messages
func FormatValidationErrors(err error) string {
	if err == nil {
		return ""
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string
		for _, e := range validationErrors {
			errorMessages = append(errorMessages, formatValidationError(e))
		}
		return strings.Join(errorMessages, "; ")
	}

	return err.Error()
}

func formatValidationError(e validator.FieldError) string {
	field := toSnakeCase(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, e.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

func toSnakeCase(s string) string {
	var result string

	for i, v := range s {
		if i > 0 && v >= 'A' && v <= 'Z' {
			result += "_"
		}

		result += string(v)
	}

	return strings.ToLower(result)
}
