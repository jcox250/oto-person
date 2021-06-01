package clients

import (
	"context"
	"fmt"

	redis "github.com/go-redis/redis/v8"
)

// Logger defines the logger required by the person service
type Logger interface {
	Info(keyvals ...interface{})
	Debug(keyvals ...interface{})
	Error(keyvals ...interface{})
}

// TODO:
// - Use interface for client
// - Implement constructor
// - Write tests
// - Possible tracing/metrics?
type Redis struct {
	client *redis.Client
	log    Logger
}

func NewRedisClient(addr string, l Logger) *Redis {
	c := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &Redis{client: c, log: l}
}

func (r *Redis) Add(ctx context.Context, key string, value interface{}) error {
	if err := r.client.Set(ctx, key, value, 0).Err(); err != nil {
		r.log.Error("msg", "failed to write key, value to redis", "key", key, "value", fmt.Sprintf("%v", value))
		return err
	}
	r.log.Debug("msg", "wrote key, value to redis", "key", key, "value", fmt.Sprintf("%v", value))
	return nil
}

func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	b, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		r.log.Error("msg", "failed to get key from redis", "key", key)
		return nil, err
	}
	r.log.Debug("msg", "got value for key from redis", "key", key, "value", string(b))
	return b, err
}
