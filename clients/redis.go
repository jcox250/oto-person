package clients

import (
	"context"

	redis "github.com/go-redis/redis/v8"
)

// TODO:
// - Use interface for client
// - Implement constructor
// - Write tests
// - Possible tracing/metrics?
type Redis struct {
	client *redis.Client
}

func NewRedisClient() *Redis {
	return &Redis{}
}

func (r *Redis) Add(ctx context.Context, key string, value interface{}) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	return r.client.Get(ctx, key).Bytes()
}
