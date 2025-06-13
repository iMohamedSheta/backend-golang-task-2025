package main

import (
	"log"
	"os"
	"taskgo/bootstrap"
	"taskgo/internal/config"
	"taskgo/internal/database/migration"
	"taskgo/internal/database/seeders"
	"taskgo/pkg/database"
	"taskgo/pkg/enums"
)

func main() {
	// Load application parts (configurations, DB connection, logger, router, validator, etc..)
	// can be updated to just connect to database and run migrations
	bootstrap.Load()

	var err error
	if len(os.Args) > 1 && os.Args[1] == "rollback" {
		if config.App.GetString("app.env", "prod") == "prod" {
			log.Fatal(enums.Red.Value() + "Rollback is not allowed in production environment" + enums.Reset.Value())
			return
		}
		// Rollback database migrations
		err = migration.RollbackMigrations(database.GetDB())
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// Run database migrations
	err = migration.RunMigrations(database.GetDB())

	if err != nil {
		log.Fatal(err)
		return
	}

	// Run database seeder
	seeders.SeedDatabase()
}
