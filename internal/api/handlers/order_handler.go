package handlers

import (
	"taskgo/internal/api/requests"
	"taskgo/internal/api/responses"
	"taskgo/internal/helpers"
	"taskgo/internal/policies"
	"taskgo/internal/repository"
	"taskgo/internal/services"
	"taskgo/internal/tasks"
	"taskgo/pkg/errors"
	"taskgo/pkg/logger"
	"taskgo/pkg/response"
	"taskgo/pkg/utils"
	"taskgo/pkg/validate"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrderHandler struct {
	orderService *services.OrderService
	orderPolicy  *policies.OrderPolicy
}

func NewOrderHandler() *OrderHandler {
	orderRepository := repository.NewOrderRepository()
	productRepository := repository.NewProductRepository()
	inventoryRepository := repository.NewInventoryRepository()
	inventoryService := services.NewInventoryService(inventoryRepository, productRepository)
	return &OrderHandler{
		orderService: services.NewOrderService(inventoryService, orderRepository, productRepository),
		orderPolicy:  &policies.OrderPolicy{},
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) error {
	var req requests.CreateOrderRequest
	if err := utils.BindToRequestAndExtractFields(c, &req); err != nil {
		return errors.NewBadRequestBindingError("", "Failed to bind create order request", err)
	}

	authUser, authorizeErr := helpers.GetAuthUser(c)
	if authorizeErr != nil {
		return authorizeErr
	}

	// sync: Policy Check
	if !h.orderPolicy.CanCreate(authUser) {
		return errors.NewUnAuthorizedError("Unauthorized", "You are not allowed to create order", nil)
	}

	// sync: Validation
	valid, validErrorsList := validate.ValidateRequest(&req)
	if !valid {
		return errors.NewValidationError(validErrorsList)
	}

	// sync: Service Call -> create order
	order, err := h.orderService.CreateOrder(&req)
	if err != nil {
		if validError, ok := errors.AsValidationError(err); ok {
			return validError
		}

		if serverError, ok := errors.AsServerError(err); ok {
			return serverError
		}

		return errors.NewServerError("Internal Server Error: Failed to create order", "Internal Server Error: Failed to create order", err)
	}

	// Async chain of tasks -> inventory check -> process payment -> order fulfillment -> after that other tasks are independent (notifications, reporting) can be handled in another way
	err = tasks.NewChain().
		Then(tasks.NewInventoryCheckChainTask(order.ID)).
		Then(tasks.NewProcessPaymentChainTask(order.ID)).
		OnQueue(tasks.QueueOrderProcessingChain).
		MaxRetries(3).
		Timeout(2 * time.Minute).
		OnSuccess(func(result interface{}) error {
			logger.Log().Info("Order processing chain completed", zap.Uint("order_id", order.ID))
			return nil
		}).
		OnFailure(func(err error) error {
			logger.Log().Error("Order processing chain failed", zap.Uint("order_id", order.ID), zap.Error(err))
			return nil
		}).
		Dispatch()

	if err != nil {
		// Dispatching failure should be handled also to mark the order status but skip for now
		logger.Log().Error("Failed to dispatch order processing chain", zap.Uint("order_id", order.ID), zap.Error(err))
		return errors.NewServerError("Internal Server Error: Failed to dispatch order processing chain", "Internal Server Error: Failed to dispatch order processing chain", err)
	}

	var orderResponse responses.OrderResponse
	orderResponse.Response(c, order)
	return nil
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	response.Json(c, "Order details retrieved successfully", gin.H{
		"order_id":   c.Param("id"),
		"product_id": "67890",
		"quantity":   2,
	}, 200)
}

func (h *OrderHandler) ListUserOrders(c *gin.Context) {
	response.Json(c, "User orders retrieved successfully", gin.H{
		"orders": []string{"Order1", "Order2", "Order3"},
	}, 200)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	response.Json(c, "Order cancelled successfully", gin.H{
		"order_id": c.Param("id"),
	}, 200)
}

func (h *OrderHandler) GetOrderStatus(c *gin.Context) {
	response.Json(c, "Order status retrieved successfully", gin.H{
		"order_id": c.Param("id"),
		"status":   "Processing",
	}, 200)
}
