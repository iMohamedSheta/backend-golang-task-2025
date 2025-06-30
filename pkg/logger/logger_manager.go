package logger

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
)

// Manager represents a logger manager.
// It manages the loggers and provides methods to get them.
type Manager struct {
	loggers       map[string]*zap.Logger
	mu            sync.RWMutex
	defaultLogger string
	once          sync.Once
}

// Create a new Logger Manager instance.
func NewManager() *Manager {
	return &Manager{
		loggers: make(map[string]*zap.Logger),
	}
}

// Get the logger with the given name (channel).
func (m *Manager) LoadDefault(path string, cfg zap.Config) *zap.Logger {
	m.once.Do(func() {
		if err := m.Register("default", path, cfg); err != nil {
			log.Fatal("Failed to load default logger:", err)
		}
		m.defaultLogger = "default"
	})
	return m.Channel("default")
}

// Register a new logger with the given name (channel).
func (m *Manager) Register(name, path string, cfg zap.Config) error {
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

// Returns a logger by channel (name).
func (m *Manager) Channel(name string) *zap.Logger {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.loggers[name]
}

// Sets the default logger name.
func (m *Manager) SetDefaultLogger(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.defaultLogger = name
}

// Returns the default logger.
func (m *Manager) Log() *zap.Logger {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.loggers[m.defaultLogger]
}

// Close calls Sync() on all registered loggers to flush and clean up.
// This implements the Shutdownable interface the ioc shutdown the service when the application is shutdown
func (m *Manager) Shutdown() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var firstErr error
	for name, logger := range m.loggers {
		if err := logger.Sync(); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			log.Printf("Failed to sync logger '%s': %v", name, err)
		}
	}
	return firstErr
}

// Returns a copy of the registered loggers.
func (m *Manager) GetLoggers() map[string]*zap.Logger {
	m.mu.RLock()
	defer m.mu.RUnlock()
	copied := make(map[string]*zap.Logger)
	for k, v := range m.loggers {
		copied[k] = v
	}
	return copied
}
