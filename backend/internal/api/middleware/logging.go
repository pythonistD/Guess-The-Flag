package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func LoggerMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			path := r.URL.Path
			rawQuery := r.URL.RawQuery
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rw, r)

			latency := time.Since(start)
			clientIP := r.RemoteAddr
			method := r.Method

			if rawQuery != "" {
				path = path + "?" + rawQuery
			}

			fields := []zap.Field{
				zap.String("client_ip", clientIP),
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("status", rw.status),
				zap.String("latency", latency.String()),
			}

			switch {
			case rw.status >= 500:
				logger.Error("HTTP request", fields...)
			case rw.status >= 400:
				logger.Warn("HTTP request", fields...)
			default:
				logger.Info("HTTP request", fields...)
			}
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}
