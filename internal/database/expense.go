package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateNewExpenseParams struct {
	Name       string    `json:"name"`
	Amount     string    `json:"amount"`
	Remarks    *string   `json:"remarks"`
	Category   string    `json:"category"`
	ClientId   string    `json:"client_id"`
	AccountId  *string   `json:"account_id"`
	CurrencyId string    `json:"currency_id"`
	ExecutedAt time.Time `json:"executed_at"`
}

func GetAllExpenses(db *pgxpool.Pool, clientId string) ([]Expense, error) {
	expenses := []Expense{}
	query := `SELECT name, currency_id, category, amount, remarks, executed_at FROM everytrack_backend.expense WHERE client_id = $1;`
	rows, queryError := db.Query(context.Background(), query, clientId)
	if queryError != nil {
		return expenses, queryError
	}

	defer rows.Close()

	for rows.Next() {
		var expense Expense
		scanError := rows.Scan(&expense.Name, &expense.CurrencyId, &expense.Category, &expense.Amount, &expense.Remarks, &expense.ExecutedAt)
		if scanError != nil {
			return expenses, scanError
		}
		expenses = append(expenses, expense)
	}

	return expenses, nil
}

func CreateNewExpense(db *pgxpool.Pool, params CreateNewExpenseParams) (bool, error) {
	query := "INSERT INTO everytrack_backend.expense (client_id, account_id, currency_id, name, category, amount, remarks, executed_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);"
	_, createError := db.Exec(
		context.Background(),
		query,
		params.ClientId,
		params.AccountId,
		params.CurrencyId,
		params.Name,
		params.Category,
		params.Amount,
		params.Remarks,
		params.ExecutedAt,
	)

	if createError != nil {
		return false, createError
	}

	return true, nil
}
