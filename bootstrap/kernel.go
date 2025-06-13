package bootstrap

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"taskgo/internal/api/routes"
	"taskgo/internal/config"
	"taskgo/internal/tasks"
	"taskgo/pkg/database"
	pkgEnums "taskgo/pkg/enums"
	"taskgo/pkg/logger"
	"taskgo/pkg/redis"
	"taskgo/pkg/utils"
	"taskgo/pkg/validate"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
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

	This file should be invoked from `main.go` via `bootstrap.Load()` and `bootstrap.Run()`.
*/

// Load application (env, config, DB, logger, validation)
func Load() {
	loadEnvConfig(".env")
	config.LoadAll()
	loadDatabaseConnection()
	loadLoggers()
	loadValidation()
	loadRedisConnections()
	loadTaskQueue()
}

// Run the application (serve HTTP or execute CLI command)
func Run() {
	startHttpServer()
}

// LoadWithEnv loads application with specific environment file
func LoadWithEnv(envFile string) {
	loadEnvConfig(envFile)
	config.LoadAll()
	loadDatabaseConnection()
	loadLoggers()
	loadValidation()
	loadRedisConnections()
	loadTaskQueue()
}

// Shutdown gracefully closes all loaded services
func Shutdown() {
	log.Println(pkgEnums.Yellow.Value() + "Shutting down services..." + pkgEnums.Reset.Value())

	// Close task queue client
	if err := tasks.Close(); err != nil {
		log.Printf("Error closing task queue client: %v", err)
	}

	// Close database connections
	if err := database.Close(); err != nil {
		log.Printf("Error closing database connections: %v", err)
	}

	// Close Redis connections
	if err := redis.Close(); err != nil {
		log.Printf("Error closing Redis connections: %v", err)
	}

	// Close logger (should be last)
	logger.Close()

	log.Println(pkgEnums.Green.Value() + "All services shut down properly" + pkgEnums.Reset.Value())
}

// Start the HTTP server with graceful shutdown
func startHttpServer() {

	shutdown_timeout := config.App.GetDuration("app.shutdown_timeout", 10*time.Second)

	url := config.App.GetString("app.url", "localhost")

	port := config.App.GetString("app.port", "8080")

	srv := &http.Server{
		Addr:    url + ":" + port,
		Handler: routes.RegisterRoutes(),
	}

	// Start the server in a goroutine
	go func() {
		log.Println("Starting HTTP server on http://" + url + ":" + port + " ...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error: " + err.Error())
		}
	}()

	// Listen for OS signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	log.Println(pkgEnums.Yellow.Value() + "Shutting down server..." + pkgEnums.Reset.Value())

	// timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdown_timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(pkgEnums.Red.Value() + "Forced to shutdown: " + err.Error() + pkgEnums.Reset.Value())
	}

	// Use the new shutdown method
	Shutdown()

	log.Println(pkgEnums.Green.Value() + "Server exited properly" + pkgEnums.Reset.Value())
}

// Load environment variables
func loadEnvConfig(envFile string) {
	if err := godotenv.Load(envFile); err != nil {
		log.Fatal("Error loading .env file", err)
	}
	log.Println("Kernel: Loaded .env file")
}

// Connect to the database
func loadDatabaseConnection() {
	driver := config.App.GetString("database.default", string(pkgEnums.PostgresDriver))

	databaseDriver := pkgEnums.DatabaseDriver(driver)

	driverConfigRaw, err := config.App.Get("database.connections." + driver)

	var driverConfig map[string]any
	var ok bool

	if driverConfig, ok = driverConfigRaw.(map[string]any); !ok {
		log.Fatal("Error getting database connection config:  " + err.Error())
	}

	if err != nil {
		// TODO: Add error handling
		log.Fatal("Error getting database connection config: " + err.Error())
	}

	dsn := utils.BuildDSN(databaseDriver, driverConfig)

	db, err := database.Connect(databaseDriver, dsn, loadGormFeatures())

	if err != nil {
		// TODO: Add error handling
		panic("Error connecting to database: " + err.Error())
	}

	configureDatabaseConnectionPool(db)

}

func configureDatabaseConnectionPool(db *sql.DB) {
	if db == nil {
		log.Fatal("Database connection is nil")
	}

	maxIdleConns := config.App.GetInt("database.max_idle_conns", 10)
	maxOpenConns := config.App.GetInt("database.max_open_conns", 100)
	connMaxLifetime := config.App.GetDuration("database.conn_max_lifetime", 30*time.Minute)

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime)
}

