package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateNewExpenseParams struct {
	Amount     string    `json:"amount"`
	Remarks    *string   `json:"remarks"`
	Category   string    `json:"category"`
	ClientId   string    `json:"client_id"`
	AccountId  *string   `json:"account_id"`
	CurrencyId string    `json:"currency_id"`
	ExecutedAt time.Time `json:"executed_at"`
}

func CreateNewExpense(db *pgxpool.Pool, params CreateNewExpenseParams) (bool, error) {
	query := "INSERT INTO everytrack_backend.expense (client_id, account_id, currency_id, category, amount, remarks, executed_at) VALUES ($1, $2, $3, $4, $5, $6, $7);"
	_, createError := db.Exec(
		context.Background(),
		query,
		params.ClientId,
		params.AccountId,
		params.CurrencyId,
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
