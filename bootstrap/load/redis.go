package load

import (
	"context"
	"errors"
	"taskgo/internal/deps"
	"taskgo/pkg/ioc"
	"taskgo/pkg/utils"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

func InitRedisCache(c *ioc.Container) {
	// Load cache connection
	err := ioc.Singleton(c, func(c *ioc.Container) (*deps.CacheClient, error) {
		defaultConnection := deps.Config().GetString("redis.default", "default")
		cacheClient, err := cacheClient("redis.connections." + defaultConnection)
		if err != nil {
			return nil, err
		}
		return deps.NewCacheClient(cacheClient), nil
	})

	if err != nil {
		utils.PrintErr("Failed to load redis cache module in the ioc container : " + err.Error())
	}
}

func InitRedisQueue(c *ioc.Container) {
	// load queue connection
	err := ioc.Singleton(c, func(c *ioc.Container) (*deps.QueueClient, error) {
		queueClient, err := queueClient("redis.connections.queue")
		if err != nil {
			return nil, err
		}
		return deps.NewQueueClient(queueClient), nil
	})

	if err != nil {
		utils.PrintErr("Failed to load redis queue module in the ioc container : " + err.Error())
	}
}

// loads the redis cache connection
func cacheClient(key string) (*redis.Client, error) {
	cfg := deps.Config().GetMap(key, nil)
	if cfg == nil {
		return nil, errors.New("redis connection config not found for " + key)
	}

	options, err := utils.ConvertRedisConfigToOptions(cfg)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

// loads the queue connection from the config
func queueClient(key string) (*asynq.Client, error) {
	ops, err := utils.GetRedisQueueClientOptionsForAsynq(key)
	if err != nil {
		return nil, err
	}

	client := asynq.NewClient(ops)

	if err := client.Ping(); err != nil {
		return nil, err
	}

	return client, nil
}
