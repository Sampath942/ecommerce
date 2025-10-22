package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetValueRedis(ctx context.Context, key string, client *redis.Client) (string, error) {
	return client.Get(ctx, key).Result()
}

func SetKeyValueRedis(ctx context.Context, key string, value string, expiration time.Duration, client *redis.Client) *redis.StatusCmd {
	return client.Set(ctx, key, value, expiration)
}

func DeleteKeyRedis(ctx context.Context, key string, client *redis.Client) (int64, error) {
	return client.Del(ctx, key).Result()
}

