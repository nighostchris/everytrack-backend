package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllCurrencies(db *pgxpool.Pool) ([]Currency, error) {
	var currencies []Currency
	query := `SELECT id, ticker, symbol FROM everytrack_backend.currency;`
	rows, queryError := db.Query(context.Background(), query)
	if queryError != nil {
		return []Currency{}, queryError
	}

	defer rows.Close()

	for rows.Next() {
		var currency Currency
		scanError := rows.Scan(&currency.Id, &currency.Ticker, &currency.Symbol)
		if scanError != nil {
			return []Currency{}, scanError
		}
		currencies = append(currencies, currency)
	}

	return currencies, nil
}

func GetDefaultCurrency(db *pgxpool.Pool) (string, error) {
	var id string
	query := `SELECT id FROM everytrack_backend.currency WHERE ticker = $1;`
	queryError := db.QueryRow(context.Background(), query, "GBP").Scan(&id)
	if queryError != nil {
		return id, queryError
	}

	return id, nil
}
