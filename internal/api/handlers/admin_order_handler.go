package handlers

import (
	"taskgo/pkg/response"

	"github.com/gin-gonic/gin"
)

// TODO implement this handler
type AdminOrderHandler struct {
}

func (h *AdminOrderHandler) ListAllOrders(c *gin.Context) {
	response.Json(c, "All orders retrieved successfully", gin.H{
		"orders": []string{"Order1", "Order2", "Order3"},
	}, 200)
}

func (h *AdminOrderHandler) GetOrderDetails(c *gin.Context) {
	response.Json(c, "Order details retrieved successfully", gin.H{
		"order_id":   c.Param("id"),
		"product_id": "67890",
		"quantity":   2,
	}, 200)
}

func (h *AdminOrderHandler) UpdateOrderStatus(c *gin.Context) {
	response.Json(c, "Order status updated successfully", gin.H{
		"order_id": c.Param("id"),
		"status":   "Shipped",
	}, 200)
}

func (h *AdminOrderHandler) DailySalesReport(c *gin.Context) {
	response.Json(c, "Daily sales report generated successfully", gin.H{
		"date":   "2023-10-01",
		"sales":  1500.00,
		"orders": 30,
	}, 200)
}

func (h *AdminOrderHandler) LowStockAlerts(c *gin.Context) {
	response.Json(c, "Low stock alerts retrieved successfully", gin.H{
		"low_stock_items": []string{"Product1", "Product2"},
	}, 200)
}
