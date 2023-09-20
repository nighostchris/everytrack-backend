package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UpdateStockPriceParams struct {
	Ticker       string `json:"ticker"`
	CountryId    string `json:"country_id"`
	CurrentPrice string `json:"current_price"`
}

func GetAllStocks(db *pgxpool.Pool) ([]Stock, error) {
	var stocks []Stock
	query := `SELECT id, country_id, currency_id, name, ticker, current_price FROM everytrack_backend.stock;`
	rows, queryError := db.Query(context.Background(), query)
	if queryError != nil {
		return []Stock{}, queryError
	}

	defer rows.Close()

	for rows.Next() {
		var stock Stock
		scanError := rows.Scan(&stock.Id, &stock.CountryId, &stock.CurrencyId, &stock.Name, &stock.Ticker, &stock.CurrentPrice)
		if scanError != nil {
			return []Stock{}, scanError
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

func GetAllStocksByCountryId(db *pgxpool.Pool, id string) ([]Stock, error) {
	stocks := []Stock{}
	query := `SELECT s.id, country_id, currency_id, s.name, ticker, current_price
	FROM everytrack_backend.stock as s
	INNER JOIN everytrack_backend.country as c ON s.country_id = c.id
	WHERE s.country_id = $1;`
	rows, queryError := db.Query(context.Background(), query, id)
	if queryError != nil {
		return []Stock{}, queryError
	}

	defer rows.Close()

	for rows.Next() {
		var stock Stock
		scanError := rows.Scan(&stock.Id, &stock.CountryId, &stock.CurrencyId, &stock.Name, &stock.Ticker, &stock.CurrentPrice)
		if scanError != nil {
			return []Stock{}, scanError
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

func UpdateStockPrice(db *pgxpool.Pool, params UpdateStockPriceParams) (bool, error) {
	query := "UPDATE everytrack_backend.stock SET current_price = $1 WHERE ticker = $2 AND country_id = $3;"
	_, updateError := db.Exec(context.Background(), query, params.CurrentPrice, params.Ticker, params.CountryId)

	if updateError != nil {
		return false, updateError
	}

	return true, nil
}
