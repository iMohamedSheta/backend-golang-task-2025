package middleware

import (
	"time"

	"taskgo/internal/deps"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Logger is a Gin middleware for logging HTTP requests using zap
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate or get request ID for tracing
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
			c.Request.Header.Set("X-Request-ID", reqID)
		}
		c.Writer.Header().Set("X-Request-ID", reqID)

		// Process request
		c.Next()

		// request termination (after the request has been processed)
		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		userAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		contentLength := c.Request.ContentLength
		host := c.Request.Host
		protocol := c.Request.Proto
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()
		errorCount := len(c.Errors)

		log := deps.Log().Channel("request_log").With(
			zap.Int("status", status),
			zap.String("latency", latency.String()),
			zap.String("client_ip", clientIP),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("user_agent", userAgent),
			zap.String("referer", referer),
			zap.Int64("content_length", contentLength),
			zap.String("host", host),
			zap.String("protocol", protocol),
			zap.String("request_id", reqID),
			zap.Int("error_count", errorCount),
		)

		if errorMessage != "" {
			log.Error(errorMessage)
		} else {
			log.Info("Handled request")
		}
	}
}
