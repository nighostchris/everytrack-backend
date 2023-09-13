package handlers

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nighostchris/everytrack-backend/internal/database"
	"go.uber.org/zap"
)

type ExchangeRatesHandler struct {
	Db     *pgxpool.Pool
	Logger *zap.Logger
}

type ExchangeRateData struct {
	Rate             string `json:"rate"`
	BaseCurrencyId   string `json:"baseCurrencyId"`
	TargetCurrencyId string `json:"targetCurrencyId"`
}

func (erh *ExchangeRatesHandler) GetAllExchangeRates(c echo.Context) error {
	erh.Logger.Info("starts")

	// Get all exchange rates from database
	exchangeRates, getExchangeRatesError := database.GetAllExchangeRates(erh.Db)

	if getExchangeRatesError != nil {
		erh.Logger.Error(fmt.Sprintf("failed to get exchange rates from database. %s", getExchangeRatesError.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Internal server error."})
	}
	erh.Logger.Debug("got exchange rates from database")

	// Construct response object
	responseData := []ExchangeRateData{}
	for _, exchangeRate := range exchangeRates {
		responseData = append(responseData, ExchangeRateData{
			Rate:             exchangeRate.Rate,
			BaseCurrencyId:   exchangeRate.BaseCurrencyId,
			TargetCurrencyId: exchangeRate.TargetCurrencyId,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"success": true, "data": responseData})
}
