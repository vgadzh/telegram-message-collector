package http

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/vgadzh/telegram-message-collector/internal/auth"
	"github.com/vgadzh/telegram-message-collector/internal/config"
	"github.com/vgadzh/telegram-message-collector/internal/http/handler"
	"github.com/vgadzh/telegram-message-collector/internal/http/middleware"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
}

func New(cfg *config.Config, jwtService *auth.JWTService, loginService *auth.LoginService, logger *zap.Logger) *Server {
	healthHandler := handler.NewHealthHandler()
	loginHandler := handler.NewLoginHandler(loginService, jwtService)
	sessionHandler := handler.NewSessionHandler()

	mux := http.NewServeMux()
	// public
	mux.HandleFunc("GET /health/live", healthHandler.Live)
	mux.HandleFunc("POST /login", loginHandler.Login)

	// private
	mux.Handle("POST /sessions/send-code", middleware.JWT(jwtService)(http.HandlerFunc(sessionHandler.SendCode)))

	handlerChain := middleware.Chain(mux, middleware.Logging(logger))
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
