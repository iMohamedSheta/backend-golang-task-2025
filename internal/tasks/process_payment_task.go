package tasks

import (
	"context"
	"fmt"
	"taskgo/internal/deps"

	"github.com/hibiken/asynq"
)

// ProcessPaymentTask implement Task interface also it's used as payload for task
type ProcessPaymentTask struct {
	OrderID uint `json:"order_id"`
}

func NewProcessPaymentTask(orderID uint) *ProcessPaymentTask {
	return &ProcessPaymentTask{OrderID: orderID}
}

func (t *ProcessPaymentTask) GetTaskType() string {
	return TypeProcessPayment
}

func (t *ProcessPaymentTask) GetPayload() interface{} {
	return *t // Return itself as payload
}

func (t *ProcessPaymentTask) CreateTask() (*asynq.Task, error) {
	return CreateAsynqTask(t, asynq.Queue(QueuePayments), asynq.MaxRetry(3))
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
func (h *ProcessPaymentHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	return processTaskPayload(ctx, t, h.handle)
}

/*
|-------------------------------------------------
|  Actual task handling code goes here:
|-------------------------------------------------
*/
func (p *ProcessPaymentHandler) handle(ctx context.Context, payload *ProcessPaymentTask) error {
	// Here is the actual payment processing logic
	// ...
	deps.Log().Channel("queue_log").Info(fmt.Sprintf("Processed payment for Order: %d", payload.OrderID))

	return nil
}
