package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateNewClientParams struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateNewClient(db *pgxpool.Pool, params CreateNewClientParams) (string, error) {
	var id string
	query := "INSERT INTO everytrack_backend.client (email, username, password) VALUES ($1, $2, $3) RETURNING id;"
	queryError := db.QueryRow(context.Background(), query, params.Email, params.Username, params.Password).Scan(&id)

	if queryError != nil {
		return id, queryError
	}

	return id, nil
}

func GetClientByEmail(db *pgxpool.Pool, email string) (Client, error) {
	var client Client
	query := "SELECT id, email, password, created_at, updated_at FROM everytrack_backend.client WHERE email = $1;"
	queryError := db.QueryRow(context.Background(), query, email).Scan(&client.Id, &client.Email, &client.Password, &client.CreatedAt, &client.UpdatedAt)

	if queryError != nil {
		return client, queryError
	}

	return client, nil
}
