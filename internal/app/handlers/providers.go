package handlers

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nighostchris/everytrack-backend/internal/database"
	"go.uber.org/zap"
)

type ProvidersHandler struct {
	Db     *pgxpool.Pool
	Logger *zap.Logger
}

type AccountType struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ProviderDetailsRecord struct {
	Name         string        `json:"name"`
	Icon         string        `json:"icon"`
	AccountTypes []AccountType `json:"accountTypes"`
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
	ph.Logger.Debug("got providers from database", requestId)

	// Construct the response object
	providerDetailsMap := make(map[string][]AccountType)
	providerIconMap := make(map[string]string)
	for _, providerDetail := range providerDetails {
		accountType := AccountType{Id: providerDetail.AccountTypeId, Name: providerDetail.AccountTypeName}
		providerDetailsMap[providerDetail.Name] = append(providerDetailsMap[providerDetail.Name], accountType)
		providerIconMap[providerDetail.Name] = providerDetail.Icon
	}
	responseData := []ProviderDetailsRecord{}
	for providerName, accountTypes := range providerDetailsMap {
		responseData = append(responseData, ProviderDetailsRecord{
			Name:         providerName,
			Icon:         providerIconMap[providerName],
			AccountTypes: accountTypes,
		})
	}
	ph.Logger.Debug(fmt.Sprintf("constructed response object - %#v", responseData), requestId)

	return c.JSON(
		http.StatusOK,
		LooseJson{"success": true, "data": responseData},
	)
}
