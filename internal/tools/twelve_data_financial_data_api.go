package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nighostchris/everytrack-backend/internal/config"
	"github.com/nighostchris/everytrack-backend/internal/database"
	"go.uber.org/zap"
)

type TickerMeta struct {
	Type             string `json:"type"`
	Symbol           string `json:"symbol"`
	Interval         string `json:"interval"`
	Currency         string `json:"currency"`
	Exchange         string `json:"exchange"`
	MicCode          string `json:"mic_code"`
	ExchangeTimezone string `json:"exchange_timezone"`
}

type TickerPriceHistory struct {
	Low      string `json:"low"`
	Open     string `json:"open"`
	High     string `json:"high"`
	Close    string `json:"close"`
	Volume   string `json:"volume"`
	Datetime string `json:"datetime"`
}

type TickerDetails struct {
	Meta   TickerMeta           `json:"meta"`
	Values []TickerPriceHistory `json:"values"`
	Status string               `json:"status"`
}

type TwelveDataFinancialDataApi struct {
	Logger *zap.Logger
	Db     *pgxpool.Pool
	Env    *config.Config
}

func (tdfda *TwelveDataFinancialDataApi) FetchLatestUSStockPrice() {
	tdfda.Logger.Info("starts")

	// Get US country id from database
	country, getCountryError := database.GetCountryByCode(tdfda.Db, "US")
	if getCountryError != nil {
		tdfda.Logger.Error(fmt.Sprintf("failed to get country details for US from database. %s", getCountryError.Error()))
		return
	}

	// Get all supported US stocks in database
	stocks, getAllStocksError := database.GetAllStocksByCountryId(tdfda.Db, country.Id)
	if getAllStocksError != nil {
		tdfda.Logger.Error(fmt.Sprintf("failed to get US stocks from database. %s", getAllStocksError.Error()))
		return
	}

	// Construct symbols list to embed in API call URL
	rawSymbols := []string{}
	for _, stock := range stocks {
		rawSymbols = append(rawSymbols, stock.Ticker)
	}
	symbol := strings.Join(rawSymbols, ",")
	tdfda.Logger.Info(fmt.Sprintf("going to fetch financial data for supported US stocks - %v", symbol))

	// Try to fetch twelve data financial data api
	rawResponse, fetchError := http.Get(fmt.Sprintf(
		"https://api.twelvedata.com/time_series?symbol=%s&interval=1min&apikey=%s",
		symbol,
		tdfda.Env.TwelveDataApiKey,
	))
	if fetchError != nil {
		tdfda.Logger.Error(fmt.Sprintf("failed to fetch twelve data financial data api. %s", fetchError.Error()))
		return
	}

	// Convert api response into byte array
	response, parseRawResponseError := io.ReadAll(rawResponse.Body)
	if parseRawResponseError != nil {
		tdfda.Logger.Error(fmt.Sprintf("failed to parse twelve data financial data api response. %s", parseRawResponseError.Error()))
		return
	}

	// Convert api response byte array into consumable json
	var data map[string]TickerDetails
	convertJsonError := json.Unmarshal(response, &data)
	if convertJsonError != nil {
		tdfda.Logger.Error(fmt.Sprintf("failed to convert twelve data financial data api response into json. %s", convertJsonError.Error()))
		return
	}

	// Update current price for stocks in database
	for _, stock := range stocks {
		stockDetails := data[stock.Ticker]
		if len(stockDetails.Values) > 0 {
			latestPrice, parseFloatError := strconv.ParseFloat(stockDetails.Values[0].Close, 64)
			if parseFloatError != nil {
				tdfda.Logger.Error(fmt.Sprintf("failed to parse latest price for %s. %s", stock.Ticker, parseFloatError.Error()))
			} else {
				_, updateStockPriceError := database.UpdateStockPrice(tdfda.Db, database.UpdateStockPriceParams{
					Ticker:       stock.Ticker,
					CountryId:    country.Id,
					CurrentPrice: fmt.Sprintf("%.2f", latestPrice),
				})
				if updateStockPriceError != nil {
					tdfda.Logger.Error(fmt.Sprintf("failed to update current price for %s. %s", stock.Ticker, convertJsonError.Error()))
				}
				tdfda.Logger.Info(fmt.Sprintf("updated new price for %s", stock.Ticker))
			}
		} else {
			tdfda.Logger.Info(fmt.Sprintf("no new price for %s to update", stock.Ticker))
		}
	}
}
