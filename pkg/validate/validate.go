package validate

import (
	"fmt"
	"taskgo/pkg/contracts"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func Validate(data map[string]interface{}, rules map[string]string, messages map[string]string) (bool, map[string]string) {
	validate := Validator()
	errors := make(map[string]string)

	for field, rule := range rules {
		err := validate.Var(data[field], rule)
		if err != nil {
			// Get the first error for the field
			if validationErrs, ok := err.(validator.ValidationErrors); ok && len(validationErrs) > 0 {
				// Get the custom error message if available
				customMessage := fmt.Sprintf("Field %s failed on '%s'", field, validationErrs[0].Tag())

				// Check if custom message exists in the messages map
				if msg, exists := messages[field]; exists {
					customMessage = msg
				}

				errors[field] = customMessage
			} else {
				errors[field] = "Validation failed"
			}
		}
	}

	return len(errors) == 0, errors
}

func ValidateRequest(r contracts.Validatable) (bool, map[string]any) {
	validate := Validator()

	if err := validate.Struct(r); err != nil {
		errors := make(map[string]any)
		messages := r.Messages()

		for _, valErr := range err.(validator.ValidationErrors) {
			// Use the JSON tag name instead of struct field name
			fieldName := valErr.Field()
			if jsonTag := valErr.Field(); jsonTag != "" {
				fieldName = toSnakeCase(valErr.Field())
			}

			key := fmt.Sprintf("%s.%s", fieldName, valErr.Tag())
			if msg, ok := messages[key]; ok {
				errors[fieldName] = msg
			} else {
				errors[fieldName] = valErr.Error()
			}
		}
		return false, errors
	}

	return true, nil
}

func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if unicode.IsUpper(r) && i != 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}
