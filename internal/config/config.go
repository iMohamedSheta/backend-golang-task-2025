package config

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type Config struct {
	mu    sync.RWMutex
	store map[string]any
}

// Contain all the loaded function to load configurations
var (
	loaders   []func()
	loadersMu sync.RWMutex
)

// Global app config instance
var App = &Config{store: make(map[string]any)}

// Set a configuration value
func (c *Config) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
}

// Get a configuration value - returns error instead of panicking
func (c *Config) Get(key string) (any, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := strings.Split(key, ".")
	current := c.store
	lastIndex := len(keys) - 1

	for i, k := range keys {
		// if we are at the last key return the value
		if i == lastIndex {
			if value, ok := current[k]; ok {
				return value, nil
			}
			return nil, fmt.Errorf("config key not found: %s", key)
		}

		// check if the key inside nested map
		if value, ok := current[k].(map[string]any); ok {
			current = value
		} else {
			return nil, fmt.Errorf("config key not found: %s (failed at path: %s)", key, strings.Join(keys[:i+1], "."))
		}
	}

	return nil, fmt.Errorf("config key not found: %s", key)
}

// GetWithDefault - gets a value with a default fallback
func (c *Config) GetWithDefault(key string, defaultValue any) any {
	value, err := c.Get(key)
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *Config) GetString(key string, defaultVal string) string {
	val, err := c.Get(key)
	if err != nil {
		return defaultVal
	}
	if str, ok := val.(string); ok {
		return str
	}
	return defaultVal
}

func (c *Config) GetBool(key string, defaultVal bool) bool {
	val, err := c.Get(key)
	if err != nil {
		return defaultVal
	}
	if boolean, ok := val.(bool); ok {
		return boolean
	}
	return defaultVal
}

func (c *Config) GetArrayOfStrings(key string, defaultVal []string) []string {
	val, err := c.Get(key)
	if err != nil {
		return defaultVal
	}
	if arr, ok := val.([]string); ok {
		return arr
	}
	return defaultVal
}

func (c *Config) GetInt(key string, defaultVal int) int {
	val, err := c.Get(key)
	if err != nil {
		return defaultVal
	}
	if v, ok := val.(int); ok {
		return v
	}
	return defaultVal
}

func (c *Config) GetDuration(key string, defaultVal time.Duration) time.Duration {
	val, err := c.Get(key)
	if err != nil {
		return defaultVal
	}
	if time, ok := val.(time.Duration); ok {
		return time
	}
	return defaultVal
}

// Load configuration from environment variables
func Env(key string, defaultValue any) any {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	if defaultValue != nil {
		return defaultValue
	}

	return nil
}

// Register new function to load config set
func Register(loader func()) {
	loadersMu.Lock()
	defer loadersMu.Unlock()
	loaders = append(loaders, loader)
}

// Load all the registered function config loaders
func LoadAll() {
	loadersMu.RLock()
	defer loadersMu.RUnlock()
	for _, load := range loaders {
		load()
	}
}
