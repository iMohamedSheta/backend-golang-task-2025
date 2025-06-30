package adapters

import "go.uber.org/zap"

type loggerAdapter struct {
	log *zap.Logger
}

func NewLoggerAdapter(log *zap.Logger) *loggerAdapter {
	return &loggerAdapter{
		log: log,
	}
}

func (l *loggerAdapter) Error(msg string, fields ...any) {
	l.log.Error(msg, zap.Any("payload", fields))
}

func (l *loggerAdapter) Info(msg string, fields ...any) {
	l.log.Info(msg, zap.Any("payload", fields))
}

func (l *loggerAdapter) Warn(msg string, fields ...any) {
	l.log.Warn(msg, zap.Any("payload", fields))
}
