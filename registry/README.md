# Registry Pattern Implementation

This package implements a centralized registry pattern for managing application-wide dependencies and services.

## Features

- **Thread-safe access** with RWMutex
- **Singleton pattern** with global getter
- **Option functions** for flexible initialization
- **Type-safe getters** for all services
- **Configuration management** integration
- **Graceful shutdown** support

## Usage

### Initialization

```go
import "github.com/Infinite-Sum-Games/Am.EventKit/registry"

// Initialize with configuration
reg := registry.New(
    registry.WithConfig(config),
    registry.WithRedisClient(redisClient),
    registry.WithOLTPPool(dbPool),
)

// Set as singleton
registry.SetSingletonObject(reg)
```

### Accessing Services

```go
// Get Redis client
redisClient := registry.GetRedisClient()

// Get database pool
dbPool := registry.GetOLTPPool()

// Get configuration
config := registry.GetConfig()

// Or use instance methods
reg := registry.GetSingletonObject()
redisClient := reg.GetRedisClient()
```

### Service Initialization

```go
// Initialize all services at once (recommended)
if err := registry.InitializeServices(); err != nil {
    log.Fatal(err)
}

// Or initialize individually
redisClient, err := registry.InitRedis()
if err != nil {
    log.Fatal(err)
}

// Initialize OLTP Database
oltpPool, err := registry.InitOLTPPool()
if err != nil {
    log.Fatal(err)
}

// Initialize OLAP Database with DuckDB
olapPool, err := registry.InitOLAPPool()
if err != nil {
    log.Printf("Warning: Failed to initialize OLAP: %v", err)
}
```

## Registry Components

### Core Registry
- `Registry` struct with thread-safe access
- Support for Redis, PostgreSQL (OLTP/OLAP), external services
- Key-value storage for application-specific data

### Database Clients
- **Redis**: go-redis/v9 client with connection pooling
- **OLTP**: pgxpool for transactional database operations
- **OLAP**: pgxpool with DuckDB integration for analytical queries (4-5 connections)

### External Services
- Mailer service interface
- Payment service interface
- Message queue interface

### Configuration
Complete configuration structure with all application settings:
- App configuration (ports, domains, etc.)
- Database settings (connection strings, pool settings)
- Redis configuration
- External service configurations
- Kill-switch settings

## Integration with Bootstrap

The registry is automatically initialized in `bootstrap.NewApp()` and integrated with the application lifecycle:

```go
func NewApp() *App {
    config, err := LoadConfig()
    // ...
    reg := registry.New(registry.WithConfig(convertConfig(config)))
    registry.SetSingletonObject(reg)
    return &App{registry: reg, appConfig: *config}
}
```

## OLAP with DuckDB Integration

The registry supports DuckDB as an OLAP engine through the pg_duckdb PostgreSQL extension:

### Configuration
```toml
[database.olap]
force_duckdb = true  # Force DuckDB execution
pool_max_conns = 5   # Small pool for analytical queries
pool_min_conns = 4
url = ""  # Uses same database, pg_duckdb handles analytics
```

### Usage
```go
// Get OLAP pool with DuckDB execution enabled
olapPool := registry.GetOLAPPool()

// Run analytical queries (automatically uses DuckDB engine)
rows, err := olapPool.Query(ctx, `
    SELECT category, COUNT(*), AVG(price)
    FROM events
    GROUP BY category
    ORDER BY COUNT(*) DESC
`)

// Validate DuckDB setup
err := registry.ValidateOLAPSetup()
```

Benefits:
- **Columnar storage** for analytical queries
- **Vectorized execution** for better performance
- **Small connection pool** (4-5 connections) for resource efficiency
- **Transparent integration** - same PostgreSQL drivers
- **Force DuckDB execution** for all analytical queries

## Migration from Global Variables

This pattern replaces the previous approach of using global variables like:
- `registry.Redis` → `registry.GetRedisClient()`
- `registry.DBPool` → `registry.GetOLTPPool()`
- `registry.DBPool` → `registry.GetOLAPPool()` (new for DuckDB analytics)

Benefits:
- Thread-safe access
- Dependency injection support
- Better testability
- Centralized lifecycle management
- Graceful shutdown handling
- Separate OLTP/OLAP database pools
- DuckDB analytics engine integration