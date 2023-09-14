package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProviderDetail struct {
	Name            string `json:"name"`
	Icon            string `json:"icon"`
	AccountTypeId   string `json:"account_type_id"`
	AccountTypeName string `json:"account_type_name"`
}

func GetAllProvidersByType(db *pgxpool.Pool, providerType string) ([]ProviderDetail, error) {
	var providerDetails []ProviderDetail
	query := `SELECT ap.name, ap.icon, apat.id as account_type_id, apat.name as account_type_name
	FROM everytrack_backend.asset_provider AS ap
	INNER JOIN everytrack_backend.asset_provider_account_type AS apat ON ap.id = apat.asset_provider_id
	WHERE ap.type = $1;`
	rows, queryError := db.Query(context.Background(), query, providerType)
	if queryError != nil {
		return []ProviderDetail{}, queryError
	}

	defer rows.Close()

	for rows.Next() {
		var providerDetail ProviderDetail
		scanError := rows.Scan(&providerDetail.Name, &providerDetail.Icon, &providerDetail.AccountTypeId, &providerDetail.AccountTypeName)
		if scanError != nil {
			return []ProviderDetail{}, scanError
		}
		providerDetails = append(providerDetails, providerDetail)
	}

	return providerDetails, nil
}
