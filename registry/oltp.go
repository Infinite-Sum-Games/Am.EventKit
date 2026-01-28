package registry

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitOLTPPool() (*pgxpool.Pool, error) {
	reg := GetSingletonObject()
	if reg == nil {
		return nil, fmt.Errorf("registry not initialized")
	}

	config := reg.GetConfig()

	dbConfig := config.Database
	dbConnectionStr := dbConfig.URL

	poolConfig, err := pgxpool.ParseConfig(dbConnectionStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// --- Pool configuration validation ---
	if dbConfig.PoolMaxConns > 0 {
		poolConfig.MaxConns = int32(dbConfig.PoolMaxConns)
	} else {
		poolConfig.MaxConns = int32(100)
	}

	if dbConfig.PoolMinConns > 0 {
		poolConfig.MinConns = int32(dbConfig.PoolMinConns)
	} else {
		poolConfig.MinConns = int32(10) // DEFAULT_MIN_CONNS
	}

	// Parse time duration strings with defaults
	if dbConfig.PoolMaxConnLifetime != "" {
		if maxConnLifeTime, err := time.ParseDuration(dbConfig.PoolMaxConnLifetime); err == nil {
			poolConfig.MaxConnLifetime = maxConnLifeTime
		} else {
			poolConfig.MaxConnLifetime = time.Hour // DEFAULT_MAX_CONN_LIFE_TIME
		}
	} else {
		poolConfig.MaxConnLifetime = time.Hour
	}

	if dbConfig.PoolMaxConnIdleTime != "" {
		if maxConnIdleTime, err := time.ParseDuration(dbConfig.PoolMaxConnIdleTime); err == nil {
			poolConfig.MaxConnIdleTime = maxConnIdleTime
		} else {
			poolConfig.MaxConnIdleTime = time.Minute * 30 // DEFAULT_MAX_CONN_IDLE_TIME
		}
	} else {
		poolConfig.MaxConnIdleTime = time.Minute * 30
	}

	if dbConfig.HealthCheckPeriod != "" {
		if healthPeriod, err := time.ParseDuration(dbConfig.HealthCheckPeriod); err == nil {
			poolConfig.HealthCheckPeriod = healthPeriod
		} else {
			poolConfig.HealthCheckPeriod = time.Minute // DEFAULT_HEALTH_PERIOD
		}
	} else {
		poolConfig.HealthCheckPeriod = time.Minute
	}

	if dbConfig.ConnectTimeout != "" {
		if connTimeout, err := time.ParseDuration(dbConfig.ConnectTimeout); err == nil {
			poolConfig.ConnConfig.ConnectTimeout = connTimeout
		} else {
			poolConfig.ConnConfig.ConnectTimeout = time.Second * 60 // DEFAULT_CONN_TIMEOUT
		}
	} else {
		poolConfig.ConnConfig.ConnectTimeout = time.Second * 60
	}

	// Actual pool initialization
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbConn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection from pool: %w", err)
	}
	defer dbConn.Release()

	if err := dbConn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database connection test failed: %w", err)
	}

	return pool, nil
}

func CloseOLTPPool(pool *pgxpool.Pool) error {
	if pool != nil {
		pool.Close()
	}
	return nil
}
