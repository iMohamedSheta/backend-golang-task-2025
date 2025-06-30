package middleware

import (
	"fmt"
	"runtime/debug"
	"taskgo/internal/deps"
	"taskgo/pkg/errors"
	"taskgo/pkg/response"

	"github.com/gin-gonic/gin"
)

func RecoveryWithLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				deps.Log().Log().Error(fmt.Sprintf(
					"\n🚨 Panic Recovered 🚨\nMethod: %s\nEndpoint: %s\nError: %v\n\nStack Trace:\n%s\n\n",
					c.Request.Method,
					c.Request.URL.Path,
					r,
					debug.Stack(),
				))

				response.ServerErrorJson(c, errors.NewServerError("", "Panic error", fmt.Errorf("%v", r)))
			}
		}()
		c.Next()
	}
}
