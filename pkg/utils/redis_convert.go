package utils

import (
	"errors"
	"fmt"
	"strconv"
	"taskgo/internal/deps"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

// ConvertRedisConfigToOptions converts the app redis config to a redis options struct
func ConvertRedisConfigToOptions(cfg map[string]any) (*redis.Options, error) {
	if isActive, ok := cfg["active"].(bool); ok && !isActive {
		return nil, errors.New("redis connection is not active")
	}

	// Prefer URL if provided
	if rawURL, ok := cfg["url"].(string); ok && rawURL != "" {
		opt, err := redis.ParseURL(rawURL)
		if err != nil {
			return nil, fmt.Errorf("invalid redis url: %w", err)
		}
		return opt, nil
	}

	// Fallback to manual config
	host, _ := cfg["host"].(string)
	password, _ := cfg["password"].(string)

	port := toInt(cfg["port"], 6379)
	db := toInt(cfg["database"], 0)
	poolSize := toInt(cfg["pool_size"], 0)

	opt := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	}

	if poolSize > 0 {
		opt.PoolSize = poolSize
	}

	if timeoutStr, ok := cfg["timeout"].(string); ok && timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			opt.DialTimeout = timeout
		}
	}

	return opt, nil
}

// Get Redis jobs client options
func GetRedisQueueClientOptionsForAsynq(key string) (*asynq.RedisClientOpt, error) {
	cfg := deps.Config().GetMap(key, nil)
	if cfg == nil {
		return nil, errors.New("redis connection config not found for " + key)
	}

	opts, err := ConvertRedisConfigToOptions(cfg)
	if err != nil {
		return nil, err
	}

	return &asynq.RedisClientOpt{
		Addr:      opts.Addr,
		Username:  opts.Username,
		Password:  opts.Password,
		DB:        opts.DB,
		PoolSize:  opts.PoolSize,
		TLSConfig: opts.TLSConfig,
	}, nil
}

func toInt(v any, def int) int {
	switch val := v.(type) {
	case int:
		return val
	case float64:
		return int(val)
	case string:
		i, err := strconv.Atoi(val)
		if err == nil {
			return i
		}
	}
	return def
}
