package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UpdateExchangeRateParams struct {
	Rate             string `json:"rate"`
	BaseCurrencyId   string `json:"base_currency_id"`
	TargetCurrencyId string `json:"target_currency_id"`
}

func GetAllExchangeRates(db *pgxpool.Pool) ([]ExchangeRate, error) {
	var exchangeRates []ExchangeRate
	query := `SELECT base_currency_id, target_currency_id, rate FROM everytrack_backend.exchange_rate;`
	rows, queryError := db.Query(context.Background(), query)
	if queryError != nil {
		return []ExchangeRate{}, queryError
	}

	defer rows.Close()

	for rows.Next() {
		var exchangeRate ExchangeRate
		scanError := rows.Scan(&exchangeRate.BaseCurrencyId, &exchangeRate.TargetCurrencyId, &exchangeRate.Rate)
		if scanError != nil {
			return []ExchangeRate{}, scanError
		}
		exchangeRates = append(exchangeRates, exchangeRate)
	}

	return exchangeRates, nil
}

func UpdateExchangeRate(db *pgxpool.Pool, params UpdateExchangeRateParams) (bool, error) {
	var existingRowCount int
	query := `SELECT count(*) FROM everytrack_backend.exchange_rate WHERE base_currency_id = $1 AND target_currency_id = $2;`
	queryError := db.QueryRow(context.Background(), query, params.BaseCurrencyId, params.TargetCurrencyId).Scan(&existingRowCount)
	if queryError != nil {
		return false, queryError
	}

	if existingRowCount > 0 {
		updateQuery := `UPDATE everytrack_backend.exchange_rate SET rate = $1 WHERE base_currency_id = $2 AND target_currency_id = $3;`
		_, updateError := db.Exec(context.Background(), updateQuery, params.Rate, params.BaseCurrencyId, params.TargetCurrencyId)

		if updateError != nil {
			return false, updateError
		}
	} else {
		createQuery := `INSERT INTO everytrack_backend.exchange_rate (base_currency_id, target_currency_id, rate) VALUES ($1, $2, $3);`
		_, createError := db.Exec(context.Background(), createQuery, params.BaseCurrencyId, params.TargetCurrencyId, params.Rate)

		if createError != nil {
			return false, createError
		}
	}

	return true, nil
}
