package enums

type OrderStatus string

const (
	// Order is created but not yet paid or inventory reserved.
	OrderStatusPending OrderStatus = "pending"

	// Payment succeeded and inventory has been successfully reserved.
	// Next step is to process and prepare the order.
	OrderStatusConfirmed OrderStatus = "confirmed"

	// Order is being processed (e.g., picked, packed).
	// Typically follows 'confirmed'.
	OrderStatusProcessing OrderStatus = "processing"

	// Order has been handed over to the delivery service.
	// Shipping/tracking info can be attached here.
	OrderStatusShipped OrderStatus = "shipped"

	// Customer has received the order.
	// Marks the final success state in normal flow.
	OrderStatusDelivered OrderStatus = "delivered"

	// Order was cancelled by user or system (before shipping).
	// Can happen in 'pending' or 'confirmed' status.
	OrderStatusCancelled OrderStatus = "cancelled"

	// Order has been refunded after payment (e.g., returned by customer).
	// Can happen after 'delivered' or in special cases after 'shipped'.
	OrderStatusRefunded OrderStatus = "refunded"
)
