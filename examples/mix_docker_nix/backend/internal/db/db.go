package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/briheet/nozarashi/examples/mix_docker_nix/backend/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	pool *pgxpool.Pool
}

func NewClient(ctx context.Context, cfg *config.Config) (*Client, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DB.DatabaseURL)
	if err != nil {
		return nil, err
	}
	poolConfig.MaxConns = cfg.DB.MaxConns
	poolConfig.MinConns = cfg.DB.MinConns
	poolConfig.MaxConnLifetime = cfg.DB.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.DB.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	// Dependency ordering starts PostgreSQL first, but it may still be initializing.
	pingCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	for {
		if err := pool.Ping(pingCtx); err == nil {
			return &Client{pool: pool}, nil
		}

		select {
		case <-pingCtx.Done():
			pool.Close()
			return nil, fmt.Errorf("connect to PostgreSQL: %w", pingCtx.Err())
		case <-time.After(500 * time.Millisecond):
		}
	}
}

func (c *Client) Conn() *pgxpool.Pool { return c.pool }

func (c *Client) ExecuteTx(ctx context.Context, fn func(pgx.Tx) error) (resultErr error) {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			resultErr = errors.Join(resultErr, fmt.Errorf("rollback transaction: %w", err))
		}
	}()
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (c *Client) Close() {
	c.pool.Close()
}
