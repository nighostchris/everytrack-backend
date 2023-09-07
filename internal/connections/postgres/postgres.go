package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(connString string) *pgxpool.Pool {
	logTemplate := "{\"level\":\"%s\",\"timestamp\":\"%s\",\"function\":\"github.com/nighostchris/everytrack-backend/internal/postgres.New\",\"message\":\"%s\"}\n"
	fmt.Printf(logTemplate, "info", time.Now().Format(time.RFC3339Nano), "initializing database connection")

	// https://pkg.go.dev/github.com/jackc/pgx/v5#ParseConfig
	connConfig, parseConfigError := pgxpool.ParseConfig(connString)
	if parseConfigError != nil {
		fmt.Printf(logTemplate, "error", time.Now().Format(time.RFC3339Nano), parseConfigError.Error())
		os.Exit(1)
	}

	// https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool#NewWithConfig
	pool, initPoolError := pgxpool.NewWithConfig(context.Background(), connConfig)
	if initPoolError != nil {
		fmt.Printf(logTemplate, "error", time.Now().Format(time.RFC3339Nano), initPoolError.Error())
		os.Exit(1)
	}
	fmt.Printf(logTemplate, "info", time.Now().Format(time.RFC3339Nano), "initialized database connection")

	return pool
}
