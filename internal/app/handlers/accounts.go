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
	"golang.org/x/exp/slices"
)

type AccountsHandler struct {
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

func (ah *AccountsHandler) GetAllAccountsByType(c echo.Context) error {
	clientId := c.Get("uid").(string)
	requestId := zap.String("requestId", c.Get("requestId").(string))
	ah.Logger.Info("starts", requestId)

	providerType := c.QueryParam("type")
	if len(providerType) == 0 {
		ah.Logger.Error("undefined provider type", requestId)
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Undefined provider type."},
		)
	}
	if !slices.Contains(ProviderTypes, providerType) {
		ah.Logger.Error(fmt.Sprintf("invalid provider type %s", providerType), requestId)
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Invalid provider type."},
		)
	}
	ah.Logger.Info(fmt.Sprintf("going to get all account summary by type %s", providerType), requestId)

	// Get all accounts by provider type from database
	accountSummary, getAccountSummaryError := database.GetAllAccountSummaryByType(ah.Db, providerType, clientId)
	if getAccountSummaryError != nil {
		ah.Logger.Error(
			fmt.Sprintf("failed to get all account summary from database. %s", getAccountSummaryError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	ah.Logger.Debug(fmt.Sprintf("got bank accounts from database - %#v", accountSummary), requestId)

	return c.JSON(http.StatusOK, LooseJson{"success": true, "data": accountSummary})
}

func (ah *AccountsHandler) CreateNewAccount(c echo.Context) error {
	data := new(CreateNewAccountRequestBody)
	clientId := c.Get("uid").(string)
	requestId := zap.String("requestId", c.Get("requestId").(string))
	ah.Logger.Info("starts", requestId)

	// Retrieve request body and validate with schema
	if bindError := c.Bind(data); bindError != nil {
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Missing required fields"},
		)
	}

	if validateError := c.Validate(data); validateError != nil {
		var ve validator.ValidationErrors
		if errors.As(validateError, &ve) {
			return c.JSON(
				http.StatusBadRequest,
				LooseJson{"success": false, "error": fmt.Sprintf("Invalid field %s", strcase.ToLowerCamel(ve[0].Field()))},
			)
		}
		ah.Logger.Error(fmt.Sprintf("invalid field. %s", validateError.Error()), requestId)
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Invalid field"},
		)
	}
	ah.Logger.Debug("validated request parameters", requestId)

	// Create a new account in database
	_, createError := database.CreateNewAccount(
		ah.Db,
		database.CreateNewAccountParams{
			ClientId:      clientId,
			CurrencyId:    data.CurrencyId,
			AccountTypeId: data.AccountTypeId,
		},
	)
	if createError != nil {
		ah.Logger.Error(
			fmt.Sprintf("failed to create new account in database. %s", createError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	ah.Logger.Debug("created a new account in database", requestId)

	return c.JSON(http.StatusOK, LooseJson{"success": true})
}

func (ah *AccountsHandler) UpdateAccount(c echo.Context) error {
	data := new(UpdateAccountRequestBody)
	requestId := zap.String("requestId", c.Get("requestId").(string))
	ah.Logger.Info("starts", requestId)

	// Retrieve request body and validate with schema
	if bindError := c.Bind(data); bindError != nil {
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Missing required fields"},
		)
	}

	if validateError := c.Validate(data); validateError != nil {
		var ve validator.ValidationErrors
		if errors.As(validateError, &ve) {
			return c.JSON(
				http.StatusBadRequest,
				LooseJson{"success": false, "error": fmt.Sprintf("Invalid field %s", strcase.ToLowerCamel(ve[0].Field()))},
			)
		}
		ah.Logger.Error(fmt.Sprintf("invalid field. %s", validateError.Error()), requestId)
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Invalid field"},
		)
	}
	ah.Logger.Debug("validated request parameters", requestId)

	// Update account in database
	_, updateError := database.UpdateAccount(
		ah.Db,
		database.UpdateAccountParams{
			Balance:       data.Balance,
			CurrencyId:    data.CurrencyId,
			AccountTypeId: data.AccountTypeId,
		},
	)
	if updateError != nil {
		ah.Logger.Error(
			fmt.Sprintf("failed to update account in database. %s", updateError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	ah.Logger.Debug("updated account in database", requestId)

	return c.JSON(http.StatusOK, LooseJson{"success": true})
}
