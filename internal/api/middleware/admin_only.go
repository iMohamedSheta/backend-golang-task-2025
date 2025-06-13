package middleware

import (
	"taskgo/pkg/response"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the user is an admin
		if !c.GetBool("is_admin") {
			response.ErrorJson(c, "Forbidden: Admin access required", "authorization_error", 403)
			c.Abort()
			return
		}
		c.Next()
	}
}
