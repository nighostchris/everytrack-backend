package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nighostchris/everytrack-backend/internal/database"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

type AccountsHandler struct {
	Db     *pgxpool.Pool
	Logger *zap.Logger
}

type CreateNewAccountRequestBody struct {
	Name            string `json:"name" validate:"required"`
	CurrencyId      string `json:"currencyId" validate:"required"`
	AssetProviderId string `json:"assetProviderId" validate:"required"`
}

type TransferBetweenAccountsRequestBody struct {
	Amount          string `json:"amount" validate:"required"`
	SourceAccountId string `json:"sourceAccountId" validate:"required"`
	TargetAccountId string `json:"targetAccountId" validate:"required"`
}

type UpdateAccountRequestBody struct {
	Balance       string `json:"balance" validate:"required"`
	CurrencyId    string `json:"currencyId" validate:"required"`
	AccountTypeId string `json:"accountTypeId" validate:"required"`
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
	ah.Logger.Debug(fmt.Sprintf("got accounts from database - %#v", accountSummary), requestId)

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

	// Check if account name in use already
	accountNameInUse, checkExistingAccountError := database.CheckExistingAccount(
		ah.Db,
		database.CheckExistingAccountParams{
			Name:            data.Name,
			ClientId:        clientId,
			AssetProviderId: data.AssetProviderId,
		},
	)
	if checkExistingAccountError != nil {
		ah.Logger.Error(
			fmt.Sprintf("failed to check existing account in database. %s", checkExistingAccountError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	if accountNameInUse {
		return c.JSON(
			http.StatusConflict,
			LooseJson{"success": false, "error": "Account name already in use."},
		)
	}

	// Create a new account in database
	_, createError := database.CreateNewAccount(
		ah.Db,
		database.CreateNewAccountParams{
			ClientId:        clientId,
			Name:            data.Name,
			CurrencyId:      data.CurrencyId,
			AssetProviderId: data.AssetProviderId,
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

func (ah *AccountsHandler) TransferBetweenAccounts(c echo.Context) error {
	data := new(TransferBetweenAccountsRequestBody)
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

	// Get source account name and balance
	sourceAccountSummary, getSourceAccountSummaryError := database.GetAccountSummary(ah.Db, data.SourceAccountId)
	if getSourceAccountSummaryError != nil {
		ah.Logger.Error(
			fmt.Sprintf("failed to get source account summary from database. %s", getSourceAccountSummaryError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	ah.Logger.Debug("got source account summary from database", requestId)

	// Get target account name and balance
	targetAccountSummary, getTargetAccountSummaryError := database.GetAccountSummary(ah.Db, data.TargetAccountId)
	if getTargetAccountSummaryError != nil {
		ah.Logger.Error(
			fmt.Sprintf("failed to get target account summary from database. %s", getTargetAccountSummaryError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	ah.Logger.Debug("got target account summary from database", requestId)

	// Converting source account balance to decimal
	sourceAccountBalanceInDecimal, parseSourceAccountBalanceError := decimal.NewFromString(sourceAccountSummary.Balance)
	if parseSourceAccountBalanceError != nil {
		ah.Logger.Error(fmt.Sprintf("failed to parse source account balance into decimal. %s", parseSourceAccountBalanceError.Error()), requestId)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}

	// Converting target account balance to decimal
	targetAccountBalanceInDecimal, parseTargetAccountBalanceError := decimal.NewFromString(targetAccountSummary.Balance)
	if parseTargetAccountBalanceError != nil {
		ah.Logger.Error(fmt.Sprintf("failed to parse target account balance into decimal. %s", parseTargetAccountBalanceError.Error()), requestId)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}

	// Converting amount to decimal
	amountInDecimal, parseAmountError := decimal.NewFromString(data.Amount)
	if parseAmountError != nil {
		ah.Logger.Error(fmt.Sprintf("failed to parse amount into decimal. %s", parseAmountError.Error()), requestId)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}

	// Calculate the final account balance after receiving / spending the transfer amount
	newSourceAccountBalance := sourceAccountBalanceInDecimal.Sub(amountInDecimal)
	newTargetAccountBalance := targetAccountBalanceInDecimal.Add(amountInDecimal)
	ah.Logger.Debug(
		fmt.Sprintf("source account balance for %s will go from %s to %s", data.SourceAccountId, sourceAccountSummary.Balance, newSourceAccountBalance.Truncate(2).String()),
		requestId,
	)
	ah.Logger.Debug(
		fmt.Sprintf("target account balance for %s will go from %s to %s", data.TargetAccountId, targetAccountSummary.Balance, newTargetAccountBalance.Truncate(2).String()),
		requestId,
	)

	// Update source account balance in database
	_, updateSourceAccountError := database.UpdateAccount(
		ah.Db,
		database.UpdateAccountParams{
			Balance:       newSourceAccountBalance.Truncate(2).String(),
			CurrencyId:    sourceAccountSummary.CurrencyId,
			AccountTypeId: sourceAccountSummary.AccountTypeId,
		},
	)
	if updateSourceAccountError != nil {
		ah.Logger.Error(
			fmt.Sprintf("failed to update source account balance in database. %s", updateSourceAccountError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	ah.Logger.Debug("updated source account balance in database", requestId)

	// Update target account balance in database
	_, updateTargetAccountError := database.UpdateAccount(
		ah.Db,
		database.UpdateAccountParams{
			Balance:       newTargetAccountBalance.Truncate(2).String(),
			CurrencyId:    targetAccountSummary.CurrencyId,
			AccountTypeId: targetAccountSummary.AccountTypeId,
		},
	)
	if updateTargetAccountError != nil {
		ah.Logger.Error(
			fmt.Sprintf("failed to update target account balance in database. %s", updateTargetAccountError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	ah.Logger.Debug("updated target account balance in database", requestId)

	// Create expense transaction for source account
	_, createSourceAccountTransactionError := database.CreateNewTransaction(ah.Db, database.CreateNewTransactionParams{
		Income:     false,
		Name:       fmt.Sprintf("Transfer to %s", targetAccountSummary.Name),
		Amount:     data.Amount,
		ClientId:   clientId,
		Category:   "bank-transfer",
		AccountId:  data.SourceAccountId,
		CurrencyId: sourceAccountSummary.CurrencyId,
		ExecutedAt: time.Now().Truncate(24 * time.Hour),
	})
	ah.Logger.Debug("created a new transaction record for source account in database", requestId)

	if createSourceAccountTransactionError != nil {
		ah.Logger.Error(
			fmt.Sprintf("failed to create new transaction record for source account in database. %s", createSourceAccountTransactionError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}

	// Create incoming transaction for target account
	_, createTargetAccountTransactionError := database.CreateNewTransaction(ah.Db, database.CreateNewTransactionParams{
		Income:     true,
		Name:       fmt.Sprintf("Received from %s", sourceAccountSummary.Name),
		Amount:     data.Amount,
		ClientId:   clientId,
		Category:   "bank-transfer",
		AccountId:  data.TargetAccountId,
		CurrencyId: targetAccountSummary.CurrencyId,
		ExecutedAt: time.Now().Truncate(24 * time.Hour),
	})
	ah.Logger.Debug("created a new transaction record for target account in database", requestId)

	if createTargetAccountTransactionError != nil {
		ah.Logger.Error(
			fmt.Sprintf("failed to create new transaction record for target account in database. %s", createTargetAccountTransactionError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}

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

func (ah *AccountsHandler) DeleteAccount(c echo.Context) error {
	clientId := c.Get("uid").(string)
	requestId := zap.String("requestId", c.Get("requestId").(string))
	ah.Logger.Info("starts", requestId)

	accountId := c.QueryParam("id")
	if len(accountId) == 0 {
		ah.Logger.Error("undefined account id", requestId)
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Undefined account id."},
		)
	}
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
	ah.Logger.Info(fmt.Sprintf("going to check if client owns the account with id %s", accountId), requestId)

	// Get all client owned accounts in database
	ownedAccounts, getOwnedAccountsError := database.GetAllAccountSummaryByType(ah.Db, providerType, clientId)
	if getOwnedAccountsError != nil {
		ah.Logger.Error(
			fmt.Sprintf("failed to get all owned accounts from database. %s", getOwnedAccountsError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	if len(ownedAccounts) == 0 {
		ah.Logger.Error(
			fmt.Sprintf("client does not own the account %s and not authorized to delete", accountId),
			requestId,
		)
		return c.JSON(
			http.StatusUnauthorized,
			LooseJson{"success": false, "error": "Unauthorized to delete account."},
		)
	}
	ah.Logger.Info(fmt.Sprintf("going to delete account with id %s", accountId), requestId)

	// Delete account in database
	_, deleteError := database.DeleteAccount(ah.Db, accountId)
	if deleteError != nil {
		ah.Logger.Error(
			fmt.Sprintf("failed to delete account in database. %s", deleteError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	ah.Logger.Debug("deleted account in database", requestId)

	return c.JSON(http.StatusOK, LooseJson{"success": true})
}
