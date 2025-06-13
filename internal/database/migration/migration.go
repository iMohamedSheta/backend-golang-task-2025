package migration

import (
	"log"

	"taskgo/internal/database/models"

	"gorm.io/gorm"
)

// RunMigrations runs all database migrations
func RunMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Inventory{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
		&models.Notification{},
		&models.AuditLog{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// RollbackMigrations drops all tables  TODO: should be prevented in production
func RollbackMigrations(db *gorm.DB) error {
	log.Println("Rolling back database migrations...")

	// Drop all tables
	err := db.Migrator().DropTable(
		&models.AuditLog{},
		&models.Notification{},
		&models.Payment{},
		&models.OrderItem{},
		&models.Order{},
		&models.Inventory{},
		&models.Product{},
		&models.User{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations rolled back successfully")
	return nil
}
