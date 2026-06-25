package middleware

import (
	"net/http"
	"strings"

	"github.com/vgadzh/telegram-message-collector/internal/auth"
	"github.com/vgadzh/telegram-message-collector/internal/observability/logctx"
	"go.uber.org/zap"
)

func JWT(authService *auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logctx.FromContext(r.Context())

			header := r.Header.Get("Authorization")
			if header == "" {
				logger.Debug("no authorization header")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(header, bearerPrefix) {
				logger.Debug("invalid authorization header prefix")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			token := strings.TrimSpace(strings.TrimPrefix(header, bearerPrefix))
			if token == "" {
				logger.Debug("token is empty")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := authService.ParseToken(token)
			if err != nil {
				logger.Debug("authentication failed", zap.Error(err))
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			logger = logger.With(
				zap.Int64("user.id", claims.UserID),
				zap.String("user.role", claims.Role),
			)

			p := auth.Principal{
				UserID: claims.UserID,
				Role:   claims.Role,
			}

			ctx := logctx.WithLogger(r.Context(), logger)
			ctx = auth.WithPrincipal(ctx, p)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
