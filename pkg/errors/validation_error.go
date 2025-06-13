package errors

import "taskgo/pkg/enums"

// ValidationError represents a validation error with a map of field errors
type ValidationError struct {
	Errors    map[string]any
	ErrorCode string
}

// ValidationErrorType implements the error interface for validation errors (error message)
func (v *ValidationError) Error() string {
	return "Validation failed"
}

// Public Error returns a user-friendly error message for public consumption
func (v *ValidationError) PublicError() string {
	return "Validation failed"
}

// ValidationError creates a new validation error with the given errors
func NewValidationError(errors map[string]any) *ValidationError {
	return &ValidationError{Errors: errors, ErrorCode: string(enums.ErrCodeValidationError)}
}

// AsValidationError tries to cast err to *ValidationError.
// Returns true and sets target if successful, otherwise false.
func AsValidationError(err error) (*ValidationError, bool) {
	ve, ok := err.(*ValidationError)
	return ve, ok
}
