package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/vgadzh/telegram-message-collector/internal/app"
	"github.com/vgadzh/telegram-message-collector/internal/config"
	"github.com/vgadzh/telegram-message-collector/internal/observability"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger, err := observability.NewLogger(cfg)
	if err != nil {
		return fmt.Errorf("create logger: %w", err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	app := app.New(ctx, cfg, logger)

	logger.Info("service starting")
	return app.Run()
}
