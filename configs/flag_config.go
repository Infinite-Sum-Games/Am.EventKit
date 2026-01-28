package configs

type FlagsmithConfig struct {
	APIURL           string `koanf:"api_url"`
	EnvironmentKey   string `koanf:"environment_key"`
	RequestTimeout   int    `koanf:"request_timeout"`
	RefreshInterval  int    `koanf:"refresh_interval"`
	DefaultToEnabled bool   `koanf:"default_to_enabled"`
}
