package bootstrap

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"taskgo/internal/adapters"
	"taskgo/internal/api/routes"
	"taskgo/internal/deps"
	chainq "taskgo/pkg/asynq_chain"
	"taskgo/pkg/ioc"
	"taskgo/pkg/utils"
	"time"

	"github.com/hibiken/asynq"
)

/*
	Package bootstrap kernel initializes and runs the core application lifecycle.

	This file is the main entry point for bootstrapping the application. It:
	- Loads essential configurations (environment variables, app settings, DB, logger, validation rules).
	- Starts the HTTP server with graceful shutdown capabilities on system signals (e.g., SIGTERM).

	Modules involved:
	- config: loads and provides access to application configuration.
	- logger: sets up the global logging system.
	- database: initializes the database connection.
	- redis: provides access to Redis connections.
	- validate: registers custom validation rules.
	- routes: provides the HTTP router.

	This file should be invoked from `main.go` via `bootstrap.NewAppBuilder()` and `bootstrap.Run()`.
*/

// Run the application (serve HTTP)
func Run() {
	startHttpServer()
}

// Shutdown gracefully closes all loaded services
func Shutdown() {
	utils.PrintWarning("Shutting down services...")

	ioc.AppContainer().ShutdownAll()

	utils.PrintSuccess("All services shut down properly")
}

// Start the HTTP server with graceful shutdown
func startHttpServer() {
	cfg := deps.Config()
	shutdown_timeout := cfg.GetDuration("app.shutdown_timeout", 10*time.Second)
	bindAddress := cfg.GetString("app.bind_address", "0.0.0.0")
	bindPort := cfg.GetString("app.bind_port", "8080")
	url := cfg.GetString("app.url", "http://localhost")
	port := cfg.GetString("app.port", "8080")

	srv := &http.Server{
		Addr:    bindAddress + ":" + bindPort,
		Handler: routes.RegisterRoutes(),
	}

	// Start the server in a goroutine
	go func() {
		utils.PrintInfo("Starting HTTP server on " + url + ":" + port + " ...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.PrintErr("Server error:", err)
			os.Exit(1)
		}
	}()

	// Listen for OS signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	utils.PrintWarning("Shutting down HTTP server...")

	// timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdown_timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		utils.PrintErr("Server forced to shutdown:", err)
		os.Exit(1)
	}

	Shutdown()

	utils.PrintSuccess("Server exited properly")
}

// GetRegisteredTaskHandlers returns all task handlers including the configured orchestrator
func GetRegisteredTaskHandlers() map[string]asynq.Handler {
	once.Do(func() {
		// Create the orchestrator first
		orchestrator := chainq.NewChainOrchestrator(
			deps.Queue().Client,
			adapters.NewLoggerAdapter(deps.Log().Log()),
		)

		// Get all individual task handlers
		individualHandlers := registerTaskHandlers()

		// Auto-register all individual handlers with the orchestrator
		for taskType, handler := range individualHandlers {
			orchestrator.RegisterHandler(taskType, handler)
		}

		// Initialize the global registeredTasks with individual handlers
		registeredTasks = make(map[string]asynq.Handler)
		for taskType, handler := range individualHandlers {
			registeredTasks[taskType] = handler
		}

		// Add the configured orchestrator to the map
		registeredTasks[chainq.TypeChainOrchestrator] = orchestrator
	})

	return registeredTasks
}
