package tasks

import (
	"context"
	"fmt"
	"taskgo/internal/deps"
	"taskgo/internal/repository"
	"taskgo/internal/services"

	"github.com/hibiken/asynq"
)

// InventoryCheckTask implement Task interface also it's used as payload for task
type InventoryCheckTask struct {
	OrderID uint `json:"order_id"`
}

func NewInventoryCheckTask(orderID uint) *InventoryCheckTask {
	return &InventoryCheckTask{OrderID: orderID}
}

func (t *InventoryCheckTask) GetTaskType() string {
	return TypeInventoryCheck
}

func (t *InventoryCheckTask) GetPayload() interface{} {
	return *t
}

func (t *InventoryCheckTask) CreateTask() (*asynq.Task, error) {
	return CreateAsynqTask(t, asynq.Queue(QueueInventoryCheck), asynq.MaxRetry(3))
}

/*
|------------------------------------------
|  Task handler: InventoryCheckHandler
|------------------------------------------
*/
type InventoryCheckHandler struct {
	inventoryService *services.InventoryService
	orderRepository  *repository.OrderRepository
}

// Return a new payment task Handler
func NewInventoryCheckHandler(inventoryService *services.InventoryService, orderRepo *repository.OrderRepository) *InventoryCheckHandler {
	return &InventoryCheckHandler{
		inventoryService: inventoryService,
		orderRepository:  orderRepo,
	}
}

// Handler method for the payment task implement Handler interface
func (h *InventoryCheckHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	return processTaskPayload(ctx, t, h.handle)
}

/*
|-------------------------------------------------
|  Actual task handling code goes here:
|-------------------------------------------------
*/
func (p *InventoryCheckHandler) handle(ctx context.Context, task *InventoryCheckTask) error {
	// Get order with order items
	order, err := p.orderRepository.GetOrderWithOrderItems(task.OrderID)
	if err != nil {
		return fmt.Errorf("failed to get order with order items: %w", err)
	}

	// Reserve inventory for all products in one transaction
	err = p.inventoryService.ReserveInventoriesAtomic(ctx, order, order.OrderItems)
	if err != nil {
		return fmt.Errorf("failed to reserve inventory:  %w", err)
	}

	deps.Log().Channel("queue_log").Info(fmt.Sprintf("Inventory check task processed for Order:  %d", task.OrderID))
	return nil
}
