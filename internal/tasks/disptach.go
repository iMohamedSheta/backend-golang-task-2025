package tasks

import (
	"fmt"
	"log"
	"taskgo/pkg/redis"

	"github.com/hibiken/asynq"
)

// Global asynq client instance
var client *asynq.Client

// Get Redis jobs client options
func GetRedisJobsClientOptions() (asynq.RedisClientOpt, error) {
	// Check if jobs connection is active
	if !redis.IsConnectionActive("jobs") {
		return asynq.RedisClientOpt{}, fmt.Errorf("jobs redis connection not active")
	}

	// Get the jobs Redis client
	jobsClient, err := redis.Jobs()
	if err != nil {
		return asynq.RedisClientOpt{}, fmt.Errorf("failed to get jobs redis client: %w", err)
	}

	// Extract connection options from the existing client
	opts := jobsClient.Options()

	return asynq.RedisClientOpt{
		Addr:      opts.Addr,
		Password:  opts.Password,
		DB:        opts.DB,
		PoolSize:  opts.PoolSize,
		TLSConfig: opts.TLSConfig,
	}, nil
}

// Initialize the global client using jobs Redis connection
func InitRedisJobsClient() error {
	redisOpt, err := GetRedisJobsClientOptions()
	if err != nil {
		return fmt.Errorf("failed to get redis options: %w", err)
	}

	client = asynq.NewClient(redisOpt)
	log.Println("Asynq client initialized with jobs Redis connection")
	return nil
}

// Dispatch function with various options
func Dispatch(task *asynq.Task, opts ...asynq.Option) error {
	if client == nil {
		return fmt.Errorf("asynq client not initialized")
	}

	info, err := client.Enqueue(task, opts...)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf("Task enqueued: ID=%s, Queue=%s, Type=%s",
		info.ID, info.Queue, info.Type)
	return nil
}

// Close closes the asynq client
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
