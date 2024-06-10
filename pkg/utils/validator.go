package utils

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	Validator *validator.Validate = validator.New(validator.WithRequiredStructEnabled())
)

func IsStrongPassword(f validator.FieldLevel) bool {
	value := f.Field().String()

	letterPattern := "[a-zA-Z]"
	numberPattern := "[0-9]"

	letterRegex := regexp.MustCompile(letterPattern)
	digitRegex := regexp.MustCompile(numberPattern)

	return letterRegex.MatchString(value) && digitRegex.MatchString(value)
}

func HumanFriendlyErrorMessage(tag string, param string) string {
	switch tag {
	case "max":
		return fmt.Sprintf("must be less than %s", param)
	case "min":
		return fmt.Sprintf("must be greater than %s", param)
	case "email":
		return "email must be a valid email"
	case "required":
		return "value cannot be null or empty"
	case "alpha":
		return "value must contain only letters"
	case "is-strong":
		return "value must contain letters and numbers"
	default:
		return ""
	}
}

func init() {
	Validator.RegisterValidation("is-strong", IsStrongPassword)
}
