package load

import (
	"database/sql"
	"taskgo/internal/deps"
	"taskgo/internal/enums"
	"taskgo/pkg/ioc"
	"taskgo/pkg/utils"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase(c *ioc.Container) {
	// Register gorm database connection
	err := ioc.Singleton(c, func(c *ioc.Container) (*deps.GormDB, error) {
		cfg := deps.Config()

		driver := enums.DatabaseDriver(cfg.GetString("database.default", "postgres"))
		dsn := utils.BuildDSN(driver, cfg.GetMap("database.connections."+string(driver), nil))

		gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, err
		}

		return deps.NewGormDB(gormDB), nil
	})

	if err != nil {
		utils.PrintErr("Failed to load gorm module in the ioc container : " + err.Error())
	}

	// Register main database connection
	err = ioc.Singleton(c, func(c *ioc.Container) (*deps.GormDBSQL, error) {
		gormDB, err := ioc.Resolve[*deps.GormDB](c)
		if err != nil {
			return nil, err
		}

		db, err := gormDB.DB.DB()
		if err != nil {
			return nil, err
		}

		configureDatabaseConnectionPool(db)

		return deps.NewSqlDB(db), nil
	})

	if err != nil {
		utils.PrintErr("Failed to load sql module in the ioc container  :  " + err.Error())
	}
}

func configureDatabaseConnectionPool(db *sql.DB) {
	maxIdleConns := deps.Config().GetInt("database.max_idle_conns", 10)
	maxOpenConns := deps.Config().GetInt("database.max_open_conns", 100)
	connMaxLifetime := deps.Config().GetDuration("database.conn_max_lifetime", 30*time.Minute)

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime)
}
