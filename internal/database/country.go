package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetCountryByCode(db *pgxpool.Pool, code string) (Country, error) {
	var country Country
	query := `SELECT id, name, code FROM everytrack_backend.country WHERE code = $1;`
	queryError := db.QueryRow(context.Background(), query, code).Scan(&country.Id, &country.Name, &country.Code)
	if queryError != nil {
		return country, queryError
	}

	return country, nil
}
