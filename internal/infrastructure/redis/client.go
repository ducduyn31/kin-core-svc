package redis

import (
	"context"
	"fmt"

	"github.com/danielng/kin-core-svc/internal/config"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

type ClientOptions struct {
	EnableTracing bool
}

func NewClient(ctx context.Context, cfg config.RedisConfig, opts ...ClientOptions) (*Client, error) {
	var opt ClientOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	redisOpt, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis URL: %w", err)
	}

	redisOpt.MaxRetries = cfg.MaxRetries
	redisOpt.PoolSize = cfg.PoolSize
	redisOpt.MinIdleConns = cfg.MinIdleConns
	redisOpt.DialTimeout = cfg.DialTimeout
	redisOpt.ReadTimeout = cfg.ReadTimeout
	redisOpt.WriteTimeout = cfg.WriteTimeout

	client := redis.NewClient(redisOpt)

	if opt.EnableTracing {
		if err := redisotel.InstrumentTracing(client); err != nil {
			return nil, fmt.Errorf("failed to instrument redis tracing: %w", err)
		}
	}

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &Client{Client: client}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	return c.Client.Ping(ctx).Err()
}

func (c *Client) Name() string {
	return "redis"
}
