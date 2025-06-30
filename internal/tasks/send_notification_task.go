package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

// SendNotificationTask implement Task, ChainableTask interfaces also it's used as payload for task
type SendNotificationTask struct {
	NotificationType string         `json:"notification_type"`
	NotifiableType   string         `json:"notifiable_type"`
	NotifiableID     uint           `json:"notifiable_id"`
	Data             map[string]any `json:"data"`
	Channel          string         `json:"channel"`
}

func (t *SendNotificationTask) GetTaskType() string {
	return TypeSendNotification
}

func (t *SendNotificationTask) GetPayload() interface{} {
	return *t // Return itself as payload
}

func (t *SendNotificationTask) CreateTask() (*asynq.Task, error) {
	return CreateAsynqTask(t, asynq.Queue(QueueNotifications), asynq.MaxRetry(3))
}

/*
|------------------------------------------
|  Task handler: SendNotificationHandler
|------------------------------------------
*/
type SendNotificationHandler struct {
}

// Return a new payment task Handler
func NewSendNotificationHandler() *SendNotificationHandler {
	return &SendNotificationHandler{}
}

// Handler method for the payment task implement Handler interface
func (h *SendNotificationHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	return processTaskPayload(ctx, t, h.handle)
}

/*
|-------------------------------------------------
|  Actual task handling code goes here:
|-------------------------------------------------
*/
func (p *SendNotificationHandler) handle(ctx context.Context, task *SendNotificationTask) error {
	fmt.Printf("Sending notification: %s\n", task.NotificationType)
	time.Sleep(time.Second * 3)
	fmt.Printf("Notification sent:  %s\n", task.NotificationType)
	return nil
}
