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

type SavingsHandler struct {
	Db     *pgxpool.Pool
	Logger *zap.Logger
}

type CreateNewAccountRequestBody struct {
	CurrencyId    string `json:"currencyId"`
	AccountTypeId string `json:"accountTypeId"`
}

type UpdateAccountRequestBody struct {
	Balance       string `json:"balance"`
	CurrencyId    string `json:"currencyId"`
	AccountTypeId string `json:"accountTypeId"`
}

func (sh *SavingsHandler) GetAllBankAccounts(c echo.Context) error {
	sh.Logger.Info("starts")

	clientId := c.Get("uid").(string)

	// Get all bank details from database
	bankAccounts, getBankAccountsError := database.GetAllBankAccounts(sh.Db, clientId)

	if getBankAccountsError != nil {
		sh.Logger.Error(fmt.Sprintf("failed to get all bank accounts from database. %s", getBankAccountsError.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Internal server error."})
	}
	sh.Logger.Debug("got bank accounts from database")

	return c.JSON(http.StatusOK, map[string]interface{}{"success": true, "data": bankAccounts})
}

func (sh *SavingsHandler) CreateNewAccount(c echo.Context) error {
	data := new(CreateNewAccountRequestBody)
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

	// Create a new account in database
	_, createError := database.CreateNewAccount(sh.Db, database.CreateNewAccountParams{ClientId: clientId, CurrencyId: data.CurrencyId, AccountTypeId: data.AccountTypeId})
	if createError != nil {
		sh.Logger.Error(fmt.Sprintf("failed to create new account in database. %s", createError.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Internal server error."})
	}
	sh.Logger.Debug("created a new account in database")

	return c.JSON(http.StatusOK, map[string]interface{}{"success": true})
}

func (sh *SavingsHandler) UpdateAccount(c echo.Context) error {
	data := new(UpdateAccountRequestBody)
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

	// Update account in database
	_, updateError := database.UpdateAccount(sh.Db, database.UpdateAccountParams{Balance: data.Balance, CurrencyId: data.CurrencyId, AccountTypeId: data.AccountTypeId})
	if updateError != nil {
		sh.Logger.Error(fmt.Sprintf("failed to update account in database. %s", updateError.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Internal server error."})
	}
	sh.Logger.Debug("updated account in database")

	return c.JSON(http.StatusOK, map[string]interface{}{"success": true})
}
