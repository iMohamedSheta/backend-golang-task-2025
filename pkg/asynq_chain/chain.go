package chainq

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"taskgo/internal/deps"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

/*
|------------------------------------------
|  Chain Orchestrator Workflow
|  (just a normal task but dispatch the next task until last one)
|------------------------------------------
|	1- Create chain of tasks as example shown
|	2- Dispatch chain (it will dispatch task (TypeChainOrchestrator"chain:orchestrator"))
|	3- This task will have all the tasks types register in the bootstrap/registers.go as handlers
|	4- and when the asynq worker (consumer/server) will call the ChainOrchestratorTask the processTask will be called
|	5- and it will check the ChainOrchestratorPayload to see the registered tasks to handle
|	6- then it will handle first task and then dispatch new task as TypeChainOrchestrator but in next step (next task)
|	7- and so on until the last task will be handled and the chain will be completed
|------------------------------------------
|	Example:
|------------------------------------------
|	tasks.NewChain().
|			Then(tasks.NewProcessPaymentTask(orderID)).
|			Then(tasks.NewCheckInventoryTask(orderID)).
|			Then(tasks.NewSendNotificationTask(orderID)).
|			OnFailure(func(err error) error {
|					// Handle chain failure
|					return nil
|			}).
|			.OnSuccess(func(result *asynq.Result) error {
|					// Handle chain success
|					return nil
|			}).
|			Dispatch()
|------------------------------------------
*/

const (
	TypeChainOrchestrator = "chain:orchestrator" // Chain orchestrator type
)

// Task interface that all tasks must implement
type Task interface {
	CreateTask() (*asynq.Task, error)
	GetTaskType() string
	GetPayload() any
}

type Logger interface {
	Info(msg string, fields ...any)
	Error(msg string, fields ...any)
	Warn(msg string, fields ...any)
}

type Chain struct {
	client     *asynq.Client
	tasks      []Task
	onSuccess  func(any) error
	onFailure  func(error) error
	maxRetries int
	timeout    time.Duration
	queue      string
	logger     Logger
}

type ChainOptions struct {
	MaxRetries   int
	Timeout      time.Duration
	DefaultQueue string
}

// NewChain creates a new task chain
func NewChain(client *asynq.Client, logger Logger, opt *ChainOptions) *Chain {
	defaultOpt := &ChainOptions{
		MaxRetries:   3,
		Timeout:      5 * time.Minute,
		DefaultQueue: "default",
	}

	if opt == nil {
		opt = defaultOpt
	} else {
		if opt.DefaultQueue == "" {
			opt.DefaultQueue = defaultOpt.DefaultQueue
		}
		if opt.Timeout == 0 {
			opt.Timeout = defaultOpt.Timeout
		}
		if opt.MaxRetries == 0 {
			opt.MaxRetries = defaultOpt.MaxRetries
		}
	}

	return &Chain{
		tasks:      make([]Task, 0),
		client:     client,
		logger:     logger,
		maxRetries: opt.MaxRetries,
		timeout:    opt.Timeout,
		queue:      opt.DefaultQueue,
	}
}

// Then adds a task to the chain
func (c *Chain) Then(task Task) *Chain {
	c.tasks = append(c.tasks, task)
	return c
}

// OnQueue sets the queue for the chain
func (c *Chain) OnQueue(queue string) *Chain {
	c.queue = queue
	return c
}

// MaxRetries sets the maximum number of retries for the entire chain
func (c *Chain) MaxRetries(retries int) *Chain {
	c.maxRetries = retries
	return c
}

// Timeout sets the timeout for each task in the chain
func (c *Chain) Timeout(timeout time.Duration) *Chain {
	c.timeout = timeout
	return c
}

// OnSuccess sets a callback for when the entire chain succeeds
func (c *Chain) OnSuccess(callback func(interface{}) error) *Chain {
	c.onSuccess = callback
	return c
}

// OnFailure sets a callback for when the chain fails
func (c *Chain) OnFailure(callback func(error) error) *Chain {
	c.onFailure = callback
	return c
}

// Dispatch dispatches the chain to the queue
func (c *Chain) Dispatch() error {
	if len(c.tasks) == 0 {
		return fmt.Errorf("no tasks in chain")
	}

	// Create the chain orchestrator payload
	chainPayload := ChainPayload{
		ChainID:     generateChainID(),
		Tasks:       c.serializeTasks(),
		CurrentStep: 0,
		MaxRetries:  c.maxRetries,
		Timeout:     c.timeout,
		Queue:       c.queue,
	}

	payload, err := json.Marshal(chainPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal chain payload: %w", err)
	}

	task := asynq.NewTask(TypeChainOrchestrator, payload,
		asynq.MaxRetry(c.maxRetries),
		asynq.Timeout(c.timeout),
		asynq.Queue(c.queue),
	)
	err = dispatchAsynqTask(c.client, c.logger, task) // Dispatch the orchestrator task
	log := deps.Log().Channel("queue_log")
	if err != nil {
		log.Error("Failed to dispatch chain", zap.Error(err))
		return err
	}

	log.Info("Chain dispatched successfully",
		zap.String("chain_id", chainPayload.ChainID),
		zap.Int("task_count", len(c.tasks)),
	)

	return nil
}

