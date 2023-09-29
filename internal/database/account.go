package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountSummary struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Balance         string `json:"balance"`
	CurrencyId      string `json:"currencyId"`
	AccountTypeId   string `json:"accountTypeId"`
	AssetProviderId string `json:"assetProviderId"`
}

type CheckExistingAccountParams struct {
	Name            string `json:"name"`
	ClientId        string `json:"client_id"`
	AssetProviderId string `json:"asset_provider_id"`
}

type CreateNewAccountParams struct {
	Name            string `json:"name"`
	ClientId        string `json:"client_id"`
	CurrencyId      string `json:"currency_id"`
	AssetProviderId string `json:"asset_provider_id"`
}

type UpdateAccountParams struct {
	Balance       string `json:"balance"`
	CurrencyId    string `json:"currency_id"`
	AccountTypeId string `json:"account_type_id"`
}

func GetAllAccountSummaryByType(db *pgxpool.Pool, providerType string, clientId string) ([]AccountSummary, error) {
	accountSummary := []AccountSummary{}
	query := `SELECT a.id, a.balance, apat.name, ap.id as asset_provider_id, apat.id as account_type_id, c.id as currency_id
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
		scanError := rows.Scan(
			&summary.Id,
			&summary.Balance,
			&summary.Name,
			&summary.AssetProviderId,
			&summary.AccountTypeId,
			&summary.CurrencyId,
		)
		if scanError != nil {
			return []AccountSummary{}, scanError
		}
		accountSummary = append(accountSummary, summary)
	}

	return accountSummary, nil
}

func GetAccountBalance(db *pgxpool.Pool, accountId string) (string, error) {
	var balance string
	getBalanceQuery := `SELECT balance FROM everytrack_backend.account WHERE id = $1;`
	getBalanceError := db.QueryRow(context.Background(), getBalanceQuery, accountId).Scan(&balance)
	if getBalanceError != nil {
		return balance, getBalanceError
	}

	return balance, nil
}

func CheckExistingAccount(db *pgxpool.Pool, params CheckExistingAccountParams) (bool, error) {
	var existingRowCount int
	query := `SELECT count(*)
	FROM everytrack_backend.asset_provider_account_type as apat
	INNER JOIN everytrack_backend.account as acc ON acc.asset_provider_account_type_id = apat.id
	WHERE apat.name = $1 AND acc.client_id = $2 AND apat.asset_provider_id = $3;`
	queryError := db.QueryRow(context.Background(), query, params.Name, params.ClientId, params.AssetProviderId).Scan(&existingRowCount)
	if queryError != nil {
		return false, queryError
	}

	if existingRowCount > 0 {
		return true, nil
	}

	return false, nil
}

func CreateNewAccount(db *pgxpool.Pool, params CreateNewAccountParams) (bool, error) {
	// Create new asset provider account type first
	var assetProviderAccountTypeId string
	insertAccountTypeQuery := "INSERT INTO everytrack_backend.asset_provider_account_type (asset_provider_id, name) VALUES ($1, $2) RETURNING id;"
	insertAccountTypeError := db.QueryRow(context.Background(), insertAccountTypeQuery, params.AssetProviderId, params.Name).Scan(&assetProviderAccountTypeId)

	if insertAccountTypeError != nil {
		return false, insertAccountTypeError
	}

	// Create new account using the newly created account type id above
	insertAccountQuery := "INSERT INTO everytrack_backend.account (client_id, asset_provider_account_type_id, currency_id, balance) VALUES ($1, $2, $3, $4);"
	_, insertAccountError := db.Exec(context.Background(), insertAccountQuery, params.ClientId, assetProviderAccountTypeId, params.CurrencyId, "0")

	if insertAccountError != nil {
		return false, insertAccountError
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

func UpdateAccountBalance(db *pgxpool.Pool, balance string, accountId string) (bool, error) {
	query := "UPDATE everytrack_backend.account SET balance = $1 WHERE id = $2;"
	_, updateError := db.Exec(context.Background(), query, balance, accountId)

	if updateError != nil {
		return false, updateError
	}

	return true, nil
}

func DeleteAccount(db *pgxpool.Pool, accountId string) (bool, error) {
	// Retrieve relevant asset_provider_account_type from database
	var accountTypeId string
	getAccountTypeIdQuery := `SELECT asset_provider_account_type_id FROM everytrack_backend.account WHERE id = $1;`
	getAccountTypeIdError := db.QueryRow(context.Background(), getAccountTypeIdQuery, accountId).Scan(&accountTypeId)
	if getAccountTypeIdError != nil {
		return false, getAccountTypeIdError
	}

	// Delete relevant account_stock records from database
	deleteStockHoldingsQuery := "DELETE FROM everytrack_backend.account_stock WHERE account_id = $1;"
	_, deleteStockHoldingsError := db.Exec(context.Background(), deleteStockHoldingsQuery, accountId)

	if deleteStockHoldingsError != nil {
		return false, deleteStockHoldingsError
	}

	// Delete account from database
	query := "DELETE FROM everytrack_backend.account WHERE id = $1;"
	_, deleteError := db.Exec(context.Background(), query, accountId)

	if deleteError != nil {
		return false, deleteError
	}

	// Delete relevant asset_provider_account_type from database since it's 1-1 relationship
	deleteAccountTypeQuery := "DELETE FROM everytrack_backend.asset_provider_account_type WHERE id = $1;"
	_, deleteAccountTypeError := db.Exec(context.Background(), deleteAccountTypeQuery, accountTypeId)

	if deleteAccountTypeError != nil {
		return false, deleteAccountTypeError
	}

	return true, nil
}
