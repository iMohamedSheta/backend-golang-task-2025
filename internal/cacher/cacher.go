package cacher

import (
	"context"
	"fmt"
	"time"

	pkgRedis "taskgo/pkg/redis"

	"github.com/redis/go-redis/v9"
)

const (
	RememberForever time.Duration = 0
)

// Remember tries to get the value from cache. If not found, it runs the fallback,
// stores the result in Redis, and returns it.
func Remember(ctx context.Context, key string, ttl time.Duration, fallback func() (any, error)) (any, error) {
	// Get Redis connection internally
	cache, err := pkgRedis.Default()
	if err != nil {
		// Optional: log or return fallback directly if Redis is down
		fmt.Printf("Redis unavailable, using fallback for key %s: %v\n", key, err)
		return fallback()
	}
	var value any
	// Try to get from cache
	value, err = cache.Get(ctx, key).Result()
	if err == redis.Nil {
		// Not found, run fallback
		value, err = fallback()
		if err != nil {
			return "", err
		}

		// Cache the value
		if err := cache.Set(ctx, key, value, ttl).Err(); err != nil {
			fmt.Printf("Warning: Failed to cache value for key %s: %v\n", key, err)
		}

		return value, nil
	} else if err != nil {
		// Redis error (not just missing key)
		return "", fmt.Errorf("failed to fetch from Redis: %w", err)
	}

	return value, nil
}
