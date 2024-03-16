package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nighostchris/everytrack-backend/internal/database"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type TransactionsHandler struct {
	Db     *pgxpool.Pool
	Logger *zap.Logger
}

type TransactionRecord struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Income     bool    `json:"income"`
	Amount     string  `json:"amount"`
	Remarks    string  `json:"remarks"`
	Category   string  `json:"category"`
	AccountId  *string `json:"accountId"`
	ExecutedAt int64   `json:"executedAt"`
	CurrencyId string  `json:"currencyId"`
}

type CreateNewTransactionRequestBody struct {
	Name       string `json:"name" validate:"required"`
	Income     string `json:"income" validate:"required"`
	Amount     string `json:"amount" validate:"required"`
	Remarks    string `json:"remarks"`
	Category   string `json:"category" validate:"required"`
	AccountId  string `json:"accountId" validate:"required"`
	CurrencyId string `json:"currencyId" validate:"required"`
	ExecutedAt int64  `json:"executedAt" validate:"required"`
}

func (th *TransactionsHandler) GetAllTransactions(c echo.Context) error {
	clientId := c.Get("uid").(string)
	requestId := zap.String("requestId", c.Get("requestId").(string))
	th.Logger.Info("starts", requestId)

	// Get all transactions from database
	transactions, getTransactionsError := database.GetAllTransactions(th.Db, clientId)
	if getTransactionsError != nil {
		th.Logger.Error(
			fmt.Sprintf("failed to get all transaction records from database. %s", getTransactionsError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	th.Logger.Debug("got transaction records from database", requestId)

	// Construct the response object
	var transactionRecords []TransactionRecord
	for _, transaction := range transactions {
		record := TransactionRecord{
			Id:         transaction.Id,
			Name:       transaction.Name,
			Income:     transaction.Income,
			Amount:     transaction.Amount,
			Category:   transaction.Category,
			CurrencyId: transaction.CurrencyId,
			Remarks:    transaction.Remarks.String,
			ExecutedAt: transaction.ExecutedAt.Unix(),
		}
		if len(transaction.AccountId.String) > 0 {
			accountId := transaction.AccountId.String
			record.AccountId = &accountId
		}
		transactionRecords = append(transactionRecords, record)
	}
	th.Logger.Debug(fmt.Sprintf("constructed response object - %#v", transactionRecords), requestId)

	return c.JSON(http.StatusOK, LooseJson{"success": true, "data": transactionRecords})
}

func (th *TransactionsHandler) CreateNewTransaction(c echo.Context) error {
	data := new(CreateNewTransactionRequestBody)
	clientId := c.Get("uid").(string)
	requestId := zap.String("requestId", c.Get("requestId").(string))
	th.Logger.Info("starts", requestId)

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
		th.Logger.Error(fmt.Sprintf("invalid field. %s", validateError.Error()), requestId)
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Invalid field"},
		)
	}

	income, parseIncomeError := strconv.ParseBool(data.Income)
	if parseIncomeError != nil {
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Invalid field income"},
		)
	}

	th.Logger.Debug("validated request parameters", requestId)

	// Construct database query parameters
	createNewTransactionDbParams := database.CreateNewTransactionParams{
		Income:     income,
		Name:       data.Name,
		Amount:     data.Amount,
		ClientId:   clientId,
		Category:   data.Category,
		AccountId:  data.AccountId,
		CurrencyId: data.CurrencyId,
		ExecutedAt: time.Unix(data.ExecutedAt, 0),
	}
	// Deal with nullable fields - remarks and accountId
	if len(data.Remarks) != 0 {
		createNewTransactionDbParams.Remarks = &data.Remarks
	}
	th.Logger.Debug(fmt.Sprintf("constructed parameters for create new transaction database query - %#v", createNewTransactionDbParams))

	// Get original account balance
	accountBalance, getAccountBalanceError := database.GetAccountBalance(th.Db, createNewTransactionDbParams.AccountId)
	if getAccountBalanceError != nil {
		th.Logger.Error(
			fmt.Sprintf(
				"failed to get account balance for %s from database. %s",
				createNewTransactionDbParams.AccountId,
				getAccountBalanceError.Error(),
			),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}

	// Converting balance to decimal
	balanceInDecimal, parseBalanceError := decimal.NewFromString(accountBalance)
	if parseBalanceError != nil {
		th.Logger.Error(fmt.Sprintf("failed to parse balance into decimal. %s", parseBalanceError.Error()), requestId)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}

	// Converting amount to decimal
	amountInDecimal, parseAmountError := decimal.NewFromString(data.Amount)
	if parseAmountError != nil {
		th.Logger.Error(fmt.Sprintf("failed to parse amount into decimal. %s", parseAmountError.Error()), requestId)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}

	// Calculate the final account balance after receiving / spending the transaction amount
	var newAccountBalance decimal.Decimal
	if income {
		newAccountBalance = balanceInDecimal.Add(amountInDecimal)
	} else {
		newAccountBalance = balanceInDecimal.Sub(amountInDecimal)
	}
	th.Logger.Debug(
		fmt.Sprintf("account balance for %s will go from %s to %s", createNewTransactionDbParams.AccountId, accountBalance, newAccountBalance.Truncate(2).String()),
		requestId,
	)

	_, updateAccountBalanceError := database.UpdateAccountBalance(
		th.Db,
		newAccountBalance.Truncate(2).String(),
		createNewTransactionDbParams.AccountId,
	)
	if updateAccountBalanceError != nil {
		th.Logger.Error(
			fmt.Sprintf("failed to update latest balance after deducting transaction amount. %s", updateAccountBalanceError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	th.Logger.Debug("going to create new transaction record in database")

	// Create new transaction record in database
	_, createError := database.CreateNewTransaction(th.Db, createNewTransactionDbParams)
	if createError != nil {
		th.Logger.Error(
			fmt.Sprintf("failed to create new transaction record in database. %s", createError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	th.Logger.Debug("created a new transaction record in database", requestId)

	return c.JSON(http.StatusOK, LooseJson{"success": true})
}

func (th *TransactionsHandler) DeleteTransaction(c echo.Context) error {
	clientId := c.Get("uid").(string)
	requestId := zap.String("requestId", c.Get("requestId").(string))
	th.Logger.Info("starts", requestId)

	transactionId := c.QueryParam("id")
	if len(transactionId) == 0 {
		th.Logger.Error("undefined transaction id", requestId)
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Undefined transaction id."},
		)
	}

	transactionAndBalance, getTransactionAndBalanceError := database.GetTransactionAndAccountBalanceById(th.Db, transactionId)
	if getTransactionAndBalanceError != nil {
		th.Logger.Error(
			fmt.Sprintf("failed to get transaction and balance from database. %s", getTransactionAndBalanceError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}

	// Converting balance to decimal
	balanceInDecimal, parseBalanceError := decimal.NewFromString(transactionAndBalance.Balance)
	if parseBalanceError != nil {
		th.Logger.Error(fmt.Sprintf("failed to parse balance into decimal. %s", parseBalanceError.Error()), requestId)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}

	// Converting amount to decimal
	amountInDecimal, parseAmountError := decimal.NewFromString(transactionAndBalance.Amount)
	if parseAmountError != nil {
		th.Logger.Error(fmt.Sprintf("failed to parse amount into decimal. %s", parseAmountError.Error()), requestId)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}

	// Calculate the final account balance after reverting the transaction amount
	var newAccountBalance decimal.Decimal
	if transactionAndBalance.Income {
		newAccountBalance = balanceInDecimal.Sub(amountInDecimal)
	} else {
		newAccountBalance = balanceInDecimal.Add(amountInDecimal)
	}
	th.Logger.Debug(
		fmt.Sprintf("account balance for %s will go from %s to %s", transactionAndBalance.AccountId, transactionAndBalance.Balance, newAccountBalance.Truncate(2).String()),
		requestId,
	)

	_, updateAccountBalanceError := database.UpdateAccountBalance(
		th.Db,
		newAccountBalance.Truncate(2).String(),
		transactionAndBalance.AccountId,
	)
	if updateAccountBalanceError != nil {
		th.Logger.Error(
			fmt.Sprintf("failed to update latest balance after reverting transaction amount. %s", updateAccountBalanceError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}

	// Delete transaction record in database
	_, deleteError := database.DeleteTransaction(th.Db, transactionId, clientId)
	if deleteError != nil {
		th.Logger.Error(
			fmt.Sprintf("failed to delete transaction in database. %s", deleteError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	th.Logger.Debug("deleted transaction record in database", requestId)

	return c.JSON(http.StatusOK, LooseJson{"success": true})
}
