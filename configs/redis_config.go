package configs

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
