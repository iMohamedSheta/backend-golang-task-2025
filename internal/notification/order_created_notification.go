package notification

import (
	"fmt"
	"time"
)

type OrderCreatedNotification struct {
	OrderID uint
}

func NewOrderCreatedNotification(orderID uint) *OrderCreatedNotification {
	return &OrderCreatedNotification{OrderID: orderID}
}

func (n *OrderCreatedNotification) Channels() []string {
	return []string{"database"}
}

func (n *OrderCreatedNotification) ToTelegram() string {
	return "🛒 New order received! Order ID: " + fmt.Sprint(n.OrderID)
}

func (n *OrderCreatedNotification) ToDatabase() string {
	return "🛒 New order received! Order ID: " + fmt.Sprint(n.OrderID)
}

func (n *OrderCreatedNotification) ShouldQueue() bool {
	return true
}

func (n *OrderCreatedNotification) ScheduledAt() *time.Time {
	return nil
}

func (n *OrderCreatedNotification) Data() map[string]any {
	return map[string]any{
		"order_id": n.OrderID,
		"channel_messages": map[string]string{
			"database": n.ToDatabase(),
			"telegram": n.ToTelegram(),
		},
	}
}
