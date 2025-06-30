package config

func init() {
	Register(redisConfig)
}

// Redis configuration for the application
func redisConfig(cfg *Config) {
	cfg.Set("redis", map[string]any{
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
				// If there is url provided, it will be used instead of the configuration
				// and you can set all the other options in the url and include new options if they are valid
				"url":       Env("REDIS_URL", nil),
				"host":      Env("REDIS_HOST", "127.0.0.1"),
				"password":  Env("REDIS_PASSWORD", nil),
				"port":      Env("REDIS_PORT", 6379),
				"database":  Env("REDIS_DB", 10),
				"active":    Env("REDIS_DEFAULT_ACTIVE", true),
				"pool_size": Env("REDIS_POOL_SIZE", 0),
				"timeout":   Env("REDIS_TIMEOUT", ""),
			},
			"queue": map[string]any{
				"url":       Env("REDIS_QUEUE_URL", nil),
				"host":      Env("REDIS_QUEUE_HOST", "127.0.0.1"),
				"password":  Env("REDIS_QUEUE_PASSWORD", nil),
				"port":      Env("REDIS_QUEUE_PORT", 6379),
				"database":  Env("REDIS_QUEUE_DB", 9),
				"active":    Env("REDIS_QUEUE_ACTIVE", true),
				"pool_size": Env("REDIS_QUEUE_POOL_SIZE", 0),
				"timeout":   Env("REDIS_QUEUE_TIMEOUT", ""),
			},
		},
	})
}
