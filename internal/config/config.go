package config

import (
	"os"
)

type Config struct {
	RabbitMq string
	Host     string
	Port     string
}

func Load() *Config {
	return &Config{
		RabbitMq: os.Getenv("RABBIT_MQ_LINK"),
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("PORT"),
	}
}
