package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nighostchris/everytrack-backend/internal/database"
	"go.uber.org/zap"
)

type SettingsHandler struct {
	Db     *pgxpool.Pool
	Logger *zap.Logger
}

type UpdateSettingsRequestBody struct {
	CurrencyId string `json:"currencyId" validate:"required"`
	Username   string `json:"username" validate:"required,max=20"`
}

func (sh *SettingsHandler) GetAllClientSettings(c echo.Context) error {
	sh.Logger.Info("starts")

	clientId := c.Get("uid").(string)
	// Get client record from database
	client, getClientError := database.GetClientById(sh.Db, clientId)
	if getClientError != nil {
		sh.Logger.Error(fmt.Sprintf("failed to get client from database. %s", getClientError.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Internal server error."})
	}
	sh.Logger.Debug("got client from database")

	return c.JSON(http.StatusOK, map[string]interface{}{"success": true, "data": map[string]interface{}{"username": client.Username, "currencyId": client.CurrencyId}})
}

func (sh *SettingsHandler) UpdateSettings(c echo.Context) error {
	data := new(UpdateSettingsRequestBody)
	sh.Logger.Info("starts")

	// Retrieve request body and validate with schema
	if bindError := c.Bind(data); bindError != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Missing required fields"})
	}

	if validateError := c.Validate(data); validateError != nil {
		var ve validator.ValidationErrors
		if errors.As(validateError, &ve) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": fmt.Sprintf("Invalid field %s", strcase.ToLowerCamel(ve[0].Field()))})
		}
		sh.Logger.Error(fmt.Sprintf("invalid field. %s", validateError.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Invalid field"})
	}
	sh.Logger.Debug("validated request parameters")

	clientId := c.Get("uid").(string)

	// Update client settings in database
	_, updateError := database.UpdateClientSettings(sh.Db, database.UpdateClientSettingsParams{Username: data.Username, CurrencyId: data.CurrencyId, ClientId: clientId})
	if updateError != nil {
		sh.Logger.Error(fmt.Sprintf("failed to update client settings in database. %s", updateError.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Internal server error."})
	}
	sh.Logger.Debug("updated client settings in database")

	return c.JSON(http.StatusOK, map[string]interface{}{"success": true})
}
