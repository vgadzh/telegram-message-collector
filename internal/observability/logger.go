package observability

import (
	"fmt"

	"github.com/vgadzh/telegram-message-collector/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	var zcfg zap.Config

	switch cfg.App.Env {
	case "local", "dev":
		zcfg = zap.NewDevelopmentConfig()
	default:
		zcfg = zap.NewProductionConfig()
	}

	level, err := zapcore.ParseLevel(cfg.App.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("parse log level: %w", err)
	}

	zcfg.Level = zap.NewAtomicLevelAt(level)

	zcfg.EncoderConfig.TimeKey = "@timestamp"
	zcfg.EncoderConfig.MessageKey = "message"
	zcfg.EncoderConfig.LevelKey = "log.level"
	zcfg.EncoderConfig.CallerKey = "log.origin.file"

	zcfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zcfg.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	zcfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	logger, err := zcfg.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}

	logger = logger.With(
		zap.String("service.name", "telegram-message-collector"),
		zap.String("service.environment", cfg.App.Env),
		zap.String("host.name", cfg.App.NodeID),
	)

	return logger, nil
}
