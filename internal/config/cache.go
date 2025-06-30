package config

func init() {
	Register(cacheConfig)
}

func cacheConfig(cfg *Config) {
	cfg.Set("cache", map[string]any{
		"default": "redis",
		"stores": map[string]any{
			"redis": map[string]any{
				"driver": "redis",
				"host":   "127.0.0.1",
				"port":   "6379",
			},
		},
	})
}
