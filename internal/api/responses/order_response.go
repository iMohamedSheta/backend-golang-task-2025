package responses

import (
	"taskgo/internal/database/models"
	"taskgo/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

type OrderResponse struct {
	Message string `json:"message" example:"Order created successfully"`
	Data    struct {
		Order OrderData `json:"order"`
	} `json:"data"`
}

type OrderData struct {
	ID                uint        `json:"id"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
	UserID            uint        `json:"user_id"`
	Status            string      `json:"status"`
	TotalAmount       float64     `json:"total_amount"`
	ShippingAddress   string      `json:"shipping_address"`
	BillingAddress    string      `json:"billing_address"`
	TrackingNumber    string      `json:"tracking_number"`
	EstimatedDelivery time.Time   `json:"estimated_delivery"`
	ActualDelivery    time.Time   `json:"actual_delivery"`
	Notes             string      `json:"notes"`
	OrderItems        []OrderItem `json:"order_items"`
}

type OrderItem struct {
	ID         uint      `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	OrderID    uint      `json:"order_id"`
	ProductID  uint      `json:"product_id"`
	Quantity   int       `json:"quantity"`
	UnitPrice  float64   `json:"unit_price"`
	TotalPrice float64   `json:"total_price"`
	Discount   float64   `json:"discount"`
	Tax        float64   `json:"tax"`
	Status     string    `json:"status"`
}

// Return Order response
func (r *OrderResponse) Response(c *gin.Context, order *models.Order) {
	r.Message = "Order created successfully"

	r.Data.Order = OrderData{
		ID:                order.ID,
		CreatedAt:         order.CreatedAt,
		UpdatedAt:         order.UpdatedAt,
		UserID:            order.UserID,
		Status:            string(order.Status),
		TotalAmount:       order.TotalAmount,
		ShippingAddress:   order.ShippingAddress,
		BillingAddress:    order.BillingAddress,
		TrackingNumber:    order.TrackingNumber,
		EstimatedDelivery: order.EstimatedDelivery,
		ActualDelivery:    order.ActualDelivery,
		Notes:             order.Notes,
		OrderItems:        mapOrderItems(order.OrderItems),
	}

	response.Json(c, r.Message, r.Data, 200)
}

func mapOrderItems(items []models.OrderItem) []OrderItem {
	var orderItems []OrderItem
	for _, item := range items {
		orderItems = append(orderItems, OrderItem{
			ID:         item.ID,
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
			OrderID:    item.OrderID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			TotalPrice: item.TotalPrice,
			Discount:   item.Discount,
			Tax:        item.Tax,
			Status:     item.Status,
		})
	}
	return orderItems
}
