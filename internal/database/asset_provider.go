package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProviderDetail struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	CountryId string `json:"countryId"`
}

func GetAllProvidersByType(db *pgxpool.Pool, providerType string) ([]ProviderDetail, error) {
	var providerDetails []ProviderDetail
	query := `SELECT id, name, icon, country_id FROM everytrack_backend.asset_provider WHERE type = $1;`
	rows, queryError := db.Query(context.Background(), query, providerType)
	if queryError != nil {
		return []ProviderDetail{}, queryError
	}

	defer rows.Close()

	for rows.Next() {
		var providerDetail ProviderDetail
		scanError := rows.Scan(
			&providerDetail.Id,
			&providerDetail.Name,
			&providerDetail.Icon,
			&providerDetail.CountryId,
		)
		if scanError != nil {
			return []ProviderDetail{}, scanError
		}
		providerDetails = append(providerDetails, providerDetail)
	}

	return providerDetails, nil
}
