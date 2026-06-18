package handler

import (
	"net/http"

	"github.com/vgadzh/telegram-message-collector/internal/observability/logctx"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	logger := logctx.FromContext(r.Context())
	logger.Info("health check")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
