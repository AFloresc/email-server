package middleware

import (
	"net/http"
)

func BodyLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 100*1024) // 100 KB
		next.ServeHTTP(w, r)
	})
}
