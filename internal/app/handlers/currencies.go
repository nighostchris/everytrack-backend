package handlers

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nighostchris/everytrack-backend/internal/database"
	"go.uber.org/zap"
)

type CurrenciesHandler struct {
	Db     *pgxpool.Pool
	Logger *zap.Logger
}

func (ch *CurrenciesHandler) GetAllCurrencies(c echo.Context) error {
	ch.Logger.Info("starts")

	// Get all bank details from database
	currencies, getCurrenciesError := database.GetAllCurrencies(ch.Db)

	if getCurrenciesError != nil {
		ch.Logger.Error(fmt.Sprintf("failed to get currencies from database. %s", getCurrenciesError.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Internal server error."})
	}
	ch.Logger.Debug("got currencies from database")

	return c.JSON(http.StatusOK, map[string]interface{}{"success": true, "data": currencies})
}
