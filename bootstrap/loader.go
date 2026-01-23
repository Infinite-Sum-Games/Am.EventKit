package bootstrap

import (
	"fmt"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	App        AppConfig        `koanf:"app"`
	Database   DatabaseConfig   `koanf:"database"`
	Redis      RedisConfig      `koanf:"redis"`
	RabbitMQ   RabbitMQConfig   `koanf:"rabbitmq"`
	Mailer     MailerConfig     `koanf:"mailer"`
	Payment    PaymentConfig    `koanf:"payment"`
	KillSwitch KillSwitchConfig `koanf:"kill_switch"`
}

type AppConfig struct {
	Env          string `koanf:"env"`
	Port         int    `koanf:"port"`
	ClientDomain string `koanf:"client_domain"`
	CookieDomain string `koanf:"cookie_domain"`
	CookieSecure bool   `koanf:"cookie_secure"`
}

type DatabaseConfig struct {
	URL                 string `koanf:"url"`
	PoolMaxConns        int    `koanf:"pool_max_conns"`
	PoolMinConns        int    `koanf:"pool_min_conns"`
	PoolMaxConnLifetime string `koanf:"pool_max_conn_lifetime"`
	PoolMaxConnIdleTime string `koanf:"pool_max_conn_idle_time"`
	HealthCheckPeriod   string `koanf:"health_check_period"`
	ConnectTimeout      string `koanf:"connect_timeout"`
}

type RedisConfig struct {
	Host         string `koanf:"host"`
	Port         string `koanf:"port"`
	Username     string `koanf:"username"`
	Password     string `koanf:"password"`
	DB           int    `koanf:"db"`
	Protocol     int    `koanf:"protocol"`
	DialTimeout  string `koanf:"dial_timeout"`
	ReadTimeout  string `koanf:"read_timeout"`
	WriteTimeout string `koanf:"write_timeout"`
}

type RabbitMQConfig struct {
	URL                  string `koanf:"url"`
	Exchange             string `koanf:"exchange"`
	QueueDurable         bool   `koanf:"queue_durable"`
	QueueAutoDelete      bool   `koanf:"queue_auto_delete"`
	PrefetchCount        int    `koanf:"prefetch_count"`
	ReconnectDelay       string `koanf:"reconnect_delay"`
	MaxReconnectAttempts int    `koanf:"max_reconnect_attempts"`
}

type MailerConfig struct {
	Host         string `koanf:"host"`
	Port         int    `koanf:"port"`
	SMTPUsername string `koanf:"smtp_username"`
	SMTPPassword string `koanf:"smtp_password"`
}

type PaymentConfig struct {
	PayUTestKey       string `koanf:"payu_test_key"`
	PayUTestSalt      string `koanf:"payu_test_salt"`
	PayUProdKey       string `koanf:"payu_prod_key"`
	PayUProdSalt      string `koanf:"payu_prod_salt"`
	PayUTestVerifyURL string `koanf:"payu_test_verify_url"`
	PayUProdVerifyURL string `koanf:"payu_prod_verify_url"`
}

type KillSwitchConfig struct {
	Features     FeaturesConfig     `koanf:"features"`
	Integrations IntegrationsConfig `koanf:"integrations"`
}

type FeaturesConfig struct {
	KillWeb          bool `koanf:"kill_web"`
	KillAdmin        bool `koanf:"kill_admin"`
	KillOrgWeb       bool `koanf:"kill_org_web"`
	KillOrgApp       bool `koanf:"kill_org_app"`
	KillLogisticsWeb bool `koanf:"kill_logistics_web"`
	KillLogisticsApp bool `koanf:"kill_logistics_app"`
}

type IntegrationsConfig struct {
	KillPayments bool `koanf:"kill_payments"`
	KillMailer   bool `koanf:"kill_mailer"`
	KillDB       bool `koanf:"kill_db"`
	KillRedis    bool `koanf:"kill_redis"`
	KillQueue    bool `koanf:"kill_queue"`
}

func LoadConfig() (*Config, error) {
	k := koanf.New(".")

	// Load config.toml
	if err := k.Load(file.Provider("config.toml"), toml.Parser()); err != nil {
		return nil, fmt.Errorf("failed to load config.toml: %w", err)
	}

	// Load kill-switch.toml
	if err := k.Load(file.Provider("kill-switch.toml"), toml.Parser()); err != nil {
		return nil, fmt.Errorf("failed to load kill-switch.toml: %w", err)
	}

	var config Config
	if err := k.Unmarshal("", &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
