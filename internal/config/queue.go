package config

func init() {
	Register(queueConfig)
}

// Redis configuration for the application
func queueConfig(cfg *Config) {
	cfg.Set("queue", map[string]any{
		// Default Redis connection to use, we will use redis connection (jobs) if it's redis
		"default":  Env("QUEUE_DEFAULT", "redis"),
		"enabled":  Env("QUEUE_ACTIVE", true),
		"required": Env("QUEUE_REQUIRED", false), // Required for app to run
		"consumer": map[string]any{
			// Worker concurrency settings
			"concurrency": Env("QUEUE_CONSUMER_CONCURRENCY", 10),

			// Queue priorities
			"queues": map[string]any{
				"critical":               Env("QUEUE_PRIORITY_CRITICAL", 6),
				"default":                Env("QUEUE_PRIORITY_DEFAULT", 3),
				"low":                    Env("QUEUE_PRIORITY_LOW", 1),
				"payments":               3,
				"inventory_check":        3,
				"order_processing_chain": 6,
				"notifications":          3,
			},

			// Retry configuration
			"retry": map[string]any{
				"max_attempts": Env("QUEUE_MAX_RETRY_ATTEMPTS", 3),
				"delay":        Env("QUEUE_RETRY_DELAY", "15s"), // delay between retries
			},

			// Health check configuration
			"health_check": map[string]any{
				"enabled":  Env("QUEUE_HEALTH_CHECK_ENABLED", true),
				"interval": Env("QUEUE_HEALTH_CHECK_INTERVAL", "30s"),
			},

			// Logging configuration for tasks
			"logging": map[string]any{
				"log_level":        Env("QUEUE_LOG_LEVEL", "info"),
				"log_failed_tasks": Env("QUEUE_LOG_FAILED_TASKS", true),
				"log_success":      Env("QUEUE_LOG_SUCCESS", true),
			},
		},
	})
}