func loadGormFeatures() *gorm.Config {
	return &gorm.Config{
		// Add any GORM specific configurations here if needed
	}
}

// Load validation rules
func loadValidation() {
	validator := validate.Validator()
	for tag, rule := range registeredRules {
		if err := validator.RegisterValidation(tag, rule); err != nil {
			logger.Log().Error("Error registering validation rule: " + tag + ": " + err.Error())
		}
	}
}

// Load all loggers
func loadLoggers() {
	defaultChannel := config.App.GetString("log.default", "app_log")

	// Get all channels from log config
	channelsRaw, _ := config.App.Get("log.channels")
	channels, ok := channelsRaw.(map[string]any)
	if !ok {
		log.Fatal("Invalid log channels configuration")
	}

	var defaultLoaded bool

	// Register each channel
	for name, channelConfigRaw := range channels {
		channelConfig, ok := channelConfigRaw.(map[string]any)
		if !ok {
			continue
		}

		path, _ := channelConfig["path"].(string)
		level, _ := channelConfig["level"].(string)

		// Convert level level to zap level
		zapLevel := zap.InfoLevel
		switch level {
		case "debug":
			zapLevel = zap.DebugLevel
		case "info":
			zapLevel = zap.InfoLevel
		case "warn":
			zapLevel = zap.WarnLevel
		case "error":
			zapLevel = zap.ErrorLevel
		}

		cfg := zap.Config{
			Level:            zap.NewAtomicLevelAt(zapLevel),
			Development:      level == "debug",
			Encoding:         "json",
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
			EncoderConfig:    zap.NewProductionEncoderConfig(),
		}

		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		if name == defaultChannel {
			logger.LoadDefault(path, cfg)
			defaultLoaded = true
		} else {
			if err := logger.Register(name, path, cfg); err != nil {
				log.Printf("Failed to register logger %s: %v", name, err)
			}
		}
	}

	if !defaultLoaded {
		log.Fatal("Default logger channel not found in configuration")
	}
}

// Load Redis connections
func loadRedisConnections() {
	redisConfigRaw, err := config.App.Get("redis")
	if err != nil {
		log.Println("No Redis configuration found, skipping Redis setup")
		return
	}

	// Parse the redis config
	redisConfigMap, ok := redisConfigRaw.(map[string]any)
	if !ok {
		log.Fatal("Invalid Redis configuration format")
	}

	// Build Redis config struct
	redisConfig := redis.Config{
		Connections: make(map[string]redis.ConnectionConfig),
	}

	// Parse default connection name
	if defaultName, exists := redisConfigMap["default"]; exists {
		if name, ok := defaultName.(string); ok {
			redisConfig.Default = name
		}
	}

	// Parse options
	if optionsRaw, exists := redisConfigMap["options"]; exists {
		if optionsMap, ok := optionsRaw.(map[string]any); ok {
			if cluster, ok := optionsMap["cluster"].(string); ok {
				redisConfig.Options.Cluster = cluster
			}
			if prefix, ok := optionsMap["prefix"].(string); ok {
				redisConfig.Options.Prefix = prefix
			}
		}
	}

	// Parse all connections
	if connectionsRaw, exists := redisConfigMap["connections"]; exists {
		if connectionsMap, ok := connectionsRaw.(map[string]any); ok {
			for name, connConfigRaw := range connectionsMap {
				if connConfigMap, ok := connConfigRaw.(map[string]any); ok {
					redisConfig.Connections[name] = redis.ParseRedisConnectionConfig(connConfigMap)
				}
			}
		}
	}

	// Load Redis connections
	if err := redis.Load(redisConfig); err != nil {
		log.Printf("Error loading Redis connections: %v", err)
	}
}

// Load task queue system
func loadTaskQueue() {
	// Check if task queue should be enabled
	if !config.App.GetBool("queue.enabled", true) {
		log.Println("Task queue disabled, skipping initialization")
		return
	}

	// Initialize the task queue client should check the queue channel but not implemented
	if err := tasks.InitRedisJobsClient(); err != nil {
		log.Printf("Error initializing task queue: %v", err)
		// Decide if this should be fatal or just log the error
		if config.App.GetBool("queue.required", false) {
			log.Fatal("Task queue is required but failed to initialize")
		}
		return
	}

	log.Println("Task queue system initialized successfully")
}
