package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	Port string `env:"PORT"`
	DSN  string `env:"DB_DSN"`
}

func GetDefaultConfig(logger *log.Logger) *Config {
	config := &Config{}
	if err := env.Parse(config); err != nil {
		logger.Fatal(err)
	}

	if config.Port == "" {
		config.Port = "3000"
	}

	return config
}
