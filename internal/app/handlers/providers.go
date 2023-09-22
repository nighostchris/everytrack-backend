package handlers

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nighostchris/everytrack-backend/internal/database"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

type ProvidersHandler struct {
	Db     *pgxpool.Pool
	Logger *zap.Logger
}

type ProviderDetailsRecord struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	CountryId string `json:"countryid"`
}

var ProviderTypes = []string{"savings", "broker", "credit"}

func (ph *ProvidersHandler) GetAllProvidersByType(c echo.Context) error {
	requestId := zap.String("requestId", c.Get("requestId").(string))
	ph.Logger.Info("starts", requestId)

	providerType := c.QueryParam("type")
	if len(providerType) == 0 {
		ph.Logger.Error("undefined provider type", requestId)
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Undefined provider type."},
		)
	}
	if !slices.Contains(ProviderTypes, providerType) {
		ph.Logger.Error(fmt.Sprintf("invalid provider type %s", providerType), requestId)
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Invalid provider type."},
		)
	}
	ph.Logger.Info(fmt.Sprintf("going to get all providers by type %s", providerType), requestId)

	// Get all provider details from database
	providerDetails, getProviderDetailsError := database.GetAllProvidersByType(ph.Db, providerType)
	if getProviderDetailsError != nil {
		ph.Logger.Error(
			fmt.Sprintf("failed to get provider details from database. %s", getProviderDetailsError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	ph.Logger.Debug(fmt.Sprintf("got providers from database - %#v", providerDetails), requestId)

	return c.JSON(
		http.StatusOK,
		LooseJson{"success": true, "data": providerDetails},
	)
}
