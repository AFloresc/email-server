package middleware

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

var visitors = make(map[string]*rate.Limiter)
var visitorsMu sync.Mutex

func getVisitor(key string) *rate.Limiter {
	visitorsMu.Lock()
	defer visitorsMu.Unlock()

	limiter, exists := visitors[key]
	if !exists {
		limiter = rate.NewLimiter(1, 3) // 1 req/seg, burst 3
		visitors[key] = limiter
	}
	return limiter
}

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// IP hash se calcula en handler; aquí usamos RemoteAddr directamente
		limiter := getVisitor(r.RemoteAddr)
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
