package enums

type ErrorCode string

const (
	ErrCodeInternalError   ErrorCode = "INTERNAL_ERROR"
	ErrCodeValidationError ErrorCode = "VALIDATION_ERROR"
	ErrCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrCodeBadRequest      ErrorCode = "BAD_REQUEST"
	ErrCodeResponse        ErrorCode = "ERR_CODE"
)

// Return the status code of the error
func (e ErrorCode) StatusCode() int {
	return map[ErrorCode]int{
		ErrCodeInternalError:   500, // Internal Server Error
		ErrCodeValidationError: 422, // Unprocessable Entity (more appropriate for validation)
		ErrCodeUnauthorized:    401, // Unauthorized
		ErrCodeForbidden:       403, // Forbidden
		ErrCodeNotFound:        404, // Not Found
		ErrCodeBadRequest:      400, // Bad Request
	}[e]
}

// Get error message for the error code
func (e ErrorCode) Message() string {
	return map[ErrorCode]string{
		ErrCodeInternalError:   "Internal Server Error",
		ErrCodeValidationError: "Validation failed",
		ErrCodeUnauthorized:    "Unauthorized Access",
		ErrCodeForbidden:       "Forbidden",
		ErrCodeNotFound:        "Resource Not Found",
		ErrCodeBadRequest:      "Bad Request",
	}[e]
}
