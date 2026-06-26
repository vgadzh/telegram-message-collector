package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/vgadzh/telegram-message-collector/internal/auth"
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

func New(ctx context.Context, cfg *config.Config, jwtService *auth.JWTService, loginService *auth.LoginService, logger *zap.Logger) *App {
	httpServer := httpserver.New(cfg, jwtService, loginService, logger)
	a := &App{
		ctx:        ctx,
		cfg:        cfg,
		logger:     logger,
		httpServer: httpServer,
	}
	return a
}

func (a *App) Run() error {
	go func() {
		<-a.ctx.Done()
		a.logger.Info("shutdown started")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.HTTP.ShutdownTimeout)
		defer cancel()

		if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
			a.logger.Error("shutdown http server", zap.Error(err))
		}
	}()

	a.logger.Info("starting http server")
	if err := a.httpServer.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("run http server: %w", err)
	}

	a.logger.Info("service stopped")
	return nil
}
