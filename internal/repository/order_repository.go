package repository

import (
	"taskgo/internal/database/models"
	"taskgo/pkg/database"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository() *OrderRepository {
	db := database.GetDB()
	return &OrderRepository{
		db: db,
	}
}

// Create a new order with order items
func (r *OrderRepository) CreateWithOrderItems(order *models.Order, orderItems []*models.OrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create the order
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		// assign the order id to the order items
		for _, item := range orderItems {
			item.OrderID = order.ID
		}

		// Bulk insert the order items
		if err := tx.Create(&orderItems).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *OrderRepository) GetOrderWithOrderItems(orderID uint) (*models.Order, error) {
	var order models.Order
	if err := r.db.Preload("OrderItems").First(&order, orderID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}
