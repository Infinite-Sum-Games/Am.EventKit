package configs

type AppConfig struct {
	Env          string `koanf:"env"`
	Port         int    `koanf:"port"`
	MCPPort      int    `koanf:"mcp_port"`
	ClientDomain string `koanf:"client_domain"`
	CookieDomain string `koanf:"cookie_domain"`
	CookieSecure bool   `koanf:"cookie_secure"`
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
