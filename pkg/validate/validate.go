package validate

import (
	"fmt"
	"taskgo/pkg/contracts"
	"taskgo/pkg/errors"
	"taskgo/pkg/reflect"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New(v *validator.Validate) *Validator {
	return &Validator{validate: v}
}

func (v *Validator) Validate(data map[string]interface{}, rules map[string]string, messages map[string]string) (bool, map[string]string) {
	validate := v.validate
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

func (v *Validator) validateRequest(r contracts.Validatable) (bool, map[string]any, error) {
	validate := v.validate

	if err := validate.Struct(r); err != nil {
		// Check if the the type of err is not a validationErrors type then return it as serverError (InvalidValidationError)
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return false, nil, err
		}

		errors := make(map[string]any)
		messages := r.Messages()

		for _, valErr := range validationErrors {
			// Get the actual JSON tag name from the struct field
			fieldName := reflect.GetJsonFieldName(r, valErr.Field())

			// Skip fields that are excluded from JSON (json:"-")
			if fieldName == "" {
				continue
			}

			// Create the message key
			key := fmt.Sprintf("%s.%s", fieldName, valErr.Tag())

			// Use custom message if available, otherwise use default error
			if msg, ok := messages[key]; ok {
				errors[fieldName] = msg
			} else {
				errors[fieldName] = valErr.Error()
			}
		}
		return false, errors, nil
	}

	return true, nil, nil
}

// ValidateRequest validates the request and returns an error if it fails.
func (v *Validator) ValidateRequest(req contracts.Validatable) error {
	valid, validationErrors, err := v.validateRequest(req)
	if err != nil {
		return errors.NewServerError("", "Internal server error: validation fail to process for the login.", err)
	}

	if !valid {
		return errors.NewValidationError(validationErrors)
	}

	return nil
}
