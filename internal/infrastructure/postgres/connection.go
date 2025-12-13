package postgres

import (
	"context"
	"fmt"

	"github.com/danielng/kin-core-svc/internal/config"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	write *pgxpool.Pool
	read  *pgxpool.Pool
}

type DBOptions struct {
	EnableTracing bool
}

func NewDB(ctx context.Context, cfg config.DatabaseConfig, opts ...DBOptions) (*DB, error) {
	var opt DBOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	writePool, err := newPool(ctx, cfg.WriteURL, cfg, opt.EnableTracing)
	if err != nil {
		return nil, fmt.Errorf("failed to create write pool: %w", err)
	}

	var readPool *pgxpool.Pool
	if cfg.ReadURL != "" && cfg.ReadURL != cfg.WriteURL {
		readPool, err = newPool(ctx, cfg.ReadURL, cfg, opt.EnableTracing)
		if err != nil {
			writePool.Close()
			return nil, fmt.Errorf("failed to create read pool: %w", err)
		}
	} else {
		readPool = writePool
	}

	return &DB{
		write: writePool,
		read:  readPool,
	}, nil
}

func (db *DB) Write() *pgxpool.Pool {
	return db.write
}

func (db *DB) Read() *pgxpool.Pool {
	return db.read
}

func (db *DB) Close() {
	if db.read != nil && db.read != db.write {
		db.read.Close()
	}
	if db.write != nil {
		db.write.Close()
	}
}

func (db *DB) Ping(ctx context.Context) error {
	return db.write.Ping(ctx)
}

func (db *DB) Name() string {
	return "postgres"
}

func newPool(ctx context.Context, url string, cfg config.DatabaseConfig, enableTracing bool) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.MaxIdleConns)
	poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime
	poolConfig.MaxConnIdleTime = cfg.ConnMaxIdleTime

	if enableTracing {
		poolConfig.ConnConfig.Tracer = otelpgx.NewTracer()
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
