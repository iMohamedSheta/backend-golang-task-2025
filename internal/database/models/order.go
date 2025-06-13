package models

import (
	"fmt"
	"strings"
	"taskgo/internal/enums"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Order is the full order details
type Order struct {
	Base
	UserID            uint              `gorm:"index;not null" json:"user_id"` // foreign key user_id
	Status            enums.OrderStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	TotalAmount       float64           `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	ShippingAddress   string            `gorm:"type:text;not null" json:"shipping_address"`
	BillingAddress    string            `gorm:"type:text;not null" json:"billing_address"`
	TrackingNumber    string            `gorm:"size:100" json:"tracking_number"`
	EstimatedDelivery time.Time         `json:"estimated_delivery"`
	ActualDelivery    time.Time         `json:"actual_delivery"`
	Notes             string            `gorm:"type:text" json:"notes"`
	OrderItems        []OrderItem       `gorm:"foreignKey:OrderID" json:"order_items"` // relationship to order items
	Payment           Payment           `gorm:"foreignKey:OrderID" json:"payment"`     // relationship to payment
	User              User              `gorm:"foreignKey:UserID" json:"user"`         // relationship to user
}

func (order *Order) BeforeCreate(tx *gorm.DB) error {
	order.TrackingNumber = order.GenerateTrackingNumber("TRN")
	return nil
}

// OrderItem represents an item in an order
type OrderItem struct {
	Base
	OrderID    uint    `gorm:"index;not null" json:"order_id"`
	ProductID  uint    `gorm:"index;not null" json:"product_id"`
	Quantity   int     `gorm:"not null" json:"quantity"`
	UnitPrice  float64 `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	TotalPrice float64 `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Discount   float64 `gorm:"type:decimal(10,2);default:0" json:"discount"`
	Tax        float64 `gorm:"type:decimal(10,2);default:0" json:"tax"`
	Status     string  `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Product    Product `gorm:"foreignKey:ProductID" json:"product"` // relationship to product
}

func (orderItem *OrderItem) CalculateUnitPrice(productPrice float64) {
	orderItem.UnitPrice = productPrice
}

func (orderItem *OrderItem) CalculateTotalPrice() {
	orderItem.TotalPrice = orderItem.UnitPrice * float64(orderItem.Quantity)
}

// GenerateTrackingNumber will generate a tracking number for the order
func (o *Order) GenerateTrackingNumber(prefix string) string {
	if prefix == "" {
		prefix = "TRN"
	}
	return strings.ToUpper(fmt.Sprintf("%s_%s", prefix, uuid.New().String()[:8]))
}
