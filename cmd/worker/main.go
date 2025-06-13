package main

import (
	"context"
	"log"
	"strconv"
	"taskgo/bootstrap"
	"taskgo/internal/config"
	"taskgo/internal/tasks"
	"taskgo/pkg/logger"
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func main() {
	bootstrap.Load()

	// Initialize and run task worker
	runTaskWorker()

	bootstrap.Shutdown()
}

func runTaskWorker() {
	// Check if queue consumer is enabled
	if !config.App.GetBool("queue.enabled", true) {
		log.Fatal("Queue is disabled, cannot start consumer")
	}

	// Get Redis options for server
	redisOpt, err := tasks.GetRedisJobsClientOptions()
	if err != nil {
		log.Fatal("Failed to get Redis options:", err)
	}

	concurrency := config.App.GetInt("queue.consumer.concurrency", 10)
	queuesRaw, _ := config.App.Get("queue.consumer.queues")

	queues := convertQueuePriorities(queuesRaw)

	// Server configuration
	serverConfig := asynq.Config{
		Concurrency:    concurrency,
		Queues:         queues,
		RetryDelayFunc: asynq.DefaultRetryDelayFunc,

		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			// Check if we should log failed tasks
			if config.App.GetBool("queue.consumer.logging.log_failed_tasks", true) {
				// Get retry information from context
				retryCount, _ := asynq.GetRetryCount(ctx)
				maxRetry, _ := asynq.GetMaxRetry(ctx)
				taskID, _ := asynq.GetTaskID(ctx)

				logger.Channel("queue_log").Error("Task processing failed",
					zap.String("task_type", task.Type()),
					zap.String("task_id", taskID),
					zap.Int("retry_count", retryCount),
					zap.Int("max_retry", maxRetry),
					zap.Error(err),
				)
			}
		}),
	}

	server := asynq.NewServer(redisOpt, serverConfig)

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
		if config.App.GetBool("queue.consumer.logging.log_success", false) {
			taskID, _ := asynq.GetTaskID(ctx)

			logger.Channel("queue_log").Info("Task completed successfully",
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
