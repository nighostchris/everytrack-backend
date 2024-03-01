package handlers

import (
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nighostchris/everytrack-backend/internal/database"
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
	AccountId  string `json:"accountId"`
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
		CurrencyId: data.CurrencyId,
		ExecutedAt: time.Unix(data.ExecutedAt, 0),
	}
	// Deal with nullable fields - remarks and accountId
	if len(data.Remarks) != 0 {
		createNewTransactionDbParams.Remarks = &data.Remarks
	}
	if len(data.AccountId) != 0 {
		createNewTransactionDbParams.AccountId = &data.AccountId
	}
	th.Logger.Debug(fmt.Sprintf("constructed parameters for create new transaction database query - %#v", createNewTransactionDbParams))

	// Check if there is enough balance to consume when accountId is not empty
	if createNewTransactionDbParams.AccountId != nil && !income {
		accountBalance, getAccountBalanceError := database.GetAccountBalance(th.Db, *createNewTransactionDbParams.AccountId)
		if getAccountBalanceError != nil {
			th.Logger.Error(
				fmt.Sprintf(
					"failed to get account balance for %s from database. %s",
					*createNewTransactionDbParams.AccountId,
					getAccountBalanceError.Error(),
				),
				requestId,
			)
			return c.JSON(
				http.StatusInternalServerError,
				LooseJson{"success": false, "error": "Internal server error."},
			)
		}

		// Converting balance to float number
		balanceInFloat, parseBalanceError := strconv.ParseFloat(accountBalance, 64)
		if parseBalanceError != nil {
			th.Logger.Error(fmt.Sprintf("failed to parse balance into float. %s", parseBalanceError.Error()), requestId)
			return c.JSON(
				http.StatusInternalServerError,
				LooseJson{"success": false, "error": "Internal server error."},
			)
		}

		// Converting amount to float number
		amountInFloat, parseAmountError := strconv.ParseFloat(data.Amount, 64)
		if parseAmountError != nil {
			th.Logger.Error(fmt.Sprintf("failed to parse amount into float. %s", parseAmountError.Error()), requestId)
			return c.JSON(
				http.StatusInternalServerError,
				LooseJson{"success": false, "error": "Internal server error."},
			)
		}

		// Compare balance with amount
		balanceComparison := big.NewFloat(amountInFloat).Cmp(big.NewFloat(balanceInFloat))
		if balanceComparison == 1 {
			return c.JSON(
				http.StatusBadRequest,
				LooseJson{"success": false, "error": "Insufficient account balance."},
			)
		}
		th.Logger.Debug(
			fmt.Sprintf("sufficient balance in account %s to pay the transaction amount", *createNewTransactionDbParams.AccountId),
			requestId,
		)

		// Calculate the final account balance after spending the transaction amount
		negativeAmount := big.NewFloat(0).Neg(big.NewFloat(amountInFloat))
		newAccountBalance := big.NewFloat(0).Add(big.NewFloat((balanceInFloat)), negativeAmount)
		newAccountBalanceInFloat, _ := newAccountBalance.Float64()

		_, updateAccountBalanceError := database.UpdateAccountBalance(
			th.Db,
			strconv.FormatFloat(newAccountBalanceInFloat, 'f', -1, 64),
			*createNewTransactionDbParams.AccountId,
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
	revertBalanceString := c.QueryParam("revertBalance")
	if len(revertBalanceString) > 0 {
		revertBalance, parseBooleanError := strconv.ParseBool(revertBalanceString)
		if parseBooleanError != nil {
			return c.JSON(
				http.StatusBadRequest,
				LooseJson{"success": false, "error": "Invalid query parameters."},
			)
		}

		if revertBalance {
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

			// Converting balance to float number
			balanceInFloat, parseBalanceError := strconv.ParseFloat(transactionAndBalance.Balance, 64)
			if parseBalanceError != nil {
				th.Logger.Error(fmt.Sprintf("failed to parse balance into float. %s", parseBalanceError.Error()), requestId)
				return c.JSON(
					http.StatusInternalServerError,
					LooseJson{"success": false, "error": "Internal server error."},
				)
			}

			// Converting amount to float number
			amountInFloat, parseAmountError := strconv.ParseFloat(transactionAndBalance.Amount, 64)
			if parseAmountError != nil {
				th.Logger.Error(fmt.Sprintf("failed to parse amount into float. %s", parseAmountError.Error()), requestId)
				return c.JSON(
					http.StatusInternalServerError,
					LooseJson{"success": false, "error": "Internal server error."},
				)
			}

			// Calculate the final account balance after reverting the transaction amount
			newAccountBalance, _ := big.NewFloat(0).Add(
				big.NewFloat((balanceInFloat)),
				big.NewFloat(amountInFloat),
			).Float64()

			_, updateAccountBalanceError := database.UpdateAccountBalance(
				th.Db,
				strconv.FormatFloat(newAccountBalance, 'f', -1, 64),
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
		}
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
