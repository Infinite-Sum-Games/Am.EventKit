package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

const (
	defaultMaxConns          = int32(100)
	defaultMinConns          = int32(10)
	defaultMaxConnLifetime   = time.Hour
	defaultMaxConnIdleTime   = time.Minute * 30
	defaultHealthCheckPeriod = time.Minute
	defaultConnectTimeout    = time.Second * 5
)

func InitDBPool() error {
	dbConnectionStr := EnvConfig.DatabaseURL

	dbConfig, err := pgxpool.ParseConfig(dbConnectionStr)
	if err != nil {
		return fmt.Errorf("Failed to parse database URL: %w", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return fmt.Errorf("Failed to create connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbConn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("Failed to acquire connection from pool: %w", err)
	}
	defer dbConn.Release()

	if err := dbConn.Ping(ctx); err != nil {
		return fmt.Errorf("Database connection test failed: %w", err)
	}

	dbPool = pool
	return nil
}