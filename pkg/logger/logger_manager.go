package logger

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
)

type LoggerManager struct {
	loggers       map[string]*zap.Logger
	mu            sync.RWMutex
	defaultLogger string
	once          sync.Once
}

var manager = &LoggerManager{
	loggers: make(map[string]*zap.Logger),
}

func LoadDefault(path string, cfg zap.Config) *zap.Logger {
	manager.once.Do(func() {
		if err := manager.Register("default", path, cfg); err != nil {
			log.Fatal("Failed to load default logger:", err)
		}
		manager.defaultLogger = "default"
	})
	return manager.Channel("default")
}

func Register(name, path string, cfg zap.Config) error {
	return manager.Register(name, path, cfg)
}

func (m *LoggerManager) Register(name, path string, cfg zap.Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.loggers[name]; exists {
		return nil
	}

	if path != "" {
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return err
		}

		cfg.OutputPaths = append(cfg.OutputPaths, path)
	}

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	m.loggers[name] = logger
	return nil
}

func Channel(name string) *zap.Logger {
	return manager.Channel(name)
}

func (m *LoggerManager) Channel(name string) *zap.Logger {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.loggers[name]
}

func Log() *zap.Logger {
	return manager.Channel(manager.defaultLogger)
}

func Close() error {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	for _, logger := range manager.loggers {
		_ = logger.Sync()
	}
	return nil
}
