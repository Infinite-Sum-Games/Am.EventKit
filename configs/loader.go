package configs

import (
	"fmt"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const CONFIG_FILE = "config.toml"

var Env *Config

func LoadConfig() (*Config, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider("config.toml"), toml.Parser()); err != nil {
		return nil, fmt.Errorf("failed to load config.toml: %w", err)
	}

	var config Config
	if err := k.Unmarshal("", &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
