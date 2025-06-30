package notify

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const TypeSendNotification = "send:notification" // Notification task type

// NotificationTask implement chainq.Task interface also it's used as payload for task
type NotificationTask struct {
	NotificationType string         `json:"notification_type"`
	NotifiableType   string         `json:"notifiable_type"`
	NotifiableID     uint           `json:"notifiable_id"`
	Data             map[string]any `json:"data"`
	Channel          string         `json:"channel"`
}

func (t *NotificationTask) GetTaskType() string {
	return TypeSendNotification
}

func (t *NotificationTask) GetPayload() any {
	return t // Return itself as payload
}

func (t *NotificationTask) CreateTask(opts ...asynq.Option) (*asynq.Task, error) {
	payload, err := json.Marshal(t.GetPayload())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return asynq.NewTask(t.GetTaskType(), payload, opts...), nil
}

/*
|------------------------------------------
|  Task handler: NotificationHandler
|------------------------------------------
*/
type NotificationHandler struct {
	Notify *Notify
}

// Return a new payment task Handler
func NewNotificationHandler(notify *Notify) *NotificationHandler {
	return &NotificationHandler{
		Notify: notify,
	}
}

// Handler method for the payment task implement Handler interface
func (h *NotificationHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var notificationTask NotificationTask
	if err := json.Unmarshal(t.Payload(), &notificationTask); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	return h.Notify.handleSendNotification(ctx, &notificationTask)
}
