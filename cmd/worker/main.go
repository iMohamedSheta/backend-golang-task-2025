package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"taskgo/bootstrap"
	"taskgo/internal/deps"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	bootstrap.NewAppBuilder(".env").
		LoadConfig().
		LoadLogger().
		LoadDatabase().
		LoadValidator().
		LoadRedisCache().
		LoadRedisQueue().
		LoadNotify().
		Boot()

	// Initialize and run task worker
	runTaskWorker()

	bootstrap.Shutdown()
}

func runTaskWorker() {
	cfg := deps.Config()

	if !cfg.GetBool("queue.enabled", true) {
		log.Fatal("Queue is disabled, cannot start consumer")
	}

	redisCfg := cfg.GetMap("redis.connections.queue", nil)
	if redisCfg == nil {
		log.Fatal("redis connection config is missing for queue consumer")
	}
	// Get Redis options for server
	redisOpt, err := convertRedisConfigToOptions(redisCfg)
	if err != nil {
		log.Fatal("Failed to get Redis options:", err)
	}

	concurrency := cfg.GetInt("queue.consumer.concurrency", 10)
	queuesRaw, _ := cfg.Get("queue.consumer.queues")

	queues := convertQueuePriorities(queuesRaw)

	// Server configuration
	serverConfig := asynq.Config{
		Concurrency:    concurrency,
		Queues:         queues,
		RetryDelayFunc: asynq.DefaultRetryDelayFunc,

		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			// Check if we should log failed tasks
			if cfg.GetBool("queue.consumer.logging.log_failed_tasks", true) {
				// Get retry information from context
				retryCount, _ := asynq.GetRetryCount(ctx)
				maxRetry, _ := asynq.GetMaxRetry(ctx)
				taskID, _ := asynq.GetTaskID(ctx)

				deps.Log().Channel("queue_log").Error("Task processing failed",
					zap.String("task_type", task.Type()),
					zap.String("task_id", taskID),
					zap.Int("retry_count", retryCount),
					zap.Int("max_retry", maxRetry),
					zap.Error(err),
				)
			}
		}),
	}

	server := asynq.NewServer(asynq.RedisClientOpt{
		Addr:      redisOpt.Addr,
		Username:  redisOpt.Username,
		Password:  redisOpt.Password,
		DB:        redisOpt.DB,
		PoolSize:  redisOpt.PoolSize,
		TLSConfig: redisOpt.TLSConfig,
	}, serverConfig)

	// Register task handlers
	mux := asynq.NewServeMux()

	// Get registered task handlers
	handlers := bootstrap.GetRegisteredTaskHandlers()
	log.Printf("Registered %d task handlers", len(handlers))
	for taskType, handler := range handlers {
		// Logging decorator
		log.Printf("Registering task handler: %s", taskType)
		wrappedHandler := wrapHandlerWithLogging(handler, taskType)
		mux.Handle(taskType, wrappedHandler)
		log.Printf("Registered task handler: %s", taskType)
	}

	log.Printf("Starting task worker with %d concurrency...", concurrency)
	log.Printf("Queue priorities: %+v", queues)

	// Run server (blocking) - asynq handles graceful shutdown internally
	if err := server.Run(mux); err != nil {
		log.Fatal("Task server error:", err)
	}
}

// Wrap handler with logging based on configuration
func wrapHandlerWithLogging(handler asynq.Handler, taskType string) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
		start := time.Now()

		// Execute the actual handler
		err := handler.ProcessTask(ctx, task)

		duration := time.Since(start)

		if err != nil {
			// Error logging is handled by ErrorHandler
			return err
		}

		// Log successful tasks if enabled
		if deps.Config().GetBool("queue.consumer.logging.log_success", false) {
			taskID, _ := asynq.GetTaskID(ctx)

			deps.Log().Channel("queue_log").Info("Task completed successfully",
				zap.String("task_type", taskType),
				zap.String("task_id", taskID),
				zap.Duration("duration", duration),
			)
		}

		return nil
	})
}

func convertQueuePriorities(raw any) map[string]int {
	result := make(map[string]int)

	queuesMap, ok := raw.(map[string]any)
	if !ok {
		return result
	}

	for k, v := range queuesMap {
		switch val := v.(type) {
		case int:
			result[k] = val
		case int64:
			result[k] = int(val)
		case float64:
			result[k] = int(val)
		case string:
			if parsed, err := strconv.Atoi(val); err == nil {
				result[k] = parsed
			}
		}
	}

	return result
}

// convertRedisConfigToOptions converts the app redis config to a redis options struct
func convertRedisConfigToOptions(cfg map[string]any) (*redis.Options, error) {
	if isActive, ok := cfg["active"].(bool); ok && !isActive {
		return nil, errors.New("redis connection is not active")
	}

	// Prefer URL if provided
	if rawURL, ok := cfg["url"].(string); ok && rawURL != "" {
		opt, err := redis.ParseURL(rawURL)
		if err != nil {
			return nil, fmt.Errorf("invalid redis url: %w", err)
		}

		return opt, nil
	}

	// Fallback to manual config
	host := cfg["host"].(string)
	port := cfg["port"].(int)
	password, _ := cfg["password"].(string)
	db, _ := cfg["database"].(int)
	poolSize, _ := cfg["pool_size"].(int)

	opt := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	}

	if poolSize > 0 {
		opt.PoolSize = poolSize
	}

	if timeoutStr, ok := cfg["timeout"].(string); ok && timeoutStr != "" {
		timeout, err := time.ParseDuration(timeoutStr)
		if err == nil {
			opt.DialTimeout = timeout
		}
	}

	return opt, nil
}
