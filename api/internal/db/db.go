package db

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

// Get default environment configurations for DB
func FromEnv() Config {
	return Config{
		Host:     getenvDefault("PG_HOST", "localhost"),
		Port:     getenvDefault("PG_PORT", "5432"),
		User:     getenvDefault("PG_USER", "postgres"),
		Password: getenvDefault("PG_PASSWORD", ""),
		Database: getenvDefault("PG_DB", "postgres"),
		SSLMode:  getenvDefault("PG_SSLMODE", "disable"),
	}
}

// Get environment configurations
func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// Open new connection to DB
func NewPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(cfg.User),
		url.QueryEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
	)

	pgxCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("open postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return pool, nil
}

// Initialize connection to BD
func Init(ctx context.Context) error {
	cfg := FromEnv()

	pool, err := NewPool(ctx, cfg)
	if err != nil {
		return err
	}

	Pool = pool
	return nil
}

// Close connection to DB
func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
