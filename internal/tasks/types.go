package tasks

const (
	// Tasks types
	TypeProcessPayment   = "process:payment"
	TypeInventoryCheck   = "inventory:check"
	TypeSendNotification = "send:notification"
)

// Queue names
const (
	QueueDefault        = "default"
	QueueCritical       = "critical"
	QueueLow            = "low"
	QueuePayments       = "payments"
	QueueInventoryCheck = "inventory_check"
	QueueNotifications  = "notifications"

	QueueOrderProcessingChain = "order_processing_chain"
)
