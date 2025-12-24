package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DBPool *pgxpool.Pool

const (
	DEFAULT_MAX_CONNS          = int32(450)
	DEFAULT_MIN_CONNS          = int32(10)
	DEFAULT_MAX_CONN_LIFE_TIME = time.Hour
	DEFAULT_MAX_CONN_IDLE_TIME = time.Minute * 30
	DEFAULT_HEALTH_PERIOD      = time.Minute
	DEFAULT_CONN_TIMEOUT       = time.Second * 60
)

func InitDBPool() error {
	dbConnectionStr := Env.DatabaseURL

	dbConfig, err := pgxpool.ParseConfig(dbConnectionStr)
	if err != nil {
		return fmt.Errorf("failed to parse database URL: %w", err)
	}

	dbConfig.MaxConns = DEFAULT_MAX_CONNS
	dbConfig.MinConns = DEFAULT_MIN_CONNS
	dbConfig.MaxConnLifetime = DEFAULT_MAX_CONN_LIFE_TIME
	dbConfig.MaxConnIdleTime = DEFAULT_MAX_CONN_IDLE_TIME
	dbConfig.HealthCheckPeriod = DEFAULT_HEALTH_PERIOD
	dbConfig.ConnConfig.ConnectTimeout = DEFAULT_CONN_TIMEOUT

	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbConn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection from pool: %w", err)
	}
	defer dbConn.Release()

	if err := dbConn.Ping(ctx); err != nil {
		return fmt.Errorf("database connection test failed: %w", err)
	}

	DBPool = pool
	return nil
}
