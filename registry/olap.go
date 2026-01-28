package registry

import (
	"context"
	"fmt"
	"time"

	"github.com/Infinite-Sum-Games/Am.EventKit/configs"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitOLAPPool() (*pgxpool.Pool, error) {
	reg := GetSingletonObject()
	if reg == nil {
		return nil, fmt.Errorf("registry not initialized")
	}

	config := reg.GetConfig()
	// config is a value, not a pointer, so we don't need nil check

	// Use OLAP-specific URL or fallback to main database URL
	dbConnectionStr := getOLAPDatabaseURL(&config)

	poolConfig, err := pgxpool.ParseConfig(dbConnectionStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OLAP database URL: %w", err)
	}

	// Configure OLAP pool settings (4-5 connections as requested)
	olapConfig := config.Database.OLAP

	if olapConfig.PoolMaxConns > 0 {
		poolConfig.MaxConns = int32(olapConfig.PoolMaxConns)
	} else {
		poolConfig.MaxConns = int32(5) // Default to 5 connections
	}

	if olapConfig.PoolMinConns > 0 {
		poolConfig.MinConns = int32(olapConfig.PoolMinConns)
	} else {
		poolConfig.MinConns = int32(4) // Default to 4 connections
	}

	// Configure timeouts for analytical workloads
	if olapConfig.PoolMaxConnLifetime != "" {
		if maxConnLifeTime, err := time.ParseDuration(olapConfig.PoolMaxConnLifetime); err == nil {
			poolConfig.MaxConnLifetime = maxConnLifeTime
		} else {
			poolConfig.MaxConnLifetime = time.Hour * 4 // Default 4 hours for analytics
		}
	} else {
		poolConfig.MaxConnLifetime = time.Hour * 4
	}

	if olapConfig.PoolMaxConnIdleTime != "" {
		if maxConnIdleTime, err := time.ParseDuration(olapConfig.PoolMaxConnIdleTime); err == nil {
			poolConfig.MaxConnIdleTime = maxConnIdleTime
		} else {
			poolConfig.MaxConnIdleTime = time.Minute * 30 // Default 30 minutes
		}
	} else {
		poolConfig.MaxConnIdleTime = time.Minute * 30
	}

	if olapConfig.HealthCheckPeriod != "" {
		if healthPeriod, err := time.ParseDuration(olapConfig.HealthCheckPeriod); err == nil {
			poolConfig.HealthCheckPeriod = healthPeriod
		} else {
			poolConfig.HealthCheckPeriod = time.Minute * 2 // Default 2 minutes
		}
	} else {
		poolConfig.HealthCheckPeriod = time.Minute * 2
	}

	if olapConfig.ConnectTimeout != "" {
		if connTimeout, err := time.ParseDuration(olapConfig.ConnectTimeout); err == nil {
			poolConfig.ConnConfig.ConnectTimeout = connTimeout
		} else {
			poolConfig.ConnConfig.ConnectTimeout = time.Second * 180 // Default 3 minutes for complex queries
		}
	} else {
		poolConfig.ConnConfig.ConnectTimeout = time.Second * 180
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create OLAP connection pool: %w", err)
	}

	// Test basic connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbConn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire OLAP connection from pool: %w", err)
	}
	defer dbConn.Release()

	if err := dbConn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("OLAP database connection test failed: %w", err)
	}

	// Force DuckDB execution if enabled
	if olapConfig.ForceDuckDB {
		forceDuckCtx, forceDuckCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer forceDuckCancel()

		_, err = dbConn.Exec(forceDuckCtx, "SET duckdb.force_execution TO true")
		if err != nil {
			// Log warning but don't fail - can still work with regular PostgreSQL
			fmt.Printf("Warning: Failed to force DuckDB execution: %v\n", err)
		} else {
			fmt.Println("DuckDB execution enabled for OLAP queries")
		}

		// Verify DuckDB extension is available
		_, err = dbConn.Exec(forceDuckCtx, "SELECT 1 FROM pg_extension WHERE extname = 'duckdb'")
		if err != nil {
			fmt.Printf("Warning: pg_duckdb extension may not be available: %v\n", err)
		}
	}

	return pool, nil
}

// getOLAPDatabaseURL returns the OLAP database URL or falls back to the main database URL
func getOLAPDatabaseURL(config *configs.Config) string {
	if config.Database.OLAP.URL != "" {
		return config.Database.OLAP.URL
	}
	return config.Database.URL
}

func CloseOLAPPool(pool *pgxpool.Pool) error {
	if pool != nil {
		pool.Close()
	}
	return nil
}
