package handler

import (
	"errors"
	"net/http"

	"github.com/vgadzh/telegram-message-collector/internal/auth"
	"github.com/vgadzh/telegram-message-collector/internal/httpx"
	"github.com/vgadzh/telegram-message-collector/internal/observability/logctx"
	"go.uber.org/zap"
)

type LoginHandler struct {
	loginService *auth.LoginService
	jwtService   *auth.JWTService
}

func NewLoginHandler(loginService *auth.LoginService, jwtService *auth.JWTService) *LoginHandler {
	return &LoginHandler{
		loginService: loginService,
		jwtService:   jwtService,
	}
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	logger := logctx.FromContext(r.Context())

	var req loginRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		logger.Debug("decode request", zap.Error(err))
		httpx.BadRequest(w, "invalid request body")
		return
	}

	u, err := h.loginService.Authenticate(r.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			logger.Debug("authentication failed")
			httpx.Unauthorized(w)
			return
		}
		logger.Error("login failed", zap.Error(err))
		httpx.InternalServerError(w)
		return
	}

	token, err := h.jwtService.GenerateToken(u.ID, string(u.Role))
	if err != nil {
		logger.Error("generate token failed", zap.Error(err))
		httpx.InternalServerError(w)
		return
	}

	resp := loginResponse{Token: token}
	if err := httpx.JSON(w, http.StatusOK, resp); err != nil {
		logger.Error("encode response", zap.Error(err))
	}
}
