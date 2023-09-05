package main

import (
	"fmt"
	"os"

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
	db, initPgError := postgres.New(env.DATABASE)

	if initPgError != nil {
		fmt.Println(initPgError.Error())
		os.Exit(1)
	}

	// Initialize web server
}
