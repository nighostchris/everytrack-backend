package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateNewExpenseParams struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	Amount     string    `json:"amount"`
	Remarks    *string   `json:"remarks"`
	Category   string    `json:"category"`
	ClientId   string    `json:"client_id"`
	AccountId  *string   `json:"account_id"`
	CurrencyId string    `json:"currency_id"`
	ExecutedAt time.Time `json:"executed_at"`
}

type ExpenseAndAccountBalance struct {
	Amount    string `json:"amount"`
	Balance   string `json:"balance"`
	AccountId string `json:"account_id"`
}

func GetAllExpenses(db *pgxpool.Pool, clientId string) ([]Expense, error) {
	expenses := []Expense{}
	query := `SELECT id, name, account_id, currency_id, category, amount, remarks, executed_at FROM everytrack_backend.expense WHERE client_id = $1;`
	rows, queryError := db.Query(context.Background(), query, clientId)
	if queryError != nil {
		return expenses, queryError
	}

	defer rows.Close()

	for rows.Next() {
		var expense Expense
		scanError := rows.Scan(
			&expense.Id,
			&expense.Name,
			&expense.AccountId,
			&expense.CurrencyId,
			&expense.Category,
			&expense.Amount,
			&expense.Remarks,
			&expense.ExecutedAt,
		)
		if scanError != nil {
			return expenses, scanError
		}
		expenses = append(expenses, expense)
	}

	return expenses, nil
}

func GetExpenseAndAccountBalanceById(db *pgxpool.Pool, expenseId string) (ExpenseAndAccountBalance, error) {
	result := ExpenseAndAccountBalance{}
	query := `SELECT amount, balance, account_id
	FROM everytrack_backend.expense as e
	INNER JOIN everytrack_backend.account as a ON a.id = e.account_id
	WHERE e.id = $1;`
	queryError := db.QueryRow(context.Background(), query, expenseId).Scan(&result.Amount, &result.Balance, &result.AccountId)
	if queryError != nil {
		return result, queryError
	}

	return result, nil
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

func DeleteExpense(db *pgxpool.Pool, expenseId string, clientId string) (bool, error) {
	query := "DELETE FROM everytrack_backend.expense WHERE id = $1 AND client_id = $2;"
	_, deleteError := db.Exec(context.Background(), query, expenseId, clientId)

	if deleteError != nil {
		return false, deleteError
	}

	return true, nil
}
