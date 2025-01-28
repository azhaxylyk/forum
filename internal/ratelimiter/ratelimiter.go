package ratelimiter

import (
	"net/http"
	"time"
)

const (
	rateLimitPerSecond = 5
	burstLimit         = 10
)

type RateLimiter struct {
	tokens    int
	lastCheck time.Time
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		tokens:    burstLimit,
		lastCheck: time.Now(),
	}
}

func (rl *RateLimiter) Allow() bool {
	now := time.Now()
	elapsed := now.Sub(rl.lastCheck).Seconds()
	rl.tokens += int(elapsed * float64(rateLimitPerSecond))
	if rl.tokens > burstLimit {
		rl.tokens = burstLimit
	}
	rl.lastCheck = now

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

var limiters = make(map[string]*RateLimiter)

func GetLimiter(ip string) *RateLimiter {
	if limiter, exists := limiters[ip]; exists {
		return limiter
	}

	limiter := NewRateLimiter()
	limiters[ip] = limiter
	return limiter
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter := GetLimiter(ip)
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
