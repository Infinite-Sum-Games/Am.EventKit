package configs

type Config struct {
	App           AppConfig           `koanf:"app"`
	Database      DatabaseConfig      `koanf:"database"`
	Redis         RedisConfig         `koanf:"redis"`
	RabbitMQ      RabbitMQConfig      `koanf:"rabbitmq"`
	Mailer        MailerConfig        `koanf:"mailer"`
	Payment       PaymentConfig       `koanf:"payment"`
	Flagsmith     FlagsmithConfig     `koanf:"flagsmith"`
	OpenTelemetry OpenTelemetryConfig `koanf:"opentelemetry"`
}
