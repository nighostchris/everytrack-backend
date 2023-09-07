package main

import (
	"fmt"
	"os"

	"github.com/nighostchris/everytrack-backend/internal/app/auth"
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
	app := server.New(env.DomainWhitelist)

	// Define routes for server
	auth.NewHandler(db, &utils.TokenUtils{Env: env, Logger: zapLogger}).BindRoutes(app.Group("/v1/auth"))

	// Start web server
	if initWebServerError := app.Start(fmt.Sprintf("%s:%d", env.WebServerHost, env.WebServerPort)); initWebServerError != nil {
		fmt.Println(initWebServerError.Error())
		os.Exit(1)
	}
}
