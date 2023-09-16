package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountSummary struct {
	Id            string `json:"id"`
	Balance       string `json:"balance"`
	CurrencyId    string `json:"currencyId"`
	AccountTypeId string `json:"accountTypeId"`
}

type CreateNewAccountParams struct {
	ClientId      string `json:"client_id"`
	CurrencyId    string `json:"currency_id"`
	AccountTypeId string `json:"account_type_id"`
}

type UpdateAccountParams struct {
	Balance       string `json:"balance"`
	CurrencyId    string `json:"currency_id"`
	AccountTypeId string `json:"account_type_id"`
}

func GetAllAccountSummaryByType(db *pgxpool.Pool, providerType string, clientId string) ([]AccountSummary, error) {
	accountSummary := []AccountSummary{}
	query := `SELECT a.id, a.balance, apat.id as account_type_id, c.id as currency_id
	FROM everytrack_backend.account AS a
	INNER JOIN everytrack_backend.asset_provider_account_type AS apat ON a.asset_provider_account_type_id = apat.id
	INNER JOIN everytrack_backend.asset_provider AS ap ON apat.asset_provider_id = ap.id
	INNER JOIN everytrack_backend.currency AS c ON c.id = a.currency_id
	WHERE ap.type = $1 AND a.client_id = $2;`
	rows, queryError := db.Query(context.Background(), query, providerType, clientId)
	if queryError != nil {
		return []AccountSummary{}, queryError
	}

	defer rows.Close()

	for rows.Next() {
		var summary AccountSummary
		scanError := rows.Scan(&summary.Id, &summary.Balance, &summary.AccountTypeId, &summary.CurrencyId)
		if scanError != nil {
			return []AccountSummary{}, scanError
		}
		accountSummary = append(accountSummary, summary)
	}

	return accountSummary, nil
}

func CreateNewAccount(db *pgxpool.Pool, params CreateNewAccountParams) (bool, error) {
	query := "INSERT INTO everytrack_backend.account (client_id, asset_provider_account_type_id, currency_id, balance) VALUES ($1, $2, $3, $4);"
	_, createError := db.Exec(context.Background(), query, params.ClientId, params.AccountTypeId, params.CurrencyId, "0")

	if createError != nil {
		return false, createError
	}

	return true, nil
}

func UpdateAccount(db *pgxpool.Pool, params UpdateAccountParams) (bool, error) {
	query := "UPDATE everytrack_backend.account SET balance = $1, currency_id = $2 WHERE asset_provider_account_type_id = $3;"
	_, updateError := db.Exec(context.Background(), query, params.Balance, params.CurrencyId, params.AccountTypeId)

	if updateError != nil {
		return false, updateError
	}

	return true, nil
}

func DeleteAccount(db *pgxpool.Pool, accountId string) (bool, error) {
	query := "DELETE FROM everytrack_backend.account WHERE id = $1;"
	_, deleteError := db.Exec(context.Background(), query, accountId)

	if deleteError != nil {
		return false, deleteError
	}

	return true, nil
}
