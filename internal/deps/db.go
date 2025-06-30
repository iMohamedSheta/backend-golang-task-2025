package deps

import (
	"database/sql"
	"fmt"
	"taskgo/pkg/ioc"

	"gorm.io/gorm"
)

/*
	Package deps provides convenient aliases to resolve core project dependencies
	(e.g., database, config, logger, etc.) from the IOC container.

	This package simplifies service access by wrapping the IOC resolution logic in
	easy-to-use functions, improving readability and reducing repetitive boilerplate.

	Instead of writing:
		ioc.Make[*gorm.DB](ioc.App())

	You can simply call:
		deps.DB()

	Note on Aliases:
	The IOC container resolves services by their Go type. This means you cannot
	have multiple services of the same type (e.g., multiple *sql.DB instances) without
	conflict. To solve this, you can define alias structs and register them as unique types.

	Example:
		type MainDB struct {
			DB *sql.DB
		}

	register it in the IOC App container:
		ioc.Register[*MainDB](etc..

	Then resolve via:
		ioc.Make[*MainDB](ioc.App()).DB

	This pattern enables safe multi-instance registration for services with the same type.
*/

/*
|--------------------------------------------------------
|	Application Dependency Container Aliases
|--------------------------------------------------------
*/

type GormDB struct {
	DB *gorm.DB
}

type GormDBSQL struct {
	DB *sql.DB
}

/*
|--------------------------------------------------------
|	Application Dependency Container Calls
|--------------------------------------------------------
*/

// Gorm returns a GormService which wraps the *gorm.DB dependency.
func Gorm() *GormDB {
	db, err := ioc.AppMake[*GormDB]()
	if err != nil {
		panic(fmt.Sprintf("Failed to make GormDB: %s", err.Error()))
	}
	return db
}

func NewGormDB(gorm *gorm.DB) *GormDB {
	return &GormDB{
		DB: gorm,
	}
}

// Sql returns a SqlService which wraps the *sql.DB dependency.
func GormSql() *GormDBSQL {
	db, err := ioc.AppMake[*GormDBSQL]()
	if err != nil {
		panic(fmt.Sprintf("Failed to make SqlDB:  %s", err.Error()))
	}
	return db
}

// GormSqlDB which is the unerlying *sql.DB for GormDB should be shutdown so we implement shutdown interface for ioc to shutdown it in the end.
func (gSqlDB *GormDBSQL) Shutdown() error {
	if gSqlDB.DB != nil {
		return gSqlDB.DB.Close()
	}
	return nil
}

func NewSqlDB(db *sql.DB) *GormDBSQL {
	return &GormDBSQL{
		DB: db,
	}
}
