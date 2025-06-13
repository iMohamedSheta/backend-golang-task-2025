package config

import (
	"time"
)

func init() {
	Register(appConfig)
}

func appConfig() {
	App.Set("app", map[string]any{
		"name":              Env("APP_NAME", "GoCrudRestApi"),
		"url":               Env("APP_URL", "localhost"),
		"port":              Env("APP_PORT", "8080"),
		"shutdown_timeout":  20 * time.Second,
		"env":               Env("APP_ENV", "dev"),
		"debug":             Env("APP_DEBUG", true),
		"secret":            Env("APP_SECRET", "hxdCTfhtkyJBVE01k8vvtaMHbzTmr401QqGl1111"),
		"global_rate_limit": Env("APP_GLOBAL_RATE_LIMIT", "100-M"), // 100 requests per minute
		"admins": []map[string]any{
			{
				"email":      "admin@admin.com",
				"password":   "123456789",
				"first_name": "Admin",
				"last_name":  "User",
				"is_active":  true,
			},
		},
	})
}
