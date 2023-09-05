package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetClientByEmail(db *pgxpool.Pool, email string) (Client, error) {
	var client Client
	query := "SELECT id, email, password, created_at, updated_at FROM client WHERE email = $1;"
	queryError := db.QueryRow(context.Background(), query, email).Scan(&client.Id, &client.Email, &client.Password, &client.CreatedAt, &client.UpdatedAt)

	if queryError != nil {
		return client, queryError
	}

	return client, nil
}
