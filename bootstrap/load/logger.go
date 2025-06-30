package load

import (
	"log"
	"taskgo/internal/deps"
	"taskgo/pkg/ioc"
	"taskgo/pkg/logger"
	"taskgo/pkg/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Initialize application loggers
func InitLogger(c *ioc.Container) {
	err := ioc.Singleton(c, func(c *ioc.Container) (*logger.Manager, error) {
		cfg := deps.Config()

		defaultChannel := cfg.GetString("log.default", "app_log")
		channels := cfg.GetMap("log.channels", nil)

		manager := logger.NewManager()
		var defaultLoaded bool

		for name, channelConfigRaw := range channels {
			if channelConfig, ok := channelConfigRaw.(map[string]any); ok {
				zapCfg := buildZapConfig(channelConfig)
				path := channelConfig["path"].(string)

				if name == defaultChannel {
					manager.LoadDefault(path, zapCfg)
					defaultLoaded = true
				} else {
					if err := manager.Register(name, path, zapCfg); err != nil {
						log.Printf("Failed to register logger %s: %v", name, err)
					}
				}
			}
		}

		ensureDefaultLoaded(manager, defaultLoaded, defaultChannel, channels)
		return manager, nil
	})

	if err != nil {
		utils.PrintErr("Failed to load logger module in the ioc container : " + err.Error())
	}
}

func buildZapConfig(channel map[string]any) zap.Config {
	levelStr, _ := channel["level"].(string)
	zapLevel := toZapLevel(levelStr)

	path, _ := channel["path"].(string)
	if path == "" {
		path = "storage/logs/app.json"
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapLevel),
		Development:      levelStr == "debug",
		Encoding:         "json",
		OutputPaths:      []string{path},
		ErrorOutputPaths: []string{path},
		EncoderConfig:    zap.NewProductionEncoderConfig(),
	}
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	return cfg
}

func toZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

func ensureDefaultLoaded(m *logger.Manager, loaded bool, defaultName string, channels map[string]any) {
	if loaded {
		return
	}

	// set first channel as default channel if there is no default channel
	for name := range channels {
		// set first channel as default then stop
		m.SetDefaultLogger(name)
		break
	}
}
