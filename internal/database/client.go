package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateNewClientParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateNewClient(db *pgxpool.Pool, params CreateNewClientParams) (string, error) {
	var id string
	query := "INSERT INTO client (email, password) VALUES ($1, $2) RETURNING id;"
	queryError := db.QueryRow(context.Background(), query, params.Email, params.Password).Scan(&id)

	if queryError != nil {
		return id, queryError
	}

	return id, nil
}

func GetClientByEmail(db *pgxpool.Pool, email string) (Client, error) {
	var client Client
	query := "SELECT id, email, password, created_at, updated_at FROM client WHERE email = $1;"
	queryError := db.QueryRow(context.Background(), query, email).Scan(&client.Id, &client.Email, &client.Password, &client.CreatedAt, &client.UpdatedAt)

	if queryError != nil {
		return client, queryError
	}

	return client, nil
}
