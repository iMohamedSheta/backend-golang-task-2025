package tests

import (
	"fmt"
	"os"
	"taskgo/bootstrap"
	"taskgo/internal/database/migration"
	"taskgo/internal/deps"
	"testing"
)

// TestMain runs before all tests
func TestMain(m *testing.M) {
	TestingSetup()
	code := m.Run()
	TestingTeardown()
	os.Exit(code)
}

// TestSetup initializes the application for testing
func TestingSetup() {
	// Load the application with test environment
	bootstrap.NewAppBuilder("../.env.testing").
		LoadConfig().Then(updateLoggerPath).
		LoadLogger().
		LoadDatabase().
		LoadValidator().
		LoadRedisCache().
		LoadRedisQueue().
		Boot()

	// Run database migrations for testing
	migrateDBForTesting()
}

// TestTeardown cleans up test resources
func TestingTeardown() {
	// Close database connections
	db, err := deps.Gorm().DB.DB()
	if err != nil {
		panic(err)
	}
	db.Close()
}

// migrateDBForTesting runs database migrations for testing
func migrateDBForTesting() {
	// Run database migrations
	migration.RunMigrations(deps.Gorm().DB)

	// Run database seeder
	// seeders.SeedDatabase()
}

func truncateTables() {
	db := deps.Gorm().DB
	migration.RollbackMigrations(db)

	// Run database migrations
	migration.RunMigrations(db)
}

func updateLoggerPath() {
	cfg := deps.Config()

	cfgRaw, err := cfg.Get("log.channels")
	if err != nil {
		return
	}

	channels, ok := cfgRaw.(map[string]any)
	if !ok {
		panic("log.channels is not a map")
	}

	for channelName, channelCfgRaw := range channels {
		channelCfg, ok := channelCfgRaw.(map[string]any)
		if !ok {
			continue // skip invalid config
		}

		originalPath, ok := channelCfg["path"].(string)
		if !ok {
			continue // skip if no path
		}

		// Update the path to point to test dir
		newPath := "../" + originalPath
		cfg.Set(fmt.Sprintf("log.channels.%s.path", channelName), newPath)
	}
}
