package main

import (
	"fmt"
	"os"

	"github.com/nighostchris/everytrack-backend/internal/app/auth"
	"github.com/nighostchris/everytrack-backend/internal/config"
	"github.com/nighostchris/everytrack-backend/internal/connections/postgres"
	"github.com/nighostchris/everytrack-backend/internal/connections/server"
	"github.com/nighostchris/everytrack-backend/internal/logger"
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

	// Define routes for server
	auth.NewHandler(db, env, zapLogger).BindRoutes(app.Group("/v1/auth"))

	// Start web server
	if initWebServerError := app.Start(fmt.Sprintf("%s:%d", env.WebServerHost, env.WebServerPort)); initWebServerError != nil {
		zapLogger.Error(initWebServerError.Error())
		os.Exit(1)
	}
	zapLogger.Info(fmt.Sprintf("web server listening on %s:%d", env.WebServerHost, env.WebServerPort))
}
