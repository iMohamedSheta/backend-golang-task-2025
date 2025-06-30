package middleware

import (
	"taskgo/internal/deps"
	"taskgo/pkg/utils"

	"github.com/gin-gonic/gin"
)

const (
	defaultOrigin         = "*"
	defaultAllowedHeaders = "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"
	defaultAllowedMethods = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawCors, err := deps.Config().Get("cors")
		if err != nil {
			c.Next()
			return
		}

		cors := rawCors.(map[string]any)

		reqOrigin := c.Request.Header.Get("Origin")
		allowedOrigins := utils.ToArrayOfStrings(cors["origin"], []string{defaultOrigin})

		allowOrigin := ""
		for _, o := range allowedOrigins {
			if o == "*" || o == reqOrigin {
				allowOrigin = o
				break
			}
		}

		if allowOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		}

		methods := utils.ToCSV(cors["methods"], defaultAllowedMethods)
		headers := utils.ToCSV(cors["allowed_headers"], defaultAllowedHeaders)

		c.Writer.Header().Set("Access-Control-Allow-Methods", methods)
		c.Writer.Header().Set("Access-Control-Allow-Headers", headers)

		if exposed := utils.ToCSV(cors["exposed_headers"], ""); exposed != "" {
			c.Writer.Header().Set("Access-Control-Expose-Headers", exposed)
		}

		if maxAge, ok := cors["max_age"].(string); ok {
			c.Writer.Header().Set("Access-Control-Max-Age", maxAge)
		}

		// Only set credentials if origin is not "*" and credentials are set to true in the config (cors policy)
		if cred, ok := cors["credentials"].(bool); ok && cred && allowOrigin != "*" {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
