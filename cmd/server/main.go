//	@title			TaskGo API
//	@version		1.0
//	@description	Order Processing System API

//	@contact.name	iMohamedSheta
//	@contact.url	https://github.com/iMohamedSheta
//	@contact.email	mohamed15.sheta15@gmail.com

//	@license.name	MIT License
//	@license.url	https://opensource.org/licenses/MIT

//	@host		localhost:8080
//	@BasePath	/api/v1

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT token.
package main

import (
	"taskgo/bootstrap"
	_ "taskgo/docs"
)

func main() {
	// Load application parts (configurations, DB connection, logger, router, validator, etc..)
	bootstrap.Load()

	// Start the application (HTTP server)
	bootstrap.Run()
}
