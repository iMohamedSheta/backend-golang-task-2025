package tests

import (
	"taskgo/bootstrap"
	"taskgo/internal/database/migration"
	"taskgo/pkg/database"
)

// TestSetup initializes the application for testing
func TestSetup() {
	// Load the application with test environment
	bootstrap.LoadWithEnv("../.env.testing")

	// Run database migrations for testing
	migrateDBForTesting()
}

// TestTeardown cleans up test resources
func TestTeardown() {
	// Close database connections
	database.Close()
}

// migrateDBForTesting runs database migrations for testing
func migrateDBForTesting() {
	// Run database migrations
	migration.RunMigrations(database.GetDB())

	// Run database seeder
	// seeders.SeedDatabase()
}

func truncateTables() {
	migration.RollbackMigrations(database.GetDB())

	// Run database migrations
	migration.RunMigrations(database.GetDB())
}
