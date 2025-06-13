package database

import (
	"database/sql"
	"fmt"
	"log"
	"runtime"
	"sync"
	"taskgo/pkg/enums"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

// Init the DB connection
func Connect(driver enums.DatabaseDriver, dsn string, config *gorm.Config) (*sql.DB, error) {
	var err error
	var sqlDB *sql.DB

	once.Do(func() {

		switch driver {
		case enums.PostgresDriver:
			db, err = gorm.Open(postgres.Open(dsn), config)

		default:
			err = fmt.Errorf("unsupported driver: %s", driver)
			return
		}

		if err != nil {
			return
		}

		sqlDB, err = db.DB()
		if err != nil {
			return
		}
	})

	if err != nil {
		return nil, err
	}

	return sqlDB, nil
}

// GetDB returns the singleton DB instance
func GetDB() *gorm.DB {
	if db == nil {
		// Capture the caller of this function
		pc, file, line, ok := runtime.Caller(1)
		if ok {
			fn := runtime.FuncForPC(pc)
			log.Printf("GetDB called from %s:%d (%s)", file, line, fn.Name())
		}

		log.Fatal("Database not initialized")
	}
	return db
}

// Closes underlying sql.DB connection pool gracefully
func Close() error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
