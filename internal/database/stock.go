package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
