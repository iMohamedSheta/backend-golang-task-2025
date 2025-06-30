package notify

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/hibiken/asynq"
)

type Notification interface {
	Channels() []string
	ShouldQueue() bool
	Data() map[string]any
	ScheduledAt() *time.Time // Optional
}

type Notifiable interface {
	GetNotifiableID() uint
}

type sendType int

const (
	SendTypeNow sendType = iota
	SendTypeQueue
)

type NotificationChannelHandler func(ctx context.Context, task *NotificationTask) error

type Logger interface {
	Info(msg string, fields ...any)
	Error(msg string, fields ...any)
	Warn(msg string, fields ...any)
}

type Notify struct {
	asynq     *asynq.Client
	asynqOpts []asynq.Option
	channels  map[string]NotificationChannelHandler
	mu        sync.Mutex
	log       Logger
}

func New(log Logger, queueClient *asynq.Client, queueOpts ...asynq.Option) *Notify {
	return &Notify{
		channels:  make(map[string]NotificationChannelHandler),
		log:       log,
		asynq:     queueClient,
		asynqOpts: queueOpts,
	}
}

func (n *Notify) RegisterChannels(channels map[string]NotificationChannelHandler) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.channels = channels
}

// Send - dispatches the notification to the queue if it should be queued otherwise sends it immediately
func (n *Notify) Send(notification Notification, notifiables ...Notifiable) error {
	var mode sendType
	if notification.ShouldQueue() {
		mode = SendTypeQueue
	} else {
		mode = SendTypeNow
	}

	return n.sendTasks(n.buildTasks(notification, notifiables...), mode, notification.ScheduledAt())
}

func (n *Notify) SendScheduled(at time.Time, notification Notification, notifiables ...Notifiable) error {
	return n.sendTasks(n.buildTasks(notification, notifiables...), SendTypeQueue, &at)
}

// SendNow - sends the notification immediately without dispatching it to the queue
func (n *Notify) SendNow(notification Notification, notifiables ...Notifiable) error {
	return n.sendTasks(n.buildTasks(notification, notifiables...), SendTypeNow, nil)
}

func (n *Notify) sendTasks(tasks []*NotificationTask, mode sendType, scheduleAt *time.Time) error {
	for _, task := range tasks {
		var err error

		switch {
		case mode == SendTypeNow:
			err = n.handleSendNotification(context.Background(), task)

		case scheduleAt != nil:
			err = n.dispatch(task, scheduleAt)

		default:
			err = n.dispatch(task, nil)
		}

		if err != nil {
			n.logErr(fmt.Sprintf("error sending task to %s: %s", task.Channel, err))
			return err
		}
	}
	return nil
}

func (n *Notify) buildTasks(notification Notification, notifiables ...Notifiable) []*NotificationTask {
	var tasks []*NotificationTask
	for _, notifiable := range notifiables {
		for _, ch := range notification.Channels() {
			tasks = append(tasks, &NotificationTask{
				NotificationType: getTypeName(notification),
				NotifiableType:   getTypeName(notifiable),
				NotifiableID:     notifiable.GetNotifiableID(),
				Channel:          ch,
				Data:             notification.Data(),
			})
		}
	}
	return tasks
}

// handleSendNotification - sends the notification task using the handler registered for the channel.
func (n *Notify) handleSendNotification(ctx context.Context, task *NotificationTask) error {
	if handler, ok := n.channels[task.Channel]; ok {
		if err := handler(ctx, task); err != nil {
			return err
		}
	} else {
		n.logErr(fmt.Sprintf("no handler registered for channel: %s", task.Channel))
		return fmt.Errorf("no handler registered for channel: %s", task.Channel)
	}

	return nil
}

func (n *Notify) dispatch(task *NotificationTask, scheduleAt *time.Time) error {
	if n.asynq == nil {
		n.logErr("asynq client not initialized")
		return fmt.Errorf("asynq client not initialized")
	}

	asynqTask, err := task.CreateTask()
	if err != nil {
		n.logErr(fmt.Sprintf("failed to create task: %s", err))
		return err
	}

	opts := n.asynqOpts
	if scheduleAt != nil {
		opts = append(opts, asynq.ProcessAt(*scheduleAt))
	}

	info, err := n.asynq.Enqueue(asynqTask, opts...)
	if err != nil {
		n.logErr(fmt.Sprintf("failed to enqueue task: %s", err))
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	n.logInfo(fmt.Sprintf("Task enqueued: ID=%s, Queue=%s, Type=%s", info.ID, info.Queue, info.Type))
	return nil
}

// getTypeName returns the name of the type of the given value.
func getTypeName(v any) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

func (n *Notify) logErr(message string, info ...any) {
	n.log.Error("notify: "+message, info...)
}

func (n *Notify) logInfo(message string, info ...any) {
	n.log.Info("notify: "+message, info...)
}
