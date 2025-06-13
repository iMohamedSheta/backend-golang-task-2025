package tasks

const (
	// Chain orchestrator type
	TypeChainOrchestrator = "chain:orchestrator"

	// Tasks types
	TypeProcessPayment = "process:payment"
	TypeInventoryCheck = "inventory:check"
)

// Queue names
const (
	QueueDefault        = "default"
	QueueCritical       = "critical"
	QueueLow            = "low"
	QueuePayments       = "payments"
	QueueInventoryCheck = "inventory_check"

	QueueOrderProcessingChain = "order_processing_chain"
)
