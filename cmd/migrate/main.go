package main

import (
	"log"
	"os"
	"taskgo/bootstrap"
	"taskgo/internal/database/migration"
	"taskgo/internal/database/seeders"
	"taskgo/internal/deps"
	"taskgo/pkg/enums"
)

func main() {
	// Load application parts for migrations to work
	bootstrap.NewAppBuilder(".env").
		LoadConfig().
		LoadLogger().
		LoadDatabase().
		Boot()

	var err error
	gormDB := deps.Gorm().DB
	if len(os.Args) > 1 && os.Args[1] == "rollback" {
		if deps.Config().GetString("app.env", "prod") == "prod" {
			log.Fatal(enums.Red.Value() + "Rollback is not allowed in production environment" + enums.Reset.Value())
			return
		}
		// Rollback database migrations
		err = migration.RollbackMigrations(gormDB)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// Run database migrations
	err = migration.RunMigrations(gormDB)

	if err != nil {
		log.Fatal(err)
		return
	}

	// Run database seeder
	seeders.SeedDatabase()
}
