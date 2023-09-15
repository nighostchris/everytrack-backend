package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nighostchris/everytrack-backend/internal/database"
	"go.uber.org/zap"
)

type StocksHandler struct {
	Db     *pgxpool.Pool
	Logger *zap.Logger
}

type StockRecord struct {
	Id           string `json:"id"`
	CountryId    string `json:"countryId"`
	CurrencyId   string `json:"currencyId"`
	Name         string `json:"name"`
	Ticker       string `json:"ticker"`
	CurrentPrice string `json:"currentPrice"`
}

type StockHolding struct {
	Id      string `json:"id"`
	StockId string `json:"stockId"`
	Unit    string `json:"unit"`
	Cost    string `json:"cost"`
}

type StockHoldingRecord struct {
	AccountId string         `json:"accountId"`
	Holdings  []StockHolding `json:"holdings"`
}

type CreateNewStockHoldingRequestBody struct {
	Unit      string `json:"unit"`
	Cost      string `json:"cost"`
	StockId   string `json:"stockId"`
	AccountId string `json:"accountId"`
}

func (sh *StocksHandler) GetAllStocks(c echo.Context) error {
	requestId := zap.String("requestId", c.Get("requestId").(string))
	sh.Logger.Info("starts", requestId)

	// Get all stocks from database
	stocks, getStocksError := database.GetAllStocks(sh.Db)
	if getStocksError != nil {
		sh.Logger.Error(
			fmt.Sprintf("failed to get stocks from database. %s", getStocksError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	sh.Logger.Debug("got stocks from database", requestId)

	// Construct the response object
	stocksJson, marshalError := json.Marshal(stocks)
	if marshalError != nil {
		sh.Logger.Error(
			fmt.Sprintf("failed to marshal stocks into raw json. %s", marshalError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	var stockRecords []StockRecord
	unmarshalError := json.Unmarshal([]byte(stocksJson), &stockRecords)
	if unmarshalError != nil {
		sh.Logger.Error(
			fmt.Sprintf("failed to unmarshal stocksJson into response object. %s", unmarshalError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	sh.Logger.Debug(fmt.Sprintf("constructed response object - %#v", stockRecords), requestId)

	return c.JSON(http.StatusOK, map[string]interface{}{"success": true, "data": stockRecords})
}

func (sh *StocksHandler) GetAllStockHoldings(c echo.Context) error {
	clientId := c.Get("uid").(string)
	requestId := zap.String("requestId", c.Get("requestId").(string))
	sh.Logger.Info("starts", requestId)

	// Get all stock holdings of user from database
	accountStocks, getAccountStocksError := database.GetAllStockHoldings(sh.Db, clientId)
	if getAccountStocksError != nil {
		sh.Logger.Error(
			fmt.Sprintf("failed to get all stock holdings from database. %s", getAccountStocksError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	sh.Logger.Debug("got stock holdings from database", requestId)

	// Construct the response object
	accountStockHoldingsMap := make(map[string][]StockHolding)
	for _, accountStock := range accountStocks {
		stockHolding := StockHolding{Id: accountStock.Id, StockId: accountStock.StockId, Unit: accountStock.Unit, Cost: accountStock.Cost}
		accountStockHoldingsMap[accountStock.AccountId] = append(accountStockHoldingsMap[accountStock.AccountId], stockHolding)
	}
	responseData := []StockHoldingRecord{}
	for accountId, stockHoldings := range accountStockHoldingsMap {
		responseData = append(responseData, StockHoldingRecord{
			AccountId: accountId,
			Holdings:  stockHoldings,
		})
	}
	sh.Logger.Debug(fmt.Sprintf("constructed response object - %#v", responseData), requestId)

	return c.JSON(http.StatusOK, map[string]interface{}{"success": true, "data": responseData})
}

func (sh *StocksHandler) CreateNewStockHolding(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"success": true, "data": ""})
}
