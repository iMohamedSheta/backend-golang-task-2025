package bootstrap

import (
	"sync"
	"taskgo/internal/deps"
	"taskgo/internal/notification/handlers"
	"taskgo/internal/providers"
	"taskgo/internal/rules"
	"taskgo/internal/tasks"
	"taskgo/pkg/ioc"
	"taskgo/pkg/notify"
	"taskgo/pkg/ws"

	"github.com/go-playground/validator/v10"
	"github.com/hibiken/asynq"
)

/*
This file is used to register all the custom stuff in the application
like the validations rules, ...etc
*/

// Global registered tasks
var registeredTasks map[string]asynq.Handler
var once sync.Once

// Register new validations rules
var registeredRules = map[string]validator.Func{
	// Add your custom validation rules here
	"unique_db":      rules.UniqueDB,
	"exists_db":      rules.ExistsDB,
	"egyptian_phone": rules.EgyptianPhone,
}

// registerTaskHandlers defines all individual task handlers
func registerTaskHandlers() map[string]asynq.Handler {
	return map[string]asynq.Handler{
		tasks.TypeProcessPayment:   deps.App[*tasks.ProcessPaymentHandler](),
		tasks.TypeInventoryCheck:   deps.App[*tasks.InventoryCheckHandler](),
		tasks.TypeSendNotification: deps.App[*notify.NotificationHandler](),
		//...
	}
}

// registerNotificationsHandlers defines all individual notification handlers
func registerNotifyChannelsHandlers() map[string]notify.NotificationChannelHandler {
	return map[string]notify.NotificationChannelHandler{
		"database": handlers.DatabaseChannelHandler,
	}
}

// registerServiceProvider - register the services providers
func registerServiceProviders(c *ioc.Container) {
	providers.RegisterHandlers(c)
	providers.RegisterRepository(c)
	providers.RegisterServices(c)
	providers.RegisterTaskHandlers(c)
}

// registerWebSocketsChannels defines all web sockets channels
func registerWebSocketsChannels(c *ioc.Container) {
	hub := deps.WS().Hub

	// Register the user notifications websocket channel
	hub.RegisterChannel(&ws.ChannelPolicy{
		Pattern: "user_notifications.*",
		CanRead: func(userID, channel string) bool {
			return channel == "user_notifications."+userID
		},
		CanWrite: func(userID, channel string) bool {
			return true
		},
	})
}
