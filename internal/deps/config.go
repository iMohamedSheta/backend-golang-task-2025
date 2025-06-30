package deps

import (
	"fmt"
	"taskgo/internal/config"
	"taskgo/pkg/ioc"
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
|	Application Dependency Container Calls
|--------------------------------------------------------
*/

// Config returns the app configuration instance
func Config() *config.Config {
	cfg, err := ioc.AppMake[*config.Config]()
	if err != nil {
		panic(fmt.Sprintf("Failed to resolve config: %v", err))
	}
	return cfg
}
