package config

import (
	"github.com/caarlos0/env/v9"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	WebServerHost string `env:"WEB_SERVER_HOST"`
	WebServerPort int16  `env:"WEB_SERVER_PORT"`
	Database      string `env:"DATABASE"`
}

func New() (Config, error) {
	config := Config{}

	if err := env.Parse(&config); err != nil {
		return config, err
	}

	return config, nil
}
