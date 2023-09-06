package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nighostchris/everytrack-backend/internal/app/auth"
	"github.com/nighostchris/everytrack-backend/internal/config"
	"github.com/nighostchris/everytrack-backend/internal/connections/postgres"
	"github.com/nighostchris/everytrack-backend/internal/connections/server"
)

func main() {
	// Initialize environment variable configs
	env, initConfigError := config.New()

	if initConfigError != nil {
		fmt.Println(initConfigError.Error())
		os.Exit(1)
	}

	// Establish database connection
	db, initPgError := postgres.New(env.Database)

	if initPgError != nil {
		fmt.Println(initPgError.Error())
		os.Exit(1)
	}

	dbPingError := db.Ping(context.Background())

	if dbPingError != nil {
		fmt.Println(dbPingError.Error())
	}

	// Initialize web server
	app := server.New()

	// Define routes for server
	auth.New(db).BindRoutes(app.Group("/v1/auth"))

	// Start web server
	if initWebServerError := app.Start(fmt.Sprintf("%s:%d", env.WebServerHost, env.WebServerPort)); initWebServerError != nil {
		fmt.Println(initWebServerError.Error())
		os.Exit(1)
	}
}
