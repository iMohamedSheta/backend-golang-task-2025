package errors

import "taskgo/pkg/enums"

// ValidationError represents a validation error with a map of field errors
type ValidationError struct {
	BaseError
	Errors map[string]any
}

// ValidationError creates a new validation error with the given errors
func NewValidationError(errors map[string]any) *ValidationError {
	return &ValidationError{
		Errors:    errors,
		BaseError: newBaseError(enums.ErrCodeValidationError, "", "", nil),
	}
}

// AsValidationError tries to cast err to *ValidationError.
// Returns true and sets target if successful, otherwise false.
func AsValidationError(err error) (*ValidationError, bool) {
	ve, ok := err.(*ValidationError)
	return ve, ok
}
