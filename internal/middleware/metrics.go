package middleware

import (
	"net/http"
	"time"

	"email-server/internal/metrics"
)

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		metrics.Inc(&metrics.M.TotalRequests)

		next.ServeHTTP(w, r)

		duration := time.Since(start).Milliseconds()
		metrics.AddTime(duration)
	})
}
