package response

import (
	"encoding/json"
	"fmt"

	"taskgo/pkg/enums"
	pkgErrors "taskgo/pkg/errors"

	"github.com/gin-gonic/gin"
)

// BaseResponse represents the basic response structure
// @Description Base response structure
type BaseResponse struct {
	Message   string `json:"message" example:"Success"`
	ErrorCode string `json:"error_code,omitempty" example:""`
}

// SuccessResponse represents a successful response with data
// @Description Successful response with data
type SuccessResponse struct {
	Message string `json:"message" example:"Success"`
	Data    any    `json:"data,omitempty"`
}

// ErrorResponse represents an error response
// @Description Error response structure
type ErrorResponse struct {
	Message   string `json:"message" example:"Error occurred"`
	ErrorCode string `json:"error_code" example:"ERR_CODE"`
}

// ValidationErrorResponse represents a validation error response
// @Description Validation error response with field errors
type ValidationErrorResponse struct {
	Message   string         `json:"message" example:"Validation failed"`
	ErrorCode string         `json:"error_code" example:"VALIDATION_ERROR"`
	Data      map[string]any `json:"data,omitempty"`
}

// ServerErrorResponse represents an internal server error response
// @Description Internal server error response
type ServerErrorResponse struct {
	Message   string `json:"message" example:"Internal Server Error"`
	ErrorCode string `json:"error_code" example:"INTERNAL_ERROR"`
}

// NotFoundResponse represents a not found error response
// @Description Resource not found error response
type NotFoundResponse struct {
	Message   string `json:"message" example:"Resource not found"`
	ErrorCode string `json:"error_code" example:"NOT_FOUND"`
}

// UnauthorizedResponse represents an unauthorized error response
// @Description Unauthorized access error response
type UnauthorizedResponse struct {
	Message   string `json:"message" example:"Unauthorized access"`
	ErrorCode string `json:"error_code" example:"UNAUTHORIZED"`
}

// BadRequestResponse represents a bad request error response
// @Description Bad request error response
type BadRequestResponse struct {
	Message   string `json:"message" example:"Bad request"`
	ErrorCode string `json:"error_code" example:"BAD_REQUEST"`
}

// Standard JSON response
func Json(c *gin.Context, message string, data any, code int) {
	resp := &SuccessResponse{
		Message: message,
		Data:    data,
	}
	c.JSON(code, resp)
}

// Generic error response
func ErrorJson(c *gin.Context, message, errorCode string, code int) {
	resp := &ErrorResponse{
		Message:   message,
		ErrorCode: errorCode,
	}
	c.JSON(code, resp)
}

// Internal server error
func ServerErrorJson(c *gin.Context, err *pkgErrors.ServerError) {
	msg := "Internal Server Error"
	if err != nil {
		msg = err.PublicError()
	}

	resp := &ServerErrorResponse{
		Message:   msg,
		ErrorCode: string(enums.ErrCodeInternalError),
	}
	c.JSON(enums.ErrCodeInternalError.StatusCode(), resp)
}

// Validation error response
func ValidationErrorJson(c *gin.Context, validationError *pkgErrors.ValidationError) {
	resp := &ValidationErrorResponse{
		Message:   validationError.PublicError(),
		ErrorCode: validationError.ErrorCode,
		Data:      validationError.Errors,
	}
	c.JSON(enums.ErrCodeValidationError.StatusCode(), resp)
}

// Bad request response
func BadRequestErrorJson(c *gin.Context, badRequestError *pkgErrors.BadRequestError) {
	resp := &BadRequestResponse{
		Message:   badRequestError.PublicError(),
		ErrorCode: string(enums.ErrCodeBadRequest),
	}
	c.JSON(enums.ErrCodeBadRequest.StatusCode(), resp)
}

// Not found response
func NotFoundJson(c *gin.Context, err *pkgErrors.NotFoundError) {
	resp := &NotFoundResponse{
		Message:   err.PublicError(),
		ErrorCode: string(enums.ErrCodeNotFound),
	}
	c.JSON(enums.ErrCodeNotFound.StatusCode(), resp)
}

// Unauthorized response
func UnauthorizedJson(c *gin.Context, err *pkgErrors.UnAuthorizedError) {
	resp := &UnauthorizedResponse{
		Message:   err.PublicError(),
		ErrorCode: string(enums.ErrCodeUnauthorized),
	}
	c.JSON(enums.ErrCodeUnauthorized.StatusCode(), resp)
}

// Bad request json binding response
func BadRequestBindingJson(c *gin.Context, badRequestBindingErr *pkgErrors.BadRequestBindingError) {
	var errorMessage string
	var fieldErrors map[string]any

	err := badRequestBindingErr.Err
	if err == nil {
		err = badRequestBindingErr
	}

	// Check if it's a JSON syntax error
	if jsonErr, ok := err.(*json.SyntaxError); ok {
		errorMessage = fmt.Sprintf("Invalid JSON syntax at position %d", jsonErr.Offset)
		resp := &BadRequestResponse{
			Message:   errorMessage,
			ErrorCode: string(enums.ErrCodeBadRequest),
		}
		c.JSON(enums.ErrCodeBadRequest.StatusCode(), resp)
		return
	}

	fieldErrors = make(map[string]any)

	// if it's type error handle it as validation error
	if jsonErr, ok := err.(*json.UnmarshalTypeError); ok {
		field := jsonErr.Field
		if field == "" {
			field = "unknown"
		}
		fieldErrors[field] = fmt.Sprintf("Expected %s but received %s",
			jsonErr.Type.String(), jsonErr.Value)
	} else {
		resp := &BadRequestResponse{
			Message:   "Invalid request body",
			ErrorCode: string(enums.ErrCodeBadRequest),
		}
		c.JSON(enums.ErrCodeBadRequest.StatusCode(), resp)
		return
	}

	validationError := &pkgErrors.ValidationError{
		ErrorCode: string(enums.ErrCodeValidationError),
		Errors:    fieldErrors,
	}

	ValidationErrorJson(c, validationError)
}
