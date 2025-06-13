package rules

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func EgyptianPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	re := regexp.MustCompile(`^01[0125][0-9]{8}$`)
	return re.MatchString(phone)
}
