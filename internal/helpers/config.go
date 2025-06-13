package helpers

import (
	"fmt"
	"taskgo/internal/config"
	"taskgo/pkg/logger"
)

// GetAppSecret returns the app secret from the configuration
func GetAppSecret() ([]byte, error) {
	secret := config.App.GetString("app.secret", "")
	if secret == "" {
		logger.Log().Error("Missing or invalid app secret in configuration")
		return nil, fmt.Errorf("missing or invalid app secret in configuration")
	}

	return []byte(secret), nil
}
