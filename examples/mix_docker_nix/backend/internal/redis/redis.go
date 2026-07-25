package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/briheet/nozarashi/examples/mix_docker_nix/backend/internal/config"
	redisclient "github.com/redis/go-redis/v9"
)

// Client wraps the Redis client used by the backend.
type Client struct {
	client *redisclient.Client
}

// NewClient connects to Redis after its service becomes ready.
func NewClient(ctx context.Context, cfg *config.Config) (*Client, error) {
	client := redisclient.NewClient(&redisclient.Options{
		Addr:        cfg.Redis.Address,
		Password:    cfg.Redis.Password,
		DB:          cfg.Redis.Database,
		DialTimeout: time.Second,
		MaxRetries:  -1,
	})

	// Dependency ordering starts Redis first, but it may still be initializing.
	pingCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	for {
		if err := client.Ping(pingCtx).Err(); err == nil {
			return &Client{client: client}, nil
		}

		select {
		case <-pingCtx.Done():
			_ = client.Close()
			return nil, fmt.Errorf("connect to Redis: %w", pingCtx.Err())
		case <-time.After(500 * time.Millisecond):
		}
	}
}

// Ping checks the active Redis connection.
func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Close releases the Redis client resources.
func (c *Client) Close() error {
	return c.client.Close()
}
