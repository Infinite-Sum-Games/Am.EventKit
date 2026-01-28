package configs

type OpenTelemetryConfig struct {
	ServiceName    string `koanf:"service_name"`
	ServiceVersion string `koanf:"service_version"`
	Environment    string `koanf:"environment"`
	Endpoint       string `koanf:"endpoint"`
	Insecure       bool   `koanf:"insecure"`
	Enabled        bool   `koanf:"enabled"`
}
