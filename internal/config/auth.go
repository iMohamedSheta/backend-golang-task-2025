package config

import "time"

func init() {
	Register(authConfig)
}

// authConfig sets the authentication configuration for the application.
func authConfig() {
	App.Set("auth", map[string]any{
		"default": "jwt",
		"jwt": map[string]any{
			"secret":    "secret",
			"issuer":    Env("APP_NAME", "TaskGo"),
			"audience":  Env("APP_NAME", "TaskGoAudience"),
			"algorithm": "HS256",
			"access_token": map[string]any{
				"expiry": 30 * time.Minute,
			},
			"refresh_token": map[string]any{
				"expiry": 168 * time.Hour, // 7 days
			},
		},
	})
}
