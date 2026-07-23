package middleware

import (
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	count     int
	resetTime time.Time
}

var (
	rateLimits = make(map[string]*rateLimiter)
	rateMu     sync.RWMutex
)

func RateLimitMiddleware(maxRequests int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			rateMu.Lock()
			now := time.Now()

			limiter, exists := rateLimits[ip]
			if !exists || now.After(limiter.resetTime) {
				rateLimits[ip] = &rateLimiter{
					count:     1,
					resetTime: now.Add(window),
				}
				rateMu.Unlock()
				next.ServeHTTP(w, r)
				return
			}

			if limiter.count >= maxRequests {
				rateMu.Unlock()
				http.Error(w, `{"error":"rate_limit_exceeded","message":"Too many requests"}`, http.StatusTooManyRequests)
				return
			}

			limiter.count++
			rateMu.Unlock()
			next.ServeHTTP(w, r)
		})
	}
}
