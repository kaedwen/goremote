package common

import (
	"go.uber.org/zap"
)

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	WithValues(fields ...zap.Field) Logger
	Sync() error
}

func NewLogger(cfg *Config) Logger {
	var c zap.Config
	if cfg.Logging.Development {
		c = zap.NewDevelopmentConfig()
	} else {
		c = zap.NewProductionConfig()
	}

	c.Level = cfg.Logging.AtomicLevel()

	return &zapLogger{zap.Must(c.Build())}
}

type zapLogger struct {
	*zap.Logger
}

func (l *zapLogger) WithValues(fields ...zap.Field) Logger {
	return &zapLogger{Logger: l.Logger.With(fields...)}
}
