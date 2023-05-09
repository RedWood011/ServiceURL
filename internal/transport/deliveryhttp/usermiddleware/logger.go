package usermiddleware

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

// LoggerMiddleware Логгер маршрутизатора
func LoggerMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("incoming request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
			)

			rw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(rw, r)

			logger.Info("outgoing response",
				zap.Int("status", rw.Status()),
				zap.Int("size", rw.BytesWritten()),
			)
		})
	}
}
