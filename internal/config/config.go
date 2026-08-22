package config

import "os"

type Config struct {
	RabbitMq string
}

func Load() *Config {
	return &Config{
		RabbitMq: os.Getenv("RABBIT_MQ_LINK"),
	}
}
