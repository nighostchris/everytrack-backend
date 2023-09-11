package main

import (
	"fmt"
	"os"

	"github.com/nighostchris/everytrack-backend/internal/app/auth"
	"github.com/nighostchris/everytrack-backend/internal/app/savings"
	"github.com/nighostchris/everytrack-backend/internal/app/settings"
	"github.com/nighostchris/everytrack-backend/internal/config"
	"github.com/nighostchris/everytrack-backend/internal/connections/postgres"
	"github.com/nighostchris/everytrack-backend/internal/connections/server"
	"github.com/nighostchris/everytrack-backend/internal/logger"
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
	settings.NewHandler(db, zapLogger, &authMiddleware).BindRoutes(app.Group("/v1/settings"))

	// Start web server
	if initWebServerError := app.Start(fmt.Sprintf("%s:%d", env.WebServerHost, env.WebServerPort)); initWebServerError != nil {
		zapLogger.Error(initWebServerError.Error())
		os.Exit(1)
	}
	zapLogger.Info(fmt.Sprintf("web server listening on %s:%d", env.WebServerHost, env.WebServerPort))
}
