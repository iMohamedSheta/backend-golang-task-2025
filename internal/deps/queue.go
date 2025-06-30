package deps

import (
	"fmt"
	"taskgo/pkg/ioc"

	"github.com/hibiken/asynq"
)

/*
|--------------------------------------------------------
|	Application Dependency Container Alias
|--------------------------------------------------------
*/

type QueueClient struct {
	Client *asynq.Client
}

/*
|--------------------------------------------------------
|	Application Dependency Container Calls
|--------------------------------------------------------
*/

func NewQueueClient(client *asynq.Client) *QueueClient {
	return &QueueClient{Client: client}
}

func Queue() *QueueClient {
	queue, err := ioc.AppMake[*QueueClient]()
	if err != nil {
		Log().Log().Error(fmt.Sprintf("QueueClient can't be resolved: %s", err.Error()))
		return nil
	}
	return queue
}

// This implements the Shutdownable interface the ioc shutdown the service when the application is shutdown
func (c *QueueClient) Shutdown() error {
	if c.Client != nil {
		return c.Client.Close()
	}
	return nil
}
