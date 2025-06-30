package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"taskgo/internal/adapters"
	"taskgo/internal/deps"
	chainq "taskgo/pkg/asynq_chain"
	"time"

	"github.com/hibiken/asynq"
)

// Dispatch a new task to the queue
func Dispatch(task chainq.Task, opts ...asynq.Option) error {
	client := deps.Queue().Client
	if client == nil {
		return fmt.Errorf("asynq client not initialized")
	}

	asynqTask, err := task.CreateTask()
	if err != nil {
		return err
	}

	info, err := client.Enqueue(asynqTask, opts...)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf("Task enqueued: ID=%s, Queue=%s, Type=%s",
		info.ID, info.Queue, info.Type)
	return nil
}

// // Dispatch a new task to the queue
// func DispatchAsynqTask(task *asynq.Task, opts ...asynq.Option) error {
// 	client := deps.Queue().Client
// 	if client == nil {
// 		return fmt.Errorf("asynq client not initialized")
// 	}

// 	info, err := client.Enqueue(task, opts...)
// 	if err != nil {
// 		return fmt.Errorf("failed to enqueue task: %w", err)
// 	}

// 	log.Printf("Task enqueued: ID=%s, Queue=%s, Type=%s",
// 		info.ID, info.Queue, info.Type)
// 	return nil
// }

/*
	Helpers to create new asynq Task or TaskHandler used inside the tasks
*/

// Create new asynq task to from task to use it inside package
func CreateAsynqTask(t chainq.Task, opts ...asynq.Option) (*asynq.Task, error) {
	payload, err := json.Marshal(t.GetPayload())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return asynq.NewTask(t.GetTaskType(), payload, opts...), nil
}

// Process task payload helper to send payload to actual handler
func processTaskPayload[T any](ctx context.Context, task *asynq.Task, handler func(context.Context, *T) error) error {
	var payload T
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	return handler(ctx, &payload)
}

// Chain - helper to create new chainq.Chain which can be used to create new chain of tasks
func Chain() *chainq.Chain {
	return chainq.NewChain(
		deps.Queue().Client,
		adapters.NewLoggerAdapter(deps.Log().Log()),
		&chainq.ChainOptions{
			MaxRetries:   3,
			Timeout:      5 * time.Second,
			DefaultQueue: "default",
		},
	)
}
