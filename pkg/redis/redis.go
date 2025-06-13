package redis

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	clients sync.Map // thread-safe map
	once    sync.Once
)

type Config struct {
	Default     string                      `json:"default"`
	Options     Options                     `json:"options"`
	Connections map[string]ConnectionConfig `json:"connections"`
}

type Options struct {
	Cluster string `json:"cluster"`
	Prefix  string `json:"prefix"`
}

type ConnectionConfig struct {
	URL      string `json:"url"`
	Host     string `json:"host"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Database int    `json:"database"`
	IsActive bool   `json:"is_active"`
	PoolSize int    `json:"pool_size"`
	Timeout  string `json:"timeout"`
}

func (c ConnectionConfig) isConnectionConfigNotActive() bool {
	return (c.URL == "" && c.Host == "") || !c.IsActive
}

// Load sets up Redis connections from configuration
func Load(config Config) error {
	var err error
	once.Do(func() {
		var defaultLoaded bool

		// Create redis clients for each connection
		for name, connConfig := range config.Connections {
			// Skip if connection config is empty/inactive
			if connConfig.isConnectionConfigNotActive() {
				log.Printf("Skipping Redis %s connection - inactive", name)
				continue
			}

			var client *redis.Client

			// Create Redis options
			opts := &redis.Options{
				Password: connConfig.Password,
				DB:       connConfig.Database,
			}

			// Set pool size if specified
			if connConfig.PoolSize > 0 {
				opts.PoolSize = connConfig.PoolSize
			}

			// Set timeout if specified
			if connConfig.Timeout != "" {
				if duration, parseErr := time.ParseDuration(connConfig.Timeout); parseErr == nil {
					opts.DialTimeout = duration
					opts.ReadTimeout = duration
					opts.WriteTimeout = duration
				}
			}

			// if url is set use it, otherwise use host:port
			if connConfig.URL != "" {
				opt, parseErr := redis.ParseURL(connConfig.URL)
				if parseErr != nil {
					log.Printf("Failed to parse %s redis URL: %v - skipping", name, parseErr)
					continue
				}
				// Merge custom options with parsed URL options
				if opts.PoolSize > 0 {
					opt.PoolSize = opts.PoolSize
				}
				if opts.DialTimeout > 0 {
					opt.DialTimeout = opts.DialTimeout
					opt.ReadTimeout = opts.ReadTimeout
					opt.WriteTimeout = opts.WriteTimeout
				}
				client = redis.NewClient(opt)
			} else {
				opts.Addr = fmt.Sprintf("%s:%d", connConfig.Host, connConfig.Port)
				client = redis.NewClient(opts)
			}

			// Test connection - skip if can't connect
			ctx := context.Background()
			if pingErr := client.Ping(ctx).Err(); pingErr != nil {
				log.Printf("Failed to ping %s redis: %v - skipping", name, pingErr)
				client.Close() // close the failed client
				continue
			}

			clients.Store(name, client)
			log.Printf("Redis %s connection established", name)

			// Track if default connection was loaded
			if name == config.Default {
				defaultLoaded = true
			}
		}

		if !defaultLoaded && config.Default != "" {
			log.Printf("Warning: Default Redis connection '%s' not found or not active", config.Default)
		}
	})
	return err
}

// GetClient returns a specific Redis client
func GetClient(name string) (*redis.Client, error) {
	value, exists := clients.Load(name)
	if !exists {
		return nil, fmt.Errorf("redis client '%s' not found or not active", name)
	}

	client, ok := value.(*redis.Client)
	if !ok {
		return nil, fmt.Errorf("invalid client type for '%s'", name)
	}

	return client, nil
}

// GetConnection returns a Redis client by connection name
func Connection(name string) (*redis.Client, error) {
	return GetClient(name)
}

// Default returns the default Redis client
func Default() (*redis.Client, error) {
	return GetClient("default")
}

// Cache returns the cache Redis client
func Cache() (*redis.Client, error) {
	return GetClient("cache")
}

// Jobs returns the jobs Redis client
func Jobs() (*redis.Client, error) {
	return GetClient("jobs")
}

// Sessions returns the sessions Redis client
func Sessions() (*redis.Client, error) {
	return GetClient("sessions")
}

// RateLimit returns the rate_limit Redis client
func RateLimit() (*redis.Client, error) {
	return GetClient("rate_limit")
}

// IsConnectionActive checks if a connection is loaded and active
func IsConnectionActive(name string) bool {
	_, exists := clients.Load(name)
	return exists
}

// GetActiveConnections returns a list of active connection names
func GetActiveConnections() []string {
	var connections []string
	clients.Range(func(key, value interface{}) bool {
		if name, ok := key.(string); ok {
			connections = append(connections, name)
		}
		return true
	})
	return connections
}

// Close closes all Redis connections
func Close() error {
	clients.Range(func(key, value interface{}) bool {
		name := key.(string)
		client := value.(*redis.Client)

		if err := client.Close(); err != nil {
			log.Printf("Error closing Redis client %s: %v", name, err)
		}
		return true // continue iteration
	})

	return nil
}

// Helper function to parse Redis connection config from map
func ParseRedisConnectionConfig(configMap map[string]any) ConnectionConfig {
	config := ConnectionConfig{}

	if url, ok := configMap["url"].(string); ok {
		config.URL = url
	}

	if host, ok := configMap["host"].(string); ok {
		config.Host = host
	}

	if password, ok := configMap["password"].(string); ok {
		config.Password = password
	}

	if port, ok := getInt(configMap, "port"); ok {
		config.Port = port
	}

	if db, ok := getInt(configMap, "database"); ok {
		config.Database = db
	}

	if isActive, ok := configMap["is_active"].(bool); ok {
		config.IsActive = isActive
	}

	if poolSize, ok := getInt(configMap, "pool_size"); ok {
		config.PoolSize = poolSize
	}

	if timeout, ok := configMap["timeout"].(string); ok {
		config.Timeout = timeout
	}

	return config
}

func getInt(configMap map[string]interface{}, key string) (int, bool) {
	if v, ok := configMap[key]; ok {
		switch val := v.(type) {
		case int:
			return val, true
		case string:
			if parsed, err := strconv.Atoi(val); err == nil {
				return parsed, true
			}
		}
	}
	return 0, false
}
