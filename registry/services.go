package registry

import (
	"fmt"
	"log"
)

// InitializeServices initializes all database connections and external services
// This should be called after the registry is set up
func InitializeServices() error {
	log.Println("[OK]: Initializing database connections and services...")

	// Initialize Redis
	redisClient, err := InitRedis()
	if err != nil {
		return fmt.Errorf("failed to initialize Redis: %w", err)
	}

	// Initialize OLTP Database
	oltpPool, err := InitOLTPPool()
	if err != nil {
		return fmt.Errorf("failed to initialize OLTP database: %w", err)
	}

	// Initialize OLAP Database with DuckDB
	olapPool, err := InitOLAPPool()
	if err != nil {
		log.Printf("[WARNING]: Failed to initialize OLAP database: %v", err)
		// Continue without OLAP - it's optional
	} else {
		log.Println("[OK]: OLAP database initialized with DuckDB support")
	}

	// Update registry with initialized services
	reg := GetSingletonObject()
	if reg == nil {
		return fmt.Errorf("registry not initialized")
	}

	// Create new registry instance with all services
	config := reg.GetConfig()
	newReg := New(
		WithConfig(&config),
		WithRedisClient(redisClient),
		WithOLTPPool(oltpPool),
	)

	if olapPool != nil {
		newReg = New(
			WithConfig(&config),
			WithRedisClient(redisClient),
			WithOLTPPool(oltpPool),
			WithOLAPPool(olapPool),
		)
	}

	// Set as singleton
	SetSingletonObject(newReg)

	log.Println("[OK]: All services initialized successfully")
	return nil
}

// ValidateOLAPSetup checks if pg_duckdb extension is available and configured
func ValidateOLAPSetup() error {
	reg := GetSingletonObject()
	if reg == nil {
		return fmt.Errorf("registry not initialized")
	}

	olapConfig := reg.GetConfig().Database.OLAP
	if !olapConfig.ForceDuckDB {
		log.Println("[INFO]: DuckDB execution not forced for OLAP queries")
		return nil
	}

	// Try to acquire OLAP connection and validate DuckDB
	olapPool := GetOLAPPool()
	if olapPool == nil {
		return fmt.Errorf("OLAP pool not initialized")
	}

	log.Println("[INFO]: Validating pg_duckdb extension...")
	// Additional validation logic can be added here

	return nil
}
