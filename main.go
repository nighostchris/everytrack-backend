package main

import (
	"fmt"
	"os"

	"github.com/nighostchris/everytrack-backend/internal/app/cron"
	"github.com/nighostchris/everytrack-backend/internal/app/handlers"
	"github.com/nighostchris/everytrack-backend/internal/config"
	"github.com/nighostchris/everytrack-backend/internal/connections/postgres"
	"github.com/nighostchris/everytrack-backend/internal/connections/server"
	"github.com/nighostchris/everytrack-backend/internal/logger"
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

	// Initialize cron jobs
	cronJobs := cron.Init(db, env, logger)
	// cronJobs.SubscribeExchangeRates()
	// cronJobs.SubscribeTwelveDataFinancialData()
	cronJobs.MonitorFuturePayments()

	// Start web server
	if initWebServerError := app.Start(fmt.Sprintf("%s:%d", env.WebServerHost, env.WebServerPort)); initWebServerError != nil {
		logger.Error(initWebServerError.Error())
		os.Exit(1)
	}
}
