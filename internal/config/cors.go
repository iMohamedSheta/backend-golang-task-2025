package config

func init() {
	Register(corsConfig)
}

func corsConfig(cfg *Config) {
	cfg.Set("cors", map[string]any{
		"origin": []string{
			"*",
		},
		"methods": []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		"allowed_headers": []string{
			"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token",
			"Authorization", "accept", "origin", "Cache-Control", "X-Requested-With",
		},
		"exposed_headers": []string{},
		"credentials":     true,    // Allow cookies, HTTP auth, etc.
		"max_age":         "86400", // Preflight cache in seconds
	})
}
