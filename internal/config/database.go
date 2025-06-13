package config

import (
	pkgEnums "taskgo/pkg/enums"
)

func init() {
	Register(databaseConfig)
}

func databaseConfig() {
	App.Set("database", map[string]any{

		// This is the default database connection should be valid connection to use.
		"default": Env("DB_CONNECTION", string(pkgEnums.PostgresDriver)),

		"connections": map[string]any{
			// Postgres connection
			string(pkgEnums.PostgresDriver): map[string]any{
				"host":     Env("DB_HOST", "localhost"),
				"port":     Env("DB_PORT", 5432),
				"user":     Env("DB_USERNAME", "root"),
				"pass":     Env("DB_PASSWORD", ""),
				"database": Env("DB_DATABASE", "go"),
				"driver":   "pgsql",
				"sslmode":  "disable",
				"timezone": "Africa/Cairo",
			},
		},
		// connection pool settings
		"max_idle_conns":    Env("DB_MAX_IDLE_CONNS", 10),
		"max_open_conns":    Env("DB_MAX_OPEN_CONNS", 100),
		"conn_max_lifetime": Env("DB_CONN_MAX_LIFETIME", "30m"),
	})
}
