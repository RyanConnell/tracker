package httpserver

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func LoggingMiddleware(log *zap.Logger, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		log.Info(r.Method,
			zap.String("path", r.URL.Path),
			zap.Duration("duration", time.Since(start)))
	}
}
