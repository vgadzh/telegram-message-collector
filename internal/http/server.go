package http

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/vgadzh/telegram-message-collector/internal/config"
	"github.com/vgadzh/telegram-message-collector/internal/http/handler"
	"github.com/vgadzh/telegram-message-collector/internal/http/middleware"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
}

func New(cfg *config.Config, logger *zap.Logger) *Server {
	healthHandler := handler.NewHealthHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", healthHandler.Live)
	handlerChain := middleware.Chain(
		mux,
		middleware.Logging(logger),
	)
	return &Server{
		httpServer: &http.Server{
			Addr:              cfg.HTTP.Host + ":" + strconv.Itoa(cfg.HTTP.Port),
			Handler:           handlerChain,
			ReadTimeout:       cfg.HTTP.ReadTimeout,
			WriteTimeout:      cfg.HTTP.WriteTimeout,
			ReadHeaderTimeout: 5 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
