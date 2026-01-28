package configs

type RabbitMQConfig struct {
	URL                  string `koanf:"url"`
	Exchange             string `koanf:"exchange"`
	QueueDurable         bool   `koanf:"queue_durable"`
	QueueAutoDelete      bool   `koanf:"queue_auto_delete"`
	PrefetchCount        int    `koanf:"prefetch_count"`
	ReconnectDelay       string `koanf:"reconnect_delay"`
	MaxReconnectAttempts int    `koanf:"max_reconnect_attempts"`
}
