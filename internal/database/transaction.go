package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateNewTransactionParams struct {
	Name       string    `json:"name"`
	Income     bool      `json:"income"`
	Amount     string    `json:"amount"`
	Remarks    *string   `json:"remarks"`
	Category   string    `json:"category"`
	ClientId   string    `json:"client_id"`
	AccountId  string    `json:"account_id"`
	CurrencyId string    `json:"currency_id"`
	ExecutedAt time.Time `json:"executed_at"`
}

type TransactionAndAccountBalance struct {
	Income    bool   `json:"income"`
	Amount    string `json:"amount"`
	Balance   string `json:"balance"`
	AccountId string `json:"account_id"`
}

func GetAllTransactions(db *pgxpool.Pool, clientId string) ([]Transaction, error) {
	transactions := []Transaction{}
	query := `SELECT id, name, income, account_id, currency_id, category, amount, remarks, executed_at FROM everytrack_backend.transaction WHERE client_id = $1;`
	rows, queryError := db.Query(context.Background(), query, clientId)
	if queryError != nil {
		return transactions, queryError
	}

	defer rows.Close()

	for rows.Next() {
		var transaction Transaction
		scanError := rows.Scan(
			&transaction.Id,
			&transaction.Name,
			&transaction.Income,
			&transaction.AccountId,
			&transaction.CurrencyId,
			&transaction.Category,
			&transaction.Amount,
			&transaction.Remarks,
			&transaction.ExecutedAt,
		)
		if scanError != nil {
			return transactions, scanError
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func GetTransactionAndAccountBalanceById(db *pgxpool.Pool, transactionId string) (TransactionAndAccountBalance, error) {
	result := TransactionAndAccountBalance{}
	query := `SELECT amount, income, balance, account_id
	FROM everytrack_backend.transaction as t
	INNER JOIN everytrack_backend.account as a ON a.id = t.account_id
	WHERE t.id = $1;`
	queryError := db.QueryRow(context.Background(), query, transactionId).Scan(&result.Amount, &result.Income, &result.Balance, &result.AccountId)
	if queryError != nil {
		return result, queryError
	}

	return result, nil
}

func CreateNewTransaction(db *pgxpool.Pool, params CreateNewTransactionParams) (bool, error) {
	query := "INSERT INTO everytrack_backend.transaction (client_id, account_id, currency_id, name, category, amount, income, remarks, executed_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);"
	_, createError := db.Exec(
		context.Background(),
		query,
		params.ClientId,
		params.AccountId,
		params.CurrencyId,
		params.Name,
		params.Category,
		params.Amount,
		params.Income,
		params.Remarks,
		params.ExecutedAt,
	)

	if createError != nil {
		return false, createError
	}

	return true, nil
}

func DeleteTransaction(db *pgxpool.Pool, transactionId string, clientId string) (bool, error) {
	query := "DELETE FROM everytrack_backend.transaction WHERE id = $1 AND client_id = $2;"
	_, deleteError := db.Exec(context.Background(), query, transactionId, clientId)

	if deleteError != nil {
		return false, deleteError
	}

	return true, nil
}
