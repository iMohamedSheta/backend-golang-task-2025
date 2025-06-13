package config

func init() {
	Register(logConfig)
}

// logConfig sets the logging configuration for the application.
func logConfig() {
	App.Set("log", map[string]any{
		"default": "app_log",

		"channels": map[string]any{
			"app_log": map[string]any{
				"driver":   "daily",
				"path":     Env("APP_LOG_PATH", "storage/logs/log.json"),
				"level":    "debug",
				"max_size": 100,
				"max_age":  30, // in days
				"backup":   false,
			},
			"request_log": map[string]any{
				"driver":   "daily",
				"path":     Env("APP_REQUEST_LOG_PATH", "storage/logs/request.json"),
				"level":    "debug",
				"max_size": 100,
				"max_age":  30, // in days
				"backup":   false,
			},
			"inventory_log": map[string]any{
				"driver":   "daily",
				"path":     Env("APP_INVENTORY_LOG_PATH", "storage/logs/inventory.json"),
				"level":    "debug",
				"max_size": 100,
				"max_age":  30, // in days
				"backup":   false,
			},
			"queue_log": map[string]any{
				"driver":   "daily",
				"path":     Env("APP_QUEUE_LOG_PATH", "storage/logs/queue.json"),
				"level":    "debug",
				"max_size": 100,
				"max_age":  30, // in days
				"backup":   false,
			},
		},
	})
}
