package config

import (
	"github.com/caarlos0/env/v9"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	// Web Server
	WebServerHost   string   `env:"WEB_SERVER_HOST"`
	WebServerPort   int16    `env:"WEB_SERVER_PORT"`
	DomainWhitelist []string `env:"DOMAIN_WHITELIST"`
	// Postgres
	Database string `env:"DATABASE"`
	// Authentication Token
	TokenExpiryInHour int    `env:"TOKEN_EXPIRY_IN_HOUR"`
	AccessTokenSecret string `env:"ACCESS_TOKEN_SECRET"`
}

func New() (Config, error) {
	config := Config{}

	if err := env.Parse(&config); err != nil {
		return config, err
	}

	return config, nil
}
