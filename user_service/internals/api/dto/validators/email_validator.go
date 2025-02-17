package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateEmail(fl validator.FieldLevel) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	email := fl.Field().String()
	return regexp.MustCompile(emailRegex).MatchString(email)
}
