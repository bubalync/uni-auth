package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type Client struct {
	client *redis.Client

	opts *redis.Options
}

func NewRedisClient(opts ...Option) *Client {
	c := &Client{
		opts: &redis.Options{
			Addr: "localhost:6379",
			DB:   0,
		},
	}

	// Custom options
	for _, opt := range opts {
		opt(c)
	}
	c.client = redis.NewClient(c.opts)

	// todo configure otel
	// https://github.com/redis/go-redis
	//if err := errors.Join(redisotel.InstrumentTracing(rdb), redisotel.InstrumentMetrics(rdb)); err != nil {
	//	log.Fatal(err)
	//}

	return c
}

func (r *Client) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *Client) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *Client) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *Client) Close() error {
	return r.client.Close()
}
