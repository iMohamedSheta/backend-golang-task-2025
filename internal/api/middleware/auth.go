package middleware

import (
	"fmt"
	"strings"

	"taskgo/internal/enums"
	"taskgo/internal/helpers"
	"taskgo/pkg/logger"
	"taskgo/pkg/response"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.ErrorJson(c, "Missing or invalid Authorization header", "missing_auth_header", http.StatusUnauthorized)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := helpers.ValidateAuthToken(enums.AccessToken, tokenStr)
		if err != nil {
			response.ErrorJson(c, err.Error(), "invalid_token", http.StatusUnauthorized)
			c.Abort()
			return
		}

		userID, err := claims.GetSubject()
		if err != nil || userID == "" {
			logger.Log().Error(fmt.Sprintf("Missing or invalid %s in token", string(enums.ContextKeyAuthId)), zap.Any("claims", claims))
			response.ErrorJson(c, fmt.Sprintf("Missing or invalid %s in token", string(enums.ContextKeyAuthId)), "invalid_token", http.StatusUnauthorized)
			c.Abort()
			return
		}

		role, err := claims.GetRole()

		if err != nil || role == "" {
			logger.Log().Error("Invalid or missing role in token", zap.Any("claims", claims))
			response.ErrorJson(c, "Invalid role in token", "invalid_token", http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set(string(enums.ContextKeyAuthId), userID)
		c.Set(string(enums.ContextKeyRole), role)

		if role == enums.RoleAdmin {
			c.Set("is_admin", true)
		} else {
			c.Set("is_admin", false)
		}

		c.Next()
	}
}
