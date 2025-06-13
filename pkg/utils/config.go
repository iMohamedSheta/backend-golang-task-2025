package utils

import (
	"fmt"
	"taskgo/pkg/enums"
)

/*
| Database configuration utilities
*/

// Constructs a PostgreSQL (DSN) from the database configuration postgres driver
func BuildPostgresDSN(driverConfig map[string]any) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%v sslmode=%s TimeZone=%s",
		driverConfig["host"],
		driverConfig["user"],
		driverConfig["pass"],
		driverConfig["database"],
		driverConfig["port"],
		driverConfig["sslmode"],
		driverConfig["timezone"],
	)
}

// Constructs a DSN from the database configuration
func BuildDSN(defaultDriver enums.DatabaseDriver, driverConfig map[string]any) string {
	switch defaultDriver {
	case enums.PostgresDriver:
		return BuildPostgresDSN(driverConfig)
	default:
		return BuildPostgresDSN(driverConfig) // Default to Postgres
	}
}
