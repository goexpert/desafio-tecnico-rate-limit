package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/goexpert/rate-limit/internal/database/redisdb"
	"github.com/goexpert/rate-limit/internal/usecase"
	"github.com/goexpert/rate-limit/internal/web/handler"
	"github.com/goexpert/rate-limit/internal/web/middleware"
)

func main() {

	limit, err := strconv.Atoi(os.Getenv("RATELIMIT"))
	if err != nil {
		panic("RATE LIMIT not defined or invalid")
	}

	interval, err := strconv.Atoi(os.Getenv("RATELIMIT_CLEANUP_INTERVAL"))
	if err != nil {
		panic("RATE LIMIT INTERVAL not defined or invalid")
	}

	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	redisdb.Init(client, ctx)

	limiter := usecase.NewIpRateLimiter(ctx, limit, time.Second*time.Duration(interval), client)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", handler.HelloWorldHandler)

	err = http.ListenAndServe(":8080", middleware.IpRateLimitMiddleware(mux, limiter))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
