package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nighostchris/everytrack-backend/internal/app/handlers"
	"github.com/nighostchris/everytrack-backend/internal/config"
	"github.com/nighostchris/everytrack-backend/internal/connections/postgres"
	"github.com/nighostchris/everytrack-backend/internal/connections/server"
	"github.com/nighostchris/everytrack-backend/internal/logger"
	"github.com/nighostchris/everytrack-backend/internal/tools"
)

func main() {
	// Initialize environment variable configs
	env := config.New()
	// Initialize logger
	logger := logger.New(env.LogLevel)
	// Establish database connection
	db := postgres.New(env.Database)
	// Initialize web server
	app := server.New(env.DomainWhitelist, logger, env)

	// Define routes for server
	handlers := handlers.Init(db, env, logger)
	handlers.BindRoutes(app)

	// Fetch exchange rates from Github Currency API every day
	tool := tools.GithubCurrencyApi{Db: db, Logger: logger}
	go func() {
		for {
			tool.FetchLatestExchangeRates()
			time.Sleep(24 * time.Hour)
		}
	}()

	// Start web server
	if initWebServerError := app.Start(fmt.Sprintf("%s:%d", env.WebServerHost, env.WebServerPort)); initWebServerError != nil {
		logger.Error(initWebServerError.Error())
		os.Exit(1)
	}
}
