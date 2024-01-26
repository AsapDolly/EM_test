package internal

import (
	"log"

	"github.com/caarlos0/env/v6"
)

// Config содержит ключевые параметры для работы программы.
type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	DBAddress     string `env:"DATABASE_DSN"`
}

// GetConfig читает данные из окружения и возвращает заполненный Config.
func GetConfig() Config {
	var cfg Config

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
