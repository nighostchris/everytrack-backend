package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nighostchris/everytrack-backend/internal/database"
	"go.uber.org/zap"
)

type GithubCurrencyApi struct {
	Db     *pgxpool.Pool
	Logger *zap.Logger
}

func (gca *GithubCurrencyApi) FetchLatestExchangeRates() {
	gca.Logger.Info("starts")

	// Get all supported currencies in database
	currencies, getCurrenciesError := database.GetAllCurrencies(gca.Db)
	if getCurrenciesError != nil {
		gca.Logger.Error(fmt.Sprintf("failed to get currencies from database. %s", getCurrenciesError.Error()))
		return
	}

	for index, currency := range currencies {
		currencyTicker := strings.ToLower(currency.Ticker)
		interestedCurrencies := []database.Currency{}
		for i, c := range currencies {
			if i != index {
				interestedCurrencies = append(interestedCurrencies, c)
			}
		}
		gca.Logger.Info(fmt.Sprintf("going to fetch github currency api for currency %s", currencyTicker))

		// Try to fetch github currency api
		rawResponse, fetchError := http.Get(fmt.Sprintf("https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1/latest/currencies/%s.json", currencyTicker))
		if fetchError != nil {
			gca.Logger.Error(fmt.Sprintf("failed to fetch github currency api. %s", fetchError.Error()))
			return
		}

		// Convert api response into byte array
		response, parseRawResponseError := io.ReadAll(rawResponse.Body)
		if parseRawResponseError != nil {
			gca.Logger.Error(fmt.Sprintf("failed to parse github currency api response. %s", parseRawResponseError.Error()))
			return
		}

		// Convert api response byte array into consumable json
		var data map[string]interface{}
		convertJsonError := json.Unmarshal(response, &data)
		if convertJsonError != nil {
			gca.Logger.Error(fmt.Sprintf("failed to convert github currency api response into json. %s", convertJsonError.Error()))
			return
		}

		// Insert exchange rate into database
		fetchedRates := data[currencyTicker].(map[string]interface{})
		for _, interestedCurrency := range interestedCurrencies {
			gca.Logger.Info(fmt.Sprintf("updating exchange rate %s:%s in database", currency.Ticker, interestedCurrency.Ticker))
			_, updateError := database.UpdateExchangeRate(gca.Db, database.UpdateExchangeRateParams{
				BaseCurrencyId:   currency.Id,
				TargetCurrencyId: interestedCurrency.Id,
				Rate:             fmt.Sprintf("%.8f", fetchedRates[strings.ToLower(interestedCurrency.Ticker)].(float64)),
			})
			if updateError != nil {
				gca.Logger.Error(fmt.Sprintf("failed to update exchange rate %s:%s in database. %s", currency.Ticker, interestedCurrency.Ticker, updateError.Error()))
			}
		}
	}
}
