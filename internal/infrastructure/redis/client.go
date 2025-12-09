package redis

import (
	"context"
	"fmt"

	"github.com/danielng/kin-core-svc/internal/config"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

func NewClient(ctx context.Context, cfg config.RedisConfig) (*Client, error) {
	opt, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis URL: %w", err)
	}

	opt.MaxRetries = cfg.MaxRetries
	opt.PoolSize = cfg.PoolSize
	opt.MinIdleConns = cfg.MinIdleConns
	opt.DialTimeout = cfg.DialTimeout
	opt.ReadTimeout = cfg.ReadTimeout
	opt.WriteTimeout = cfg.WriteTimeout

	client := redis.NewClient(opt)

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
