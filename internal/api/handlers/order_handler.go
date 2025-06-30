package handlers

import (
	"taskgo/internal/api/requests"
	"taskgo/internal/api/responses"
	"taskgo/internal/deps"
	"taskgo/internal/helpers"
	"taskgo/internal/notification"
	"taskgo/internal/policies"
	"taskgo/internal/services"
	"taskgo/internal/tasks"
	"taskgo/pkg/errors"
	"taskgo/pkg/logger"
	"taskgo/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrderHandler struct {
	Handler
	orderService *services.OrderService
	orderPolicy  *policies.OrderPolicy
	log          *logger.Manager
}

// NewOrderHandler return a new OrderHandler
func NewOrderHandler(orderService *services.OrderService, orderPolicy *policies.OrderPolicy) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		orderPolicy:  orderPolicy,
		log:          deps.Log(),
	}
}

// @Summary     Create order
// @Description Create a new order
// @Tags        Orders
// @Accept      json
// @Produce     json
//
// @Param       request  body      requests.CreateOrderRequest       true  "Create order request"
//
// @Success     200      {object}  responses.CreateOrderResponse     "Order created successfully"
// @Failure     400      {object}  response.BadRequestResponse       "Bad Request"
// @Failure     401      {object}  response.UnauthorizedResponse     "Unauthorized Action"
// @Failure     422      {object}  response.ValidationErrorResponse  "Validation Error"
// @Failure     500      {object}  response.ServerErrorResponse      "Internal Server Error"
//
// @Router      /orders [post]
func (h *OrderHandler) CreateOrder(gin *gin.Context) error {
	var req requests.CreateOrderRequest
	var err error

	if err = h.BindBodyAndExtractToRequest(gin, &req); err != nil {
		return errors.NewBadRequestBindingError("", "Failed to bind create order request", err)
	}

	authUser, authorizeErr := helpers.GetAuthUser(gin)
	if authorizeErr != nil {
		return authorizeErr
	}

	if !h.orderPolicy.CanCreate(authUser) {
		return errors.NewForbiddenError("you can't create order", "user doesn't have authorization to create a new order", nil)
	}

	if err = deps.Validator().ValidateRequest(&req); err != nil {
		return err
	}

	order, err := h.orderService.CreateOrder(gin.Request.Context(), &req)
	if err != nil {
		return err
	}

	// Async chain of tasks -> inventory check -> process payment -> order fulfillment -> after that other tasks are independent (notifications, reporting) can be handled in another way
	err = tasks.Chain().
		Then(tasks.NewInventoryCheckTask(order.ID)).
		Then(tasks.NewProcessPaymentTask(order.ID)).
		OnQueue(tasks.QueueOrderProcessingChain).
		MaxRetries(3).
		Timeout(3 * time.Minute).
		OnSuccess(func(result interface{}) error {
			h.log.Channel("default").Info("Order processing chain completed", zap.Uint("order_id", order.ID))
			// Should dispatch notification task
			err := deps.Notify().Send(notification.NewOrderCreatedNotification(order.ID), authUser)
			if err != nil {
				h.log.Channel("default").Error("Failed to dispatch notification task", zap.Uint("order_id", order.ID), zap.Error(err))
				return err
			}
			return nil
		}).
		OnFailure(func(err error) error {
			h.log.Channel("default").Error("Order processing chain failed", zap.Uint("order_id", order.ID), zap.Error(err))
			return nil
		}).
		Dispatch()

	if err != nil {
		return errors.NewServerError("Internal Server Error: Failed to process the order.", "Internal Server Error: Failed to dispatch order processing chain", err)
	}

	responses.SendCreateOrderResponse(gin, order)
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
