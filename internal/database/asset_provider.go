package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BankDetail struct {
	Name            string `json:"name"`
	AccountTypeId   string `json:"account_type_id"`
	AccountTypeName string `json:"account_type_name"`
}

func GetAllBankDetails(db *pgxpool.Pool) ([]BankDetail, error) {
	var bankDetails []BankDetail
	query := `SELECT ap.name, apat.id as account_type_id, apat.name as account_type_name
	FROM everytrack_backend.asset_provider AS ap
	INNER JOIN everytrack_backend.asset_provider_account_type AS apat ON ap.id = apat.asset_provider_id
	WHERE ap.type = 'bank';`
	rows, queryError := db.Query(context.Background(), query)
	if queryError != nil {
		return []BankDetail{}, queryError
	}

	defer rows.Close()

	for rows.Next() {
		var bankDetail BankDetail
		scanError := rows.Scan(&bankDetail.Name, &bankDetail.AccountTypeId, &bankDetail.AccountTypeName)
		if scanError != nil {
			return []BankDetail{}, scanError
		}
		bankDetails = append(bankDetails, bankDetail)
	}

	return bankDetails, nil
}
