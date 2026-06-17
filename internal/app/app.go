package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/vgadzh/telegram-message-collector/internal/config"
	httpserver "github.com/vgadzh/telegram-message-collector/internal/http"
	"go.uber.org/zap"
)

type App struct {
	ctx        context.Context
	cfg        *config.Config
	logger     *zap.Logger
	httpServer *httpserver.Server
}

func New(ctx context.Context, cfg *config.Config, logger *zap.Logger) (*App, error) {
	httpServer := httpserver.New(cfg, logger)
	a := &App{
		ctx:        ctx,
		cfg:        cfg,
		logger:     logger,
		httpServer: httpServer,
	}
	return a, nil
}

func (a *App) Run() error {
	go func() {
		a.logger.Info("http server started")

		err := a.httpServer.Run()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("http server stopped", zap.Error(err))
		}
	}()

	<-a.ctx.Done()
	a.logger.Info("shutdown started")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}

	return nil
}
