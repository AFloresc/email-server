package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"email-server/internal/metrics"
)

var blockedUserAgents = []string{
	"curl",
	"wget",
	"python",
	"python-requests",
	"httpclient",
	"java",
	"libwww",
	"go-http-client",
	"bot",
	"spider",
	"crawler",
}

func isSuspiciousUserAgent(ua string) bool {
	ua = strings.ToLower(ua)
	for _, blocked := range blockedUserAgents {
		if strings.Contains(ua, blocked) {
			return true
		}
	}
	return false
}

func UserAgent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")

		if ua == "" || isSuspiciousUserAgent(ua) {
			slog.Warn("Suspicious User-Agent blocked", "userAgent", ua)
			metrics.Inc(&metrics.M.UserAgentBlocks)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
