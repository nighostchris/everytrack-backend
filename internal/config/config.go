package config

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v9"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	// Web Server
	WebServerHost   string   `env:"WEB_SERVER_HOST,notEmpty"`
	WebServerPort   int16    `env:"WEB_SERVER_PORT,notEmpty"`
	DomainWhitelist []string `env:"DOMAIN_WHITELIST,notEmpty"`
	// Postgres
	Database string `env:"DATABASE,notEmpty"`
	// Authentication Token
	TokenExpiryInHour int    `env:"TOKEN_EXPIRY_IN_HOUR,notEmpty"`
	AccessTokenSecret string `env:"ACCESS_TOKEN_SECRET,notEmpty"`
	// Logger
	LogLevel string `env:"LOG_LEVEL,notEmpty"`
	// External API
	TwelveDataApiKey string `env:"TWELVE_DATA_API_KEY,notEmpty"`
}

func New() *Config {
	logTemplate := "{\"level\":\"%s\",\"timestamp\":\"%s\",\"function\":\"github.com/nighostchris/everytrack-backend/internal/config.New\",\"message\":\"%s\"}\n"
	config := Config{}
	fmt.Printf(logTemplate, "info", time.Now().Format(time.RFC3339Nano), "initializing environment variables config")

	// Try to parse environment variables into config
	if err := env.Parse(&config); err != nil {
		fmt.Printf(logTemplate, "error", time.Now().Format(time.RFC3339Nano), err.Error())
		os.Exit(1)
	}
	fmt.Printf(logTemplate, "info", time.Now().Format(time.RFC3339Nano), "initialized environment variables config")

	return &config
}
