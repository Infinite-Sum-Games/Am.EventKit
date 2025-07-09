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
	Environment     string `mapstructure:"env"`
	Port            int    `mapstructure:"port"`
	DatabaseURL     string `mapstructure:"database_url"`
	PrivateKeyPath  string `mapstructure:"private_key_path"`  // <- flattened directly
	PublicKeyPath   string `mapstructure:"public_key_path"`   // <- flattened directly
}

var Env *EnvConfig

func LoadConfig() (*EnvConfig, error) {
	// Look for "env.toml" in current directory
	viper.SetConfigName("env")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("env.toml file either doesn't exist or is unreachable")
		}
		return nil, fmt.Errorf("unexpected error occurred while loading environment variables :%w", err)
	}

	// Default to DEVELOPMENT if ENV is unset
	activeEnv := strings.ToUpper(os.Getenv("ENV"))
	if activeEnv == "" {
		activeEnv = "DEVELOPMENT"
	}

	// Validate allowed environments
	allowedEnvs := map[string]bool{
		"PRODUCTION":  true,
		"DEVELOPMENT": true,
	}
	if !allowedEnvs[activeEnv] {
		return nil, fmt.Errorf("invalid ENV value: %s (allowed: PRODUCTION, DEVELOPMENT)", activeEnv)
	}

	// Unmarshal based on selected env block
	config := &EnvConfig{}
	if err := viper.UnmarshalKey(activeEnv, config); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables from env.toml :%w", err)
	}

	// Validate loaded config
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configurations found in env.toml :%w", err)
	}

	return config, nil
}

func validateConfig(envConfig *EnvConfig) error {
	return v.ValidateStruct(envConfig,
		v.Field(&envConfig.Environment,
			v.Required,
			v.In("PRODUCTION", "DEVELOPMENT"),
		),
		v.Field(&envConfig.Port,
			v.Required,
			v.Min(1),
			v.Max(65535),
		),
		v.Field(&envConfig.DatabaseURL,
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
		v.Field(&envConfig.PrivateKeyPath, v.Required),
		v.Field(&envConfig.PublicKeyPath, v.Required),
	)
}
