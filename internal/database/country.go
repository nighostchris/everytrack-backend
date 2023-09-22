package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllCountries(db *pgxpool.Pool) ([]Country, error) {
	countries := []Country{}
	query := `SELECT id, name, code FROM everytrack_backend.country;`
	rows, queryError := db.Query(context.Background(), query)
	if queryError != nil {
		return countries, queryError
	}

	defer rows.Close()

	for rows.Next() {
		var country Country
		scanError := rows.Scan(&country.Id, &country.Name, &country.Code)
		if scanError != nil {
			return countries, scanError
		}
		countries = append(countries, country)
	}

	return countries, nil
}

func GetCountryByCode(db *pgxpool.Pool, code string) (Country, error) {
	var country Country
	query := `SELECT id, name, code FROM everytrack_backend.country WHERE code = $1;`
	queryError := db.QueryRow(context.Background(), query, code).Scan(&country.Id, &country.Name, &country.Code)
	if queryError != nil {
		return country, queryError
	}

	return country, nil
}
