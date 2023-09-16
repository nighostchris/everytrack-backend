package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateNewStockHoldingParams struct {
	AccountId string `json:"account_id"`
	StockId   string `json:"stock_id"`
	Unit      string `json:"unit"`
	Cost      string `json:"cost"`
}

func GetAllStockHoldings(db *pgxpool.Pool, clientId string) ([]AccountStock, error) {
	accountStocks := []AccountStock{}
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

func CreateNewStockHolding(db *pgxpool.Pool, params CreateNewStockHoldingParams) (bool, error) {
	query := "INSERT INTO everytrack_backend.account_stock (account_id, stock_id, unit, cost) VALUES ($1, $2, $3, $4);"
	_, createError := db.Exec(context.Background(), query, params.AccountId, params.StockId, params.Unit, params.Cost)

	if createError != nil {
		return false, createError
	}

	return true, nil
}
