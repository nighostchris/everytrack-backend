package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateNewFuturePaymentParams struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Income      bool      `json:"income"`
	Amount      string    `json:"amount"`
	Remarks     *string   `json:"remarks"`
	Rolling     bool      `json:"rolling"`
	ClientId    string    `json:"client_id"`
	AccountId   string    `json:"account_id"`
	CurrencyId  string    `json:"currency_id"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

func GetAllFuturePayments(db *pgxpool.Pool, clientId string) ([]FuturePayment, error) {
	futurePayments := []FuturePayment{}
	query := `SELECT id, account_id, currency_id, name, amount, income, rolling, remarks, scheduled_at FROM everytrack_backend.future_payment WHERE client_id = $1;`
	rows, queryError := db.Query(context.Background(), query, clientId)
	if queryError != nil {
		return futurePayments, queryError
	}

	defer rows.Close()

	for rows.Next() {
		var futurePayment FuturePayment
		scanError := rows.Scan(
			&futurePayment.Id,
			&futurePayment.AccountId,
			&futurePayment.CurrencyId,
			&futurePayment.Name,
			&futurePayment.Amount,
			&futurePayment.Income,
			&futurePayment.Rolling,
			&futurePayment.Remarks,
			&futurePayment.ScheduledAt,
		)
		if scanError != nil {
			return futurePayments, scanError
		}
		futurePayments = append(futurePayments, futurePayment)
	}

	return futurePayments, nil
}

func CreateNewFuturePayment(db *pgxpool.Pool, params CreateNewFuturePaymentParams) (bool, error) {
	query := "INSERT INTO everytrack_backend.future_payment (client_id, account_id, currency_id, name, amount, income, rolling, remarks, scheduled_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);"
	_, createError := db.Exec(
		context.Background(),
		query,
		params.ClientId,
		params.AccountId,
		params.CurrencyId,
		params.Name,
		params.Amount,
		params.Income,
		params.Rolling,
		params.Remarks,
		params.ScheduledAt,
	)

	if createError != nil {
		return false, createError
	}

	return true, nil
}

func DeleteFuturePayment(db *pgxpool.Pool, futurePaymentId string, clientId string) (bool, error) {
	query := "DELETE FROM everytrack_backend.future_payment WHERE id = $1 AND client_id = $2;"
	_, deleteError := db.Exec(context.Background(), query, futurePaymentId, clientId)

	if deleteError != nil {
		return false, deleteError
	}

	return true, nil
}
