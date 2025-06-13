package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"taskgo/pkg/logger"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

// ProcessPaymentTask implement ChainableTask interface also it's used as payload for task
type ProcessPaymentTask struct {
	OrderID uint `json:"order_id"`
}

func NewProcessPaymentChainTask(orderID uint) *ProcessPaymentTask {
	return &ProcessPaymentTask{OrderID: orderID}
}

func (t *ProcessPaymentTask) GetTaskType() string {
	return TypeProcessPayment
}

func (t *ProcessPaymentTask) GetPayload() interface{} {
	return *t // Return itself as payload
}

func (t *ProcessPaymentTask) CreateTask() (*asynq.Task, error) {
	payload, err := json.Marshal(t.GetPayload())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(t.GetTaskType(), payload), nil
}

//------------------------------------------------------------------------------------------------------------------

// Create a new ProcessPaymentTask
func NewProcessPaymentTask(orderID uint) (*asynq.Task, error) {
	payload, err := json.Marshal(ProcessPaymentTask{
		OrderID: orderID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeProcessPayment, payload), nil
}

// Dispatch a new ProcessPaymentTask to queue
func DispatchProcessPaymentTask(orderID uint) error {
	processPaymentTask, err := NewProcessPaymentTask(orderID)
	if err != nil {
		logger.Channel("queue_log").Error("Failed to create new process payment task", zap.Error(err))
		return err
	}

	err = Dispatch(processPaymentTask, asynq.Queue(QueuePayments))
	if err != nil {
		logger.Channel("queue_log").Error("Failed to dispatch task", zap.Error(err))
		return err
	}
	logger.Channel("queue_log").Info("Dispatched task", zap.String("type", processPaymentTask.Type()), zap.String("queue", QueuePayments))
	return nil
}

/*
|------------------------------------------
|  Task handler: ProcessPaymentHandler
|------------------------------------------
*/
type ProcessPaymentHandler struct {
}

// Return a new payment task Handler
func NewProcessPaymentHandler() *ProcessPaymentHandler {
	return &ProcessPaymentHandler{}
}

// Handler method for the payment task implement Handler interface
func (p *ProcessPaymentHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload ProcessPaymentTask
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	logger.Channel("queue_log").Info(fmt.Sprintf("Processing payment for Order: %d", payload.OrderID))

	// Your payment processing logic here
	if err := p.handle(payload); err != nil {
		return fmt.Errorf("payment processing failed: %w", err)
	}

	logger.Channel("queue_log").Info(fmt.Sprintf("Payment processing completed for Order:  %d", payload.OrderID))
	return nil
}

/*
|-------------------------------------------------
|  Actual task handling code goes here:
|-------------------------------------------------
*/
func (p *ProcessPaymentHandler) handle(payload ProcessPaymentTask) error {
	// Here is the actual payment processing logic
	// ...
	logger.Channel("queue_log").Info(fmt.Sprintf("Processed payment for Order: %d", payload.OrderID))

	return nil
}
