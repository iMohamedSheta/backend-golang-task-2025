package enums

type NotificationType string

const (
	NotificationTypeOrder     NotificationType = "order"
	NotificationTypePayment   NotificationType = "payment"
	NotificationTypeSystem    NotificationType = "system"
	NotificationTypeInventory NotificationType = "inventory"
)

type NotificationStatus string

const (
	NotificationStatusPending NotificationStatus = "pending"
	NotificationStatusSent    NotificationStatus = "sent"
	NotificationStatusFailed  NotificationStatus = "failed"
	NotificationStatusRead    NotificationStatus = "read"
)

type NotificationChannel string

const (
	NotificationChannelDatabase  NotificationChannel = "database"
	NotificationChannelEmail     NotificationChannel = "email"
	NotificationChannelSMS       NotificationChannel = "sms"
	NotificationChannelWebSocket NotificationChannel = "ws"
)
