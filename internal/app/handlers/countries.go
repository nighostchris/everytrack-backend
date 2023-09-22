package handlers

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nighostchris/everytrack-backend/internal/database"
	"go.uber.org/zap"
)

type CountriesHandler struct {
	Db     *pgxpool.Pool
	Logger *zap.Logger
}

func (ch *CountriesHandler) GetAllCountries(c echo.Context) error {
	ch.Logger.Info("starts")

	// Get all countries from database
	countries, getCountriesError := database.GetAllCountries(ch.Db)

	if getCountriesError != nil {
		ch.Logger.Error(fmt.Sprintf("failed to get countries from database. %s", getCountriesError.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Internal server error."})
	}
	ch.Logger.Debug(fmt.Sprintf("got countries from database - %#v", countries))

	return c.JSON(http.StatusOK, map[string]interface{}{"success": true, "data": countries})
}
