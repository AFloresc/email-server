package middleware

import (
	"net/http"
	"sync"
	"time"
)

var lastRequestTime = make(map[string]time.Time)
var cooldownMu sync.Mutex

func checkCooldown(key string, cooldown time.Duration) bool {
	cooldownMu.Lock()
	defer cooldownMu.Unlock()

	lastTime, exists := lastRequestTime[key]
	if !exists {
		lastRequestTime[key] = time.Now()
		return true
	}

	if time.Since(lastTime) < cooldown {
		return false
	}

	lastRequestTime[key] = time.Now()
	return true
}

func Cooldown(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !checkCooldown(r.RemoteAddr, 30*time.Second) {
			http.Error(w, "Please wait before sending another message", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
