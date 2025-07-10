package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/spf13/viper"
)

type EnvConfig struct {
	Environment   string `mapstrucutre:"env"`
	Port          int    `mapstructure:"port"`
	DatabaseURL   string `mapstructure:"database_url"`
	RedisHost     string `mapstructure:"redis_host"`
	RedisPort     int    `mapstructure:"redis_port"`
	RedisUsername string `mapstructure:"redis_username"`
	RedisPassword string `mapstructure:"redis_password"`
}

var Env *EnvConfig

func LoadConfig() (*EnvConfig, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("env")
	viper.SetConfigFile("toml")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			return nil, fmt.Errorf("env.toml file either doesn't exist or is unreachable")
		}
		// Config file was found but another error was produced
		return nil, fmt.Errorf("unexpected error occurred while loading environment variables :%w", err)
	}

	allowedEnvs := map[string]bool{
		"PRODUCTION":  true,
		"DEVELOPMENT": true,
	}

	activeEnv := strings.ToUpper(os.Getenv("ENV"))
	if activeEnv == "" {
		activeEnv = "DEVELOPMENT"
	}
	// validating the ENV from source
	if !allowedEnvs[activeEnv] {
		return nil, fmt.Errorf("invalid ENV value: %s (allowed: PRODUCTION, DEVELOPMENT)", activeEnv)
	}

	config := &EnvConfig{}
	if err := viper.UnmarshalKey(activeEnv, config); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables from env.toml :%w", err)
	}

	// validate the config
	err := validateConfig(config)
	if err != nil {
		return nil, fmt.Errorf("invalid configurations found in .env :%w", err)
	}

	return config, nil
}

func validateConfig(env *EnvConfig) error {
	return v.ValidateStruct(env,
		v.Field(&env.Environment, v.Required, v.In("PRODUCTION", "DEVELOPMENT")),
		v.Field(&env.Port, v.Required, v.Min(1), v.Max(65535)),
		v.Field(&env.DatabaseURL,
			v.Required,
			v.Length(5, 100),
			is.URL,
			v.By(func(value any) error {
				s, _ := value.(string)
				if !strings.HasPrefix(s, "postgres") {
					return errors.New("database URL must start with 'postgres'")
				}
				if _, err := url.ParseRequestURI(s); err != nil {
					return errors.New("database URL must be a valid URI")
				}
				return nil
			}),
		),
		v.Field(&env.RedisHost, v.Required),
		v.Field(&env.RedisPort, v.Required, v.Min(1), v.Max(65535)),
		// OPTIONAL
		// v.Field(&env.RedisUsername, v.Required),
		// v.Field(&env.RedisPassword, v.Required),
	)
}
