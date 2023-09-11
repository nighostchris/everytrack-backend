package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BankAccount struct {
	Balance       string `json:"balance"`
	Currency      string `json:"currency"`
	AccountTypeId string `json:"account_type_id"`
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

func GetAllBankAccounts(db *pgxpool.Pool, clientId string) ([]BankAccount, error) {
	var bankAccounts []BankAccount
	query := `SELECT a.balance, apat.id as account_type_id, c.symbol as currency
	FROM everytrack_backend.account AS a
	INNER JOIN everytrack_backend.asset_provider_account_type AS apat ON a.asset_provider_account_type_id = apat.id
	INNER JOIN everytrack_backend.asset_provider AS ap ON apat.asset_provider_id = ap.id
	INNER JOIN everytrack_backend.currency AS c ON c.id = a.currency_id
	WHERE ap.type = 'bank' AND a.client_id = $1;`
	rows, queryError := db.Query(context.Background(), query, clientId)
	if queryError != nil {
		return []BankAccount{}, queryError
	}

	defer rows.Close()

	for rows.Next() {
		var bankAccount BankAccount
		scanError := rows.Scan(&bankAccount.Balance, &bankAccount.AccountTypeId, &bankAccount.Currency)
		if scanError != nil {
			return []BankAccount{}, scanError
		}
		bankAccounts = append(bankAccounts, bankAccount)
	}

	return bankAccounts, nil
}

func CreateNewAccount(db *pgxpool.Pool, params CreateNewAccountParams) (bool, error) {
	query := "INSERT INTO everytrack_backend.account (client_id, asset_provider_account_type_id, currency_id, balance) VALUES ($1, $2, $3, $4);"
	_, createError := db.Exec(context.Background(), query, params.ClientId, params.AccountTypeId, params.CurrencyId, 0)

	if createError != nil {
		return false, createError
	}

	return true, nil
}

func UpdateAccount(db *pgxpool.Pool, params UpdateAccountParams) (bool, error) {
	query := "UPDATE everytrack_backend.account SET balance = $1, asset_provider_account_type_id = $2, currency_id = $3;"
	_, updateError := db.Exec(context.Background(), query, params.Balance, params.AccountTypeId, params.CurrencyId)

	if updateError != nil {
		return false, updateError
	}

	return true, nil
}
