package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(connString string) (*pgxpool.Pool, error) {
	// https://pkg.go.dev/github.com/jackc/pgx/v5#ParseConfig
	connConfig, error := pgxpool.ParseConfig(connString)

	if error != nil {
		return nil, error
	}

	// https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool#NewWithConfig
	pool, error := pgxpool.NewWithConfig(context.Background(), connConfig)

	if error != nil {
		return nil, error
	}

	return pool, nil
}
