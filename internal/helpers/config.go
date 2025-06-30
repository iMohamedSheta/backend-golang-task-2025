package helpers

import (
	"taskgo/internal/deps"
	"taskgo/pkg/errors"
)

// GetAppSecret returns the app secret from the configuration
func GetAppSecret() ([]byte, error) {
	secret := deps.Config().GetString("app.secret", "")
	if secret == "" {
		deps.Log().Log().Error("Missing or invalid app secret in configuration")
		return nil, errors.NewServerError("internal server error", "internal server error: Missing or invalid app secret in configuration", nil)
	}

	return []byte(secret), nil
}
