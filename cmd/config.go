package cmd

import (
	"fmt"

	"github.com/spf13/viper"
)

type EnvConfig struct {
	Environment string `mapstrucutre:"env"`
	Port        int    `mapstructure:"port"`
	DatabaseURL string `mapstructure:"database_url"`
}

var Env *EnvConfig

func LoadConfig() (*EnvConfig, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("env")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			return nil, fmt.Errorf("configuration file is not found")
		}
			// Config file was found but another error was produced
			return nil, fmt.Errorf("fatal error config file: %w", err)
	}
	config := &EnvConfig{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// validate the config here

	return config, nil
}
