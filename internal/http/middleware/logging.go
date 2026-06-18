package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vgadzh/telegram-message-collector/internal/observability/logctx"
	"go.uber.org/zap"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func newStatusRecorder(w http.ResponseWriter) *statusRecorder {

	return &statusRecorder{
		ResponseWriter: w,
		status:         http.StatusOK,
	}

}

func (w *statusRecorder) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func Logging(baseLogger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := uuid.NewString()

			logger := baseLogger.With(
				zap.String("request_id", requestID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
			)

			start := time.Now()

			logger.Info("http request started")
			sw := newStatusRecorder(w)
			ctx := logctx.WithLogger(r.Context(), logger)
			next.ServeHTTP(sw, r.WithContext(ctx))

			logger.Info(
				"http request completed",
				zap.Int("status", sw.status),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}
