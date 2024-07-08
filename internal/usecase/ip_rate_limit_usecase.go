package usecase

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/goexpert/rate-limit/internal/database"
)

// IpRateLimiter struct to hold rate limiting data
type IpRateLimiter struct {
	ctx      context.Context
	requests map[string]int
	mu       sync.Mutex
	limit    int
	interval time.Duration
	client   *redis.Client
}

// NewIpRateLimiter creates a new rate limiter
func NewIpRateLimiter(ctx context.Context, limit int, interval time.Duration, client *redis.Client) *IpRateLimiter {
	rl := &IpRateLimiter{
		ctx:      ctx,
		requests: make(map[string]int),
		limit:    limit,
		interval: interval,
		client:   client,
	}
	go rl.cleanup()
	return rl
}

// Allow checks if the request is allowed
func (rl *IpRateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	result, err := rl.client.Get(rl.ctx, ip).Result()
	var ipRequests *database.IpRequests
	json.Unmarshal([]byte(result), &ipRequests)

	if err != nil {
		json, _ := json.Marshal(database.NewRequest(ip, 1))
		rl.client.Set(rl.ctx, ip, json, 0)
		return true
	}

	//	requestsNow := rl.requests[ip]
	requestsNow := ipRequests.Qty

	if requestsNow >= rl.limit {
		return false
	}

	//	rl.requests[ip]++

	requestsNow++
	json, _ := json.Marshal(database.NewRequest(ip, requestsNow))
	rl.client.Set(rl.ctx, ip, json, 0)
	return true
}

// cleanup resets the rate limit counts at regular intervals
func (rl *IpRateLimiter) cleanup() error {
	for {
		rl.mu.Lock()
		// for k := range rl.requests {
		// 	delete(rl.requests, k)
		// }
		pattern := "*.*.*.*"
		// cursor := uint64(0)

		lista, err := rl.client.Keys(rl.ctx, pattern).Result()
		if err != nil {
			return err
		}
		if len(lista) > 0 {
			for _, key := range lista {
				_, err = rl.client.Del(rl.ctx, key).Result()
				if err != nil {
					return err
				}
			}
		}
		rl.mu.Unlock()
		time.Sleep(rl.interval)
	}
}
