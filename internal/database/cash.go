package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateNewCashRecordParams struct {
	Amount     string `json:"amount"`
	ClientId   string `json:"client_id"`
	CurrencyId string `json:"currency_id"`
}

type UpdateCashRecordParams struct {
	Id         string `json:"id"`
	Amount     string `json:"amount"`
	ClientId   string `json:"client_id"`
	CurrencyId string `json:"currency_id"`
}

func GetAllCash(db *pgxpool.Pool, clientId string) ([]Cash, error) {
	cashRecords := []Cash{}
	query := `SELECT id, currency_id, amount FROM everytrack_backend.cash WHERE client_id = $1;`
	rows, queryError := db.Query(context.Background(), query, clientId)
	if queryError != nil {
		return cashRecords, queryError
	}

	defer rows.Close()

	for rows.Next() {
		var cashRecord Cash
		scanError := rows.Scan(
			&cashRecord.Id,
			&cashRecord.CurrencyId,
			&cashRecord.Amount,
		)
		if scanError != nil {
			return cashRecords, scanError
		}
		cashRecords = append(cashRecords, cashRecord)
	}

	return cashRecords, nil
}

func CreateNewCashRecord(db *pgxpool.Pool, params CreateNewCashRecordParams) (bool, error) {
	query := "INSERT INTO everytrack_backend.cash (client_id, currency_id, amount) VALUES ($1, $2, $3);"
	_, createError := db.Exec(
		context.Background(),
		query,
		params.ClientId,
		params.CurrencyId,
		params.Amount,
	)

	if createError != nil {
		return false, createError
	}

	return true, nil
}

func UpdateCashRecord(db *pgxpool.Pool, params UpdateCashRecordParams) (bool, error) {
	query := "UPDATE everytrack_backend.cash SET amount = $1, currency_id = $2 WHERE id = $3 AND client_id = $4;"
	_, updateError := db.Exec(context.Background(), query, params.Amount, params.CurrencyId, params.Id, params.ClientId)

	if updateError != nil {
		return false, updateError
	}

	return true, nil
}

func DeleteCashRecord(db *pgxpool.Pool, cashId string, clientId string) (bool, error) {
	query := "DELETE FROM everytrack_backend.cash WHERE id = $1 AND client_id = $2;"
	_, deleteError := db.Exec(context.Background(), query, cashId, clientId)

	if deleteError != nil {
		return false, deleteError
	}

	return true, nil
}
