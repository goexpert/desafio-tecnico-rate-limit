package redisdb

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/goexpert/rate-limit/internal/database"
)

func Init(client *redis.Client, ctx context.Context) {
	pong, err := client.Ping(context.Background()).Result()
	fmt.Println(pong, err)

	lista := database.HidrateListaTokens()
	for _, tokenLimit := range lista {
		json, _ := json.Marshal(tokenLimit)
		client.Set(ctx, tokenLimit.Name, json, 0).Err()
	}
}
