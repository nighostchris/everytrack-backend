package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateNewClientParams struct {
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	CurrencyId string `json:"currency_id"`
}

type UpdateClientSettingsParams struct {
	Username   string `json:"username"`
	ClientId   string `json:"client_id"`
	CurrencyId string `json:"currency_id"`
}

func CreateNewClient(db *pgxpool.Pool, params CreateNewClientParams) (string, error) {
	var id string
	query := "INSERT INTO everytrack_backend.client (email, username, password, currency_id) VALUES ($1, $2, $3, $4) RETURNING id;"
	queryError := db.QueryRow(context.Background(), query, params.Email, params.Username, params.Password, params.CurrencyId).Scan(&id)

	if queryError != nil {
		return id, queryError
	}

	return id, nil
}

func GetClientByEmail(db *pgxpool.Pool, email string) (Client, error) {
	var client Client
	query := "SELECT id, email, password, currency_id, created_at, updated_at FROM everytrack_backend.client WHERE email = $1;"
	queryError := db.QueryRow(context.Background(), query, email).Scan(&client.Id, &client.Email, &client.Password, &client.CurrencyId, &client.CreatedAt, &client.UpdatedAt)

	if queryError != nil {
		return client, queryError
	}

	return client, nil
}

func GetClientById(db *pgxpool.Pool, id string) (Client, error) {
	var client Client
	query := "SELECT id, email, username, password, currency_id, created_at, updated_at FROM everytrack_backend.client WHERE id = $1;"
	queryError := db.QueryRow(context.Background(), query, id).Scan(&client.Id, &client.Email, &client.Username, &client.Password, &client.CurrencyId, &client.CreatedAt, &client.UpdatedAt)

	if queryError != nil {
		return client, queryError
	}

	return client, nil
}

func UpdateClientSettings(db *pgxpool.Pool, params UpdateClientSettingsParams) (bool, error) {
	query := "UPDATE everytrack_backend.client SET username = $1, currency_id = $2 WHERE id = $3;"
	_, updateError := db.Exec(context.Background(), query, params.Username, params.CurrencyId, params.ClientId)

	if updateError != nil {
		return false, updateError
	}

	return true, nil
}
