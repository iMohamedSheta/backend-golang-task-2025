package logger

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"taskgo/pkg/enums"

	"go.uber.org/zap"
)

// logger is the global logger instance
var (
	App  *zap.Logger
	once sync.Once
)

func LoadLogger(logPath string, cfg zap.Config) *zap.Logger {
	once.Do(func() {
		logDir := filepath.Dir(logPath)
		if _, err := os.Stat(logDir); os.IsNotExist(err) {
			if err := os.MkdirAll(logDir, 0755); err != nil {
				log.Fatal(enums.Red.Value() + "Failed to create log directory: " + err.Error() + enums.Reset.Value())
			}
		}

		logger, err := cfg.Build()
		if err != nil {
			panic("Failed to build logger: " + err.Error())
		}
		App = logger
	})

	return App
}

// // Get the logger instance
// func Log() *zap.Logger {
// 	return App
// }
