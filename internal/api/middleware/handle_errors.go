package middleware

import (
	"runtime"
	"strings"
	"taskgo/internal/enums"
	"taskgo/pkg/errors"
	"taskgo/pkg/logger"
	"taskgo/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HandlerFuncWithError func(*gin.Context) error

func HandleErrors(handler HandlerFuncWithError) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := handler(c); err != nil {
			// Log the error first for debugging/monitoring
			globalErrorHandler(c, err)

			// Handle specific error types
			switch e := err.(type) {
			case *errors.ValidationError:
				validationErrorHandler(c, e)
			case *errors.BadRequestError:
				badRequestErrorHandler(c, e)
			case *errors.BadRequestBindingError:
				badRequestBindingErrorHandler(c, e)
			case *errors.UnAuthorizedError:
				unAuthorizedErrorHandler(c, e)
			case *errors.NotFoundError:
				notFoundErrorHandler(c, e)
			case *errors.ServerError:
				serverErrorHandler(c, e)
			default:
				unknownErrorHandler(c, err)
			}

			// Prevent further middleware execution if it's an error
			c.Abort()
			return
		}
	}
}

func globalErrorHandler(c *gin.Context, err error) {
	appRequestErrorLogger(c, err)
}

func validationErrorHandler(c *gin.Context, err *errors.ValidationError) {
	response.ValidationErrorJson(c, err)
}

func badRequestErrorHandler(c *gin.Context, err *errors.BadRequestError) {
	response.BadRequestErrorJson(c, err)
}

func badRequestBindingErrorHandler(c *gin.Context, err *errors.BadRequestBindingError) {
	response.BadRequestBindingJson(c, err)
}

func notFoundErrorHandler(c *gin.Context, err *errors.NotFoundError) {
	response.NotFoundJson(c, err)
}

func unAuthorizedErrorHandler(c *gin.Context, err *errors.UnAuthorizedError) {
	response.UnauthorizedJson(c, err)
}

func serverErrorHandler(c *gin.Context, err *errors.ServerError) {
	response.ServerErrorJson(c, err)
}

func unknownErrorHandler(c *gin.Context, err error) {
	response.ServerErrorJson(c, errors.NewServerError("Internal Server Error", "An unexpected error occurred", err))
}

// Helper function to get error type as string
func getErrorType(err error) string {
	switch err.(type) {
	case *errors.ValidationError:
		return "ValidationError"
	case *errors.BadRequestError:
		return "BadRequestError"
	case *errors.UnAuthorizedError:
		return "UnAuthorizedError"
	case *errors.NotFoundError:
		return "NotFoundError"
	case *errors.ServerError:
		return "ServerError"
	default:
		return "UnknownError"
	}
}

func appRequestErrorLogger(c *gin.Context, err error) {
	// Get stack trace
	stackBuf := make([]byte, 1024*8) // 8KB buffer
	stackSize := runtime.Stack(stackBuf, false)
	stack := string(stackBuf[:stackSize])

	// Get caller information
	_, file, line, ok := runtime.Caller(2)
	var caller string
	if ok {
		// Extract just the filename
		parts := strings.Split(file, "/")
		if len(parts) > 0 {
			caller = parts[len(parts)-1]
		}
	}

	requestID := c.GetString("request_id")
	if requestID == "" {
		requestID = c.GetHeader("X-Request-ID")
	}

	userID := c.GetString(string(enums.ContextKeyAuthId))
	if userID == "" {
		userID = "anonymous"
	}

	// Log comprehensive error information
	logger.Log().Error("Request error occurred",
		// Error
		zap.String("error", err.Error()),
		zap.String("error_type", getErrorType(err)),

		// Request context
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("raw_query", c.Request.URL.RawQuery),
		zap.String("user_agent", c.GetHeader("User-Agent")),
		zap.String("client_ip", c.ClientIP()),
		zap.String("request_id", requestID),
		zap.String("user_id", userID),

		// Request headers
		zap.Any("headers", sanitizeHeaders(c.Request.Header)),

		// Timing
		zap.String("timestamp", time.Now().Format("2006-01-02 15:04:05")),

		// Code location
		zap.String("file", caller),
		zap.Int("line", line),

		// Stack trace (you might want to make this conditional based on log level)
		zap.String("stack_trace", stack),

		// Additional context if available
		zap.Any("request_body_size", c.Request.ContentLength),
		zap.String("content_type", c.GetHeader("Content-Type")),
		zap.String("accept", c.GetHeader("Accept")),
		zap.String("referer", c.GetHeader("Referer")),
	)

	// Also log request parameters if they exist
	if len(c.Params) > 0 {
		params := make(map[string]string)
		for _, param := range c.Params {
			params[param.Key] = param.Value
		}
		logger.Log().Info("Request parameters",
			zap.Any("params", params),
			zap.String("request_id", requestID),
		)
	}

	// Log query parameters if they exist
	if len(c.Request.URL.Query()) > 0 {
		logger.Log().Info("Query parameters",
			zap.Any("query_params", sanitizeQueryParams(c.Request.URL.Query())),
			zap.String("request_id", requestID),
		)
	}
}

// Helper function to sanitize headers
func sanitizeHeaders(headers map[string][]string) map[string][]string {
	sanitized := make(map[string][]string)
	sensitiveHeaders := map[string]bool{
		"authorization": true,
		"cookie":        true,
		"x-api-key":     true,
		"x-auth-token":  true,
	}

	for key, values := range headers {
		lowerKey := strings.ToLower(key)
		if sensitiveHeaders[lowerKey] {
			sanitized[key] = []string{"[REDACTED]"}
		} else {
			sanitized[key] = values
		}
	}
	return sanitized
}

// Helper function to sanitize query parameters
func sanitizeQueryParams(params map[string][]string) map[string][]string {
	sanitized := make(map[string][]string)
	sensitiveParams := map[string]bool{
		"password": true,
		"api_key":  true,
		"token":    true,
		"secret":   true,
	}

	for key, values := range params {
		lowerKey := strings.ToLower(key)
		if sensitiveParams[lowerKey] {
			sanitized[key] = []string{"[REDACTED]"}
		} else {
			sanitized[key] = values
		}
	}
	return sanitized
}
