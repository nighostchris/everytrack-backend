package main

import (
	"context"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nighostchris/everytrack-backend/internal/config"
	"github.com/nighostchris/everytrack-backend/internal/connections/postgres"
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
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// TO FIX: modify it to consume whitelist domains from .env
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
	}))

	v1 := e.Group("/api/v1")

	auth := v1.Group("/auth")

	// auth.POST("/login")
	// auth.POST("/signup")
	// auth.POST("/refresh")
	// auth.POST("/verify")

	if initWebServerError := e.Start(fmt.Sprintf("%s:%d", env.WebServerHost, env.WebServerPort)); initWebServerError != nil {
		fmt.Println(initWebServerError.Error())
		os.Exit(1)
	}
}