type ChainPayload struct {
	ChainID     string                 `json:"chain_id"`
	Tasks       []SerializedTask       `json:"tasks"`
	CurrentStep int                    `json:"current_step"`
	MaxRetries  int                    `json:"max_retries"`
	Timeout     time.Duration          `json:"timeout"`
	Queue       string                 `json:"queue"`
	Context     map[string]interface{} `json:"context,omitempty"` // Shared data between tasks
}

type SerializedTask struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// serializeTasks converts Task to SerializedTask so it can be stored inside the ChainPayload
func (c *Chain) serializeTasks() []SerializedTask {
	serialized := make([]SerializedTask, len(c.tasks))
	for i, task := range c.tasks {
		serialized[i] = SerializedTask{
			Type:    task.GetTaskType(),
			Payload: task.GetPayload(),
		}
	}
	return serialized
}

// generateChainID creates a unique ID for the chain
func generateChainID() string {
	prefix := "CHAIN"
	return strings.ToUpper(fmt.Sprintf("%s_%s", prefix, uuid.New().String()))
}

// ChainOrchestrator handles the execution of chained tasks
type ChainOrchestrator struct {
	handlers map[string]asynq.Handler // Map of task type to handler
	client   *asynq.Client
	logger   Logger
}

func NewChainOrchestrator(client *asynq.Client, logger Logger) *ChainOrchestrator {
	return &ChainOrchestrator{
		handlers: make(map[string]asynq.Handler),
		client:   client,
		logger:   logger,
	}
}

// RegisterHandler registers a task handler
func (co *ChainOrchestrator) RegisterHandler(taskType string, handler asynq.Handler) {
	co.handlers[taskType] = handler
}

// ProcessTask is what happen when the orchestrator is called to process from the asynq worker
func (co *ChainOrchestrator) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload ChainPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal chain payload: %w", err)
	}

	log := deps.Log().Channel("queue_log")
	log.Info("Processing chain step",
		zap.String("chain_id", payload.ChainID),
		zap.Int("step", payload.CurrentStep+1),
		zap.Int("total", len(payload.Tasks)),
	)

	// Check if we've completed all tasks
	if payload.CurrentStep >= len(payload.Tasks) {
		log.Info("Chain completed successfully",
			zap.String("chain_id", payload.ChainID),
		)
		return nil
	}

	// Get current task
	currentTask := payload.Tasks[payload.CurrentStep]

	// Find handler for current task
	handler, exists := co.handlers[currentTask.Type]
	if !exists {
		return fmt.Errorf("no handler registered for task type: %s", currentTask.Type)
	}

	// Create asynq task for the current step
	taskPayload, err := json.Marshal(currentTask.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(currentTask.Type, taskPayload)

	// Execute the current task
	if err := handler.ProcessTask(ctx, task); err != nil {
		log.Error("Chain task failed",
			zap.String("chain_id", payload.ChainID),
			zap.String("task_type", currentTask.Type),
			zap.Int("step", payload.CurrentStep+1),
			zap.Error(err),
		)
		return fmt.Errorf("chain failed at step %d (%s): %w",
			payload.CurrentStep+1, currentTask.Type, err)
	}

	// Task succeeded, move to next step
	payload.CurrentStep++

	// If there are more tasks, dispatch the next step
	if payload.CurrentStep < len(payload.Tasks) {
		return co.dispatchNextStep(payload)
	}

	// Chain completed
	log.Info("Chain completed successfully",
		zap.String("chain_id", payload.ChainID),
	)
	return nil
}

// dispatchNextStep dispatches the next step in the chain
func (co *ChainOrchestrator) dispatchNextStep(payload ChainPayload) error {
	nextPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal next step payload: %w", err)
	}

	task := asynq.NewTask(TypeChainOrchestrator, nextPayload,
		asynq.MaxRetry(payload.MaxRetries),
		asynq.Timeout(payload.Timeout),
		asynq.Queue(payload.Queue),
		asynq.ProcessIn(1*time.Second), // delay between steps
	)

	return dispatchAsynqTask(co.client, co.logger, task)
}

func dispatchAsynqTask(client *asynq.Client, log Logger, task *asynq.Task, opts ...asynq.Option) error {
	if client == nil {
		log.Error("asynq client not initialized")
		return fmt.Errorf("asynq client not initialized")
	}

	info, err := client.Enqueue(task, opts...)
	if err != nil {
		log.Error(fmt.Sprintf("failed to enqueue task: %v", err))
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info(fmt.Sprintf("Task enqueued: ID=%s, Queue=%s, Type=%s", info.ID, info.Queue, info.Type))
	return nil
}
