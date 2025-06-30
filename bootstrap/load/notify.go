package load

import (
	"taskgo/internal/adapters"
	"taskgo/internal/deps"
	"taskgo/internal/tasks"
	"taskgo/pkg/ioc"
	"taskgo/pkg/notify"
	"taskgo/pkg/utils"
	"time"

	"github.com/hibiken/asynq"
)

func InitNotify(c *ioc.Container, channelsHandlers map[string]notify.NotificationChannelHandler) {
	err := ioc.Singleton(c, func(c *ioc.Container) (*notify.Notify, error) {
		zapLogger := deps.Log().Log()
		logger := adapters.NewLoggerAdapter(zapLogger)
		queue, err := ioc.AppMake[*deps.QueueClient]()
		var client *asynq.Client
		if err != nil {
			deps.Log().Log().Error("Failed to load queue client in the ioc container: " + err.Error())
			client = nil
		} else {
			client = queue.Client
		}

		notify := notify.New(logger, client, asynq.Queue(tasks.QueueNotifications), asynq.Unique(5*time.Minute), asynq.MaxRetry(3))
		notify.RegisterChannels(channelsHandlers)
		return notify, nil
	})

	if err != nil {
		utils.PrintErr("Failed to load notify module as singleton in the ioc container: " + err.Error())
	}
}
