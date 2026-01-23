package bootstrap

import (
	"log"
)

func NewApp() *Config {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	return config
}
