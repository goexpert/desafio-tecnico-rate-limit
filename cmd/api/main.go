package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/goexpert/rate-limit/internal/usecase"
	"github.com/goexpert/rate-limit/internal/web/handler"
	"github.com/goexpert/rate-limit/internal/web/middleware"
)

func main() {

	limiter := usecase.NewIpRateLimiter(10, time.Minute)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", handler.HelloWorldHandler)

	err := http.ListenAndServe(":8080", middleware.IpRateLimitMiddleware(mux, limiter))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
