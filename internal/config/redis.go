package config

func init() {
	Register(redisConfig)
}

// Redis configuration for the application
func redisConfig() {
	App.Set("redis", map[string]any{
		// Default Redis connection to use
		"default": Env("REDIS_DEFAULT", "default"),

		// Global Redis options
		"options": map[string]any{
			"cluster": Env("REDIS_CLUSTER", "redis"),
			"prefix":  Env("REDIS_PREFIX", nil),
		},

		// Redis connections
		"connections": map[string]any{
			"default": map[string]any{
				"url":       Env("REDIS_URL", nil),
				"host":      Env("REDIS_HOST", "127.0.0.1"),
				"password":  Env("REDIS_PASSWORD", nil),
				"port":      Env("REDIS_PORT", 6379),
				"database":  Env("REDIS_DB", 10),
				"is_active": Env("REDIS_DEFAULT_ACTIVE", true),
				"pool_size": Env("REDIS_POOL_SIZE", 0),
				"timeout":   Env("REDIS_TIMEOUT", ""),
			},
			"jobs": map[string]any{
				"url":       Env("REDIS_JOB_URL", nil),
				"host":      Env("REDIS_JOB_HOST", "127.0.0.1"),
				"password":  Env("REDIS_JOB_PASSWORD", nil),
				"port":      Env("REDIS_JOB_PORT", 6379),
				"database":  Env("REDIS_JOB_DB", 9),
				"is_active": Env("REDIS_JOB_ACTIVE", true),
				"pool_size": Env("REDIS_JOB_POOL_SIZE", 0),
				"timeout":   Env("REDIS_JOB_TIMEOUT", ""),
			},
		},
	})
}
