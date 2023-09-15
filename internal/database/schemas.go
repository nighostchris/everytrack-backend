package database

import "time"

type Account struct {
	Id                         string    `json:"id"`
	ClientId                   string    `json:"client_id"`
	AssetProviderAccountTypeId string    `json:"asset_provider_account_type_id"`
	CurrencyId                 string    `json:"currency_id"`
	Balance                    string    `json:"balance"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

type AccountStock struct {
	Id        string `json:"id"`
	AccountId string `json:"account_id"`
	StockId   string `json:"stock_id"`
	Unit      string `json:"unit"`
	Cost      string `json:"cost"`
}

type AssetProvider struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	Type      string    `json:"type"`
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
	Id         string    `json:"id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	CurrencyId string    `json:"currency_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Currency struct {
	Id     string `json:"id"`
	Ticker string `json:"ticker"`
	Symbol string `json:"symbol"`
}

type ExchangeRate struct {
	Id               string `json:"id"`
	BaseCurrencyId   string `json:"base_currency_id"`
	TargetCurrencyId string `json:"target_currency_id"`
	Rate             string `json:"rate"`
}

type Stock struct {
	Id           string `json:"id"`
	CountryId    string `json:"country_id"`
	CurrencyId   string `json:"currency_id"`
	Name         string `json:"name"`
	Ticker       string `json:"ticker"`
	CurrentPrice string `json:"current_price"`
}
