package logger

import (
	"fmt"

	"github.com/sudeeya/avito-assignment/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	Debug = "DEBUG"
	Info  = "INFO"
	Error = "ERROR"
	Fatal = "FATAL"
)

func SetGlobalLogger(cfg config.LogConfig) error {
	loggerCfg := zap.NewDevelopmentConfig()

	switch cfg.LogLevel {
	case Debug:
		loggerCfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case Info:
		loggerCfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case Error:
		loggerCfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case Fatal:
		loggerCfg.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	default:
		return fmt.Errorf("unknown log level: %s", cfg.LogLevel)
	}

	logger, err := loggerCfg.Build()
	if err != nil {
		return fmt.Errorf("building logger: %w", err)
	}

	zap.ReplaceGlobals(logger)

	return nil
}
