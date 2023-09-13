package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nighostchris/everytrack-backend/internal/app/auth"
	"github.com/nighostchris/everytrack-backend/internal/app/currency"
	"github.com/nighostchris/everytrack-backend/internal/app/savings"
	"github.com/nighostchris/everytrack-backend/internal/app/settings"
	"github.com/nighostchris/everytrack-backend/internal/config"
	"github.com/nighostchris/everytrack-backend/internal/connections/postgres"
	"github.com/nighostchris/everytrack-backend/internal/connections/server"
	"github.com/nighostchris/everytrack-backend/internal/logger"
	"github.com/nighostchris/everytrack-backend/internal/tools"
	"github.com/nighostchris/everytrack-backend/internal/utils"
)

func main() {
	// Initialize environment variable configs
	env := config.New()
	// Initialize logger
	zapLogger := logger.New(env.LogLevel)
	// Establish database connection
	db := postgres.New(env.Database)
	// Initialize web server
	app := server.New(env.DomainWhitelist, zapLogger)
	// Initialize auth middleware
	authMiddleware := server.AuthMiddleware{Logger: zapLogger, TokenUtils: &utils.TokenUtils{Env: env, Logger: zapLogger}}

	// Define routes for server
	auth.NewHandler(db, env, zapLogger, &authMiddleware).BindRoutes(app.Group("/v1/auth"))
	savings.NewHandler(db, zapLogger, &authMiddleware).BindRoutes(app.Group("/v1/savings"))
	currency.NewHandler(db, zapLogger, &authMiddleware).BindRoutes(app.Group("/v1/currency"))
	settings.NewHandler(db, zapLogger, &authMiddleware).BindRoutes(app.Group("/v1/settings"))

	// Fetch exchange rates from Github Currency API every day
	tool := tools.GithubCurrencyApi{Db: db, Logger: zapLogger}
	go func() {
		for {
			tool.FetchLatestExchangeRates()
			time.Sleep(24 * time.Hour)
		}
	}()

	// Start web server
	if initWebServerError := app.Start(fmt.Sprintf("%s:%d", env.WebServerHost, env.WebServerPort)); initWebServerError != nil {
		zapLogger.Error(initWebServerError.Error())
		os.Exit(1)
	}
}
