package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllStockHoldings(db *pgxpool.Pool, clientId string) ([]AccountStock, error) {
	var accountStocks []AccountStock
	query := `SELECT accs.id, account_id, stock_id, unit, cost
	FROM everytrack_backend.account_stock AS accs
	INNER JOIN everytrack_backend.account AS a ON a.id = accs.account_id
	WHERE a.client_id = $1;`
	rows, queryError := db.Query(context.Background(), query, clientId)
	if queryError != nil {
		return []AccountStock{}, queryError
	}

	defer rows.Close()

	for rows.Next() {
		var accountStock AccountStock
		scanError := rows.Scan(&accountStock.Id, &accountStock.AccountId, &accountStock.StockId, &accountStock.Unit, &accountStock.Cost)
		if scanError != nil {
			return []AccountStock{}, scanError
		}
		accountStocks = append(accountStocks, accountStock)
	}

	return accountStocks, nil
}
