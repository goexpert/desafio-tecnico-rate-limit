package middleware

import (
	"net"
	"net/http"

	"github.com/goexpert/rate-limit/internal/usecase"
)

// rateLimitMiddleware applies rate limiting to incoming requests
func IpRateLimitMiddleware(next http.Handler, limiter *usecase.IpRateLimiter) http.Handler {

	result := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ip := r.RemoteAddr
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		if !limiter.Allow(ip) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})

	return result
}
