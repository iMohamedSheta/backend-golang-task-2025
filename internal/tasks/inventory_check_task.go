package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"taskgo/internal/repository"
	"taskgo/internal/services"
	"taskgo/pkg/logger"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

// InventoryCheckTask implement ChainableTask interface also it's used as payload for task
type InventoryCheckTask struct {
	OrderID uint `json:"order_id"`
}

func NewInventoryCheckChainTask(orderID uint) *InventoryCheckTask {
	return &InventoryCheckTask{OrderID: orderID}
}

func (t *InventoryCheckTask) GetTaskType() string {
	return TypeInventoryCheck
}

func (t *InventoryCheckTask) GetPayload() interface{} {
	return *t // Return itself as payload
}

func (t *InventoryCheckTask) CreateTask() (*asynq.Task, error) {
	payload, err := json.Marshal(t.GetPayload())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(t.GetTaskType(), payload), nil
}

//------------------------------------------------------------------------------------------------------------------

// Dispatch a new inventory check task
func NewInventoryCheckTask(orderID uint) (*asynq.Task, error) {
	payload, err := json.Marshal(InventoryCheckTask{
		OrderID: orderID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeInventoryCheck, payload), nil
}

func DispatchInventoryCheckTask(orderID uint) error {
	task, err := NewInventoryCheckTask(orderID)
	if err != nil {
		logger.Channel("queue_log").Error("Failed to create new inventory check task", zap.Error(err))
		return err
	}

	err = Dispatch(task, asynq.Queue(QueueInventoryCheck))
	if err != nil {
		logger.Channel("queue_log").Error("Failed to dispatch task", zap.Error(err))
		return err
	}
	logger.Channel("queue_log").Info("Dispatched task", zap.String("type", task.Type()), zap.String("queue", QueueInventoryCheck))
	return nil
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
func NewInventoryCheckHandler() *InventoryCheckHandler {
	inventoryRepo := repository.NewInventoryRepository()
	productRepo := repository.NewProductRepository()
	orderRepo := repository.NewOrderRepository()
	return &InventoryCheckHandler{
		inventoryService: services.NewInventoryService(
			inventoryRepo,
			productRepo,
		),
		orderRepository: orderRepo,
	}
}

// Handler method for the payment task implement Handler interface
func (p *InventoryCheckHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var task InventoryCheckTask
	if err := json.Unmarshal(t.Payload(), &task); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	logger.Channel("queue_log").Info(fmt.Sprintf("Inventory check task received for Order:  %d", task.OrderID))

	// Your payment processing logic here
	if err := p.handle(task); err != nil {
		return fmt.Errorf("failed to check inventory for order: %w", err)
	}

	logger.Channel("queue_log").Info(fmt.Sprintf("Inventory check task processed for Order:  %d", task.OrderID))
	return nil
}

/*
|-------------------------------------------------
|  Actual task handling code goes here:
|-------------------------------------------------
*/
func (p *InventoryCheckHandler) handle(task InventoryCheckTask) error {
	// Get order with order items
	order, err := p.orderRepository.GetOrderWithOrderItems(task.OrderID)
	if err != nil {
		return fmt.Errorf("failed to get order with order items: %w", err)
	}

	// Reserve inventory for all products in one transaction
	err = p.inventoryService.ReserveInventoriesAtomic(order, order.OrderItems)
	if err != nil {
		return fmt.Errorf("failed to reserve inventory:  %w", err)
	}

	logger.Channel("queue_log").Info(fmt.Sprintf("Inventory check task processed for Order:  %d", task.OrderID))
	return nil
}
