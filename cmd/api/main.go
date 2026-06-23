package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/vgadzh/telegram-message-collector/internal/app"
	"github.com/vgadzh/telegram-message-collector/internal/bootstrap"
	"github.com/vgadzh/telegram-message-collector/internal/config"
	"github.com/vgadzh/telegram-message-collector/internal/migrate"
	"github.com/vgadzh/telegram-message-collector/internal/observability"
	"github.com/vgadzh/telegram-message-collector/internal/postgres"
	"github.com/vgadzh/telegram-message-collector/internal/repo"
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

	pg, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		return err
	}
	defer pg.Close()

	if err := migrate.Up(ctx, pg.Pool(), "./migrations"); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	userRepo := repo.NewUserRepo(pg.Pool())

	if err := bootstrap.BootstrapAdmin(ctx, userRepo, cfg.Admin); err != nil {
		return fmt.Errorf("bootstrap admin: %w", err)
	}

	app := app.New(ctx, cfg, logger)

	logger.Info("service starting")
	return app.Run()
}
