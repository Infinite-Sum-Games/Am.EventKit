package configs

type DatabaseConfig struct {
	URL                 string     `koanf:"url"`
	PoolMaxConns        int        `koanf:"pool_max_conns"`
	PoolMinConns        int        `koanf:"pool_min_conns"`
	PoolMaxConnLifetime string     `koanf:"pool_max_conn_lifetime"`
	PoolMaxConnIdleTime string     `koanf:"pool_max_conn_idle_time"`
	HealthCheckPeriod   string     `koanf:"health_check_period"`
	ConnectTimeout      string     `koanf:"connect_timeout"`
	OLAP                OLAPConfig `koanf:"olap"`
}

type OLAPConfig struct {
	URL                 string `koanf:"url"`
	ForceDuckDB         bool   `koanf:"force_duckdb"`
	PoolMaxConns        int    `koanf:"pool_max_conns"`
	PoolMinConns        int    `koanf:"pool_min_conns"`
	PoolMaxConnLifetime string `koanf:"pool_max_conn_lifetime"`
	PoolMaxConnIdleTime string `koanf:"pool_max_conn_idle_time"`
	HealthCheckPeriod   string `koanf:"health_check_period"`
	ConnectTimeout      string `koanf:"connect_timeout"`
}
