package middleware

import (
	"encoding/base64"
	"fmt"
	"strings"

	"taskgo/internal/deps"
	"taskgo/internal/enums"
	"taskgo/internal/services"
	"taskgo/pkg/errors"
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
		jwtService := deps.App[*services.JwtService]()
		claims, err := jwtService.ValidateAuthToken(enums.AccessToken, tokenStr)
		if err != nil {
			response.ErrorJson(c, err.Error(), "invalid_token", http.StatusUnauthorized)
			c.Abort()
			return
		}

		log := deps.Log()

		userID, err := claims.GetSubject()
		if err != nil || userID == "" {
			log.Log().Error(fmt.Sprintf("Missing or invalid %s in token", string(enums.ContextKeyAuthId)), zap.Any("claims", claims))
			response.ErrorJson(c, fmt.Sprintf("Missing or invalid %s in token", string(enums.ContextKeyAuthId)), "invalid_token", http.StatusUnauthorized)
			c.Abort()
			return
		}

		role, err := claims.GetRole()

		if err != nil || role == "" {
			log.Log().Error("Invalid or missing role in token", zap.Any("claims", claims))
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

func WebSocketAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		webSocketProtocolHeader := c.GetHeader("Sec-WebSocket-Protocol")

		tokens := strings.Split(webSocketProtocolHeader, ",")
		var authHeader string
		var err error

		if len(tokens) == 2 {
			header := tokens[0]
			authHeader, err = decodeBase64URL(strings.Trim(tokens[1], " "))
			if err != nil || header != "Authorization" {
				response.UnauthorizedJson(c, errors.NewUnAuthorizedError("Invalid or missing token", "Invalid or missing token", err))
				c.Abort()
				return
			}
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.ErrorJson(c, "Missing or invalid Authorization header", "missing_auth_header", http.StatusUnauthorized)
			c.Abort()
			return
		}

		jwtService := deps.App[*services.JwtService]()
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwtService.ValidateAuthToken(enums.AccessToken, tokenStr)
		if err != nil {
			response.UnauthorizedJson(c, errors.NewUnAuthorizedError("Invalid or missing token", "Invalid access token", err))
			c.Abort()
			return
		}

		userID, err := claims.GetSubject()
		if err != nil || userID == "" {
			response.UnauthorizedJson(c, errors.NewUnAuthorizedError("Invalid or missing token", "Invalid or missing user ID in token", err))
			c.Abort()
			return
		}

		role, err := claims.GetRole()
		if err != nil || role == "" {
			response.UnauthorizedJson(c, errors.NewUnAuthorizedError("Invalid or missing token", "Invalid or missing role in token", err))
			c.Abort()
			return
		}

		c.Set(string(enums.ContextKeyAuthId), userID)
		c.Set(string(enums.ContextKeyRole), role)
		c.Set("is_admin", role == enums.RoleAdmin)
		c.Next()
	}
}

func decodeBase64URL(encoded string) (string, error) {
	decodedBytes, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}
