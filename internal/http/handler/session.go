package handler

import (
	"net/http"

	"github.com/vgadzh/telegram-message-collector/internal/httpx"
)

type SessionHandler struct{}

func NewSessionHandler() *SessionHandler {
	return &SessionHandler{}
}

func (h *SessionHandler) SendCode(w http.ResponseWriter, r *http.Request) {
	httpx.Error(w, http.StatusNotImplemented, "not implemented")
}
