package database

import "time"

type Account struct {
	Id                         string    `json:"id"`
	ClientId                   string    `json:"client_id"`
	AssetProviderAccountTypeId string    `json:"asset_provider_account_type_id"`
	Balance                    string    `json:"balance"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

type AssetProvider struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AssetProviderAccountType struct {
	Id              string    `json:"id"`
	AssetProviderId string    `json:"asset_provider_id"`
	Name            string    `json:"name"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Client struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
