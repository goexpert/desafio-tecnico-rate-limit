package usecase

import (
	"sync"
	"time"
)

// IpRateLimiter struct to hold rate limiting data
type IpRateLimiter struct {
	requests map[string]int
	mu       sync.Mutex
	limit    int
	interval time.Duration
}

// NewIpRateLimiter creates a new rate limiter
func NewIpRateLimiter(limit int, interval time.Duration) *IpRateLimiter {
	rl := &IpRateLimiter{
		requests: make(map[string]int),
		limit:    limit,
		interval: interval,
	}
	go rl.cleanup()
	return rl
}

// Allow checks if the request is allowed
func (rl *IpRateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	requestsNow := rl.requests[ip]

	if requestsNow >= rl.limit {
		return false
	}

	rl.requests[ip]++
	return true
}

// cleanup resets the rate limit counts at regular intervals
func (rl *IpRateLimiter) cleanup() {
	for {
		time.Sleep(rl.interval)
		rl.mu.Lock()
		for k := range rl.requests {
			delete(rl.requests, k)
		}
		rl.mu.Unlock()
	}
}
