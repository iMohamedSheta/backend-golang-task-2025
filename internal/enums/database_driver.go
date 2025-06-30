package enums

// DatabaseDriver represents the type of database driver used in the application.
type DatabaseDriver string

const (
	PostgresDriver DatabaseDriver = "postgres"
	RedisDriver    DatabaseDriver = "redis"
)
