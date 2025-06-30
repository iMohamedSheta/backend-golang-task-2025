package deps

import (
	"fmt"
	"taskgo/pkg/ioc"

	"github.com/redis/go-redis/v9"
)

/*
|--------------------------------------------------------
|	Application Dependency Container Alias
|--------------------------------------------------------
*/

type CacheClient struct {
	Redis *redis.Client
}

/*
|--------------------------------------------------------
|	Application Dependency Container Calls
|--------------------------------------------------------
*/

func NewCacheClient(client *redis.Client) *CacheClient {
	return &CacheClient{Redis: client}
}

func Cache() *CacheClient {
	cache, err := ioc.AppMake[*CacheClient]()
	if err != nil {
		Log().Log().Error(fmt.Sprintf("CacheClient Dependency Container Error: %s", err.Error()))
		return nil
	}
	return cache
}

func (c *CacheClient) Shutdown() error {
	if c.Redis != nil {
		return c.Redis.Close()
	}
	return nil
}
