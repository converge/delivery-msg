package config

import (
	"github.com/caarlos0/env/v8"
	"github.com/charmbracelet/log"
)

type Config struct {
	DatabaseUrl  string `env:"DATABASE_URL"`
	NATSUser     string `env:"NATS_USER,required"`
	NATSPassword string `env:"NATS_PASSWORD,required"`
}

func ReadConfig() *Config {
	config := Config{}
	if err := env.Parse(&config); err != nil {
		log.Error(err)
	}

	return &config
}
