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

type ExpensesHandler struct {
	Db     *pgxpool.Pool
	Logger *zap.Logger
}

type ExpenseRecord struct {
	Amount     string  `json:"amount"`
	Remarks    string  `json:"remarks"`
	Category   string  `json:"category"`
	AccountId  *string `json:"accountId"`
	ExecutedAt int64   `json:"executedAt"`
	CurrencyId string  `json:"currencyId"`
}

type CreateNewExpenseRequestBody struct {
	Amount     string `json:"amount"`
	Remarks    string `json:"remarks"`
	Category   string `json:"category"`
	AccountId  string `json:"accountId"`
	CurrencyId string `json:"currencyId"`
	ExecutedAt int64  `json:"executedAt"`
}

func (eh *ExpensesHandler) GetAllExpenses(c echo.Context) error {
	clientId := c.Get("uid").(string)
	requestId := zap.String("requestId", c.Get("requestId").(string))
	eh.Logger.Info("starts", requestId)

	// Get all expenses from database
	expenses, getExpensesError := database.GetAllExpenses(eh.Db, clientId)
	if getExpensesError != nil {
		eh.Logger.Error(
			fmt.Sprintf("failed to get all expense records from database. %s", getExpensesError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	eh.Logger.Debug("got expense records from database", requestId)

	// Construct the response object
	var expenseRecords []ExpenseRecord
	for _, expense := range expenses {
		record := ExpenseRecord{
			Amount:     expense.Amount,
			Category:   expense.Category,
			CurrencyId: expense.CurrencyId,
			Remarks:    expense.Remarks.String,
			ExecutedAt: expense.ExecutedAt.Unix(),
		}
		if len(expense.AccountId.String) > 0 {
			record.AccountId = &expense.AccountId.String
		}
		expenseRecords = append(expenseRecords, record)
	}
	eh.Logger.Debug(fmt.Sprintf("constructed response object - %#v", expenseRecords), requestId)

	return c.JSON(http.StatusOK, LooseJson{"success": true, "data": expenseRecords})
}

func (eh *ExpensesHandler) CreateNewExpense(c echo.Context) error {
	data := new(CreateNewExpenseRequestBody)
	clientId := c.Get("uid").(string)
	requestId := zap.String("requestId", c.Get("requestId").(string))
	eh.Logger.Info("starts", requestId)

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
		eh.Logger.Error(fmt.Sprintf("invalid field. %s", validateError.Error()), requestId)
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Invalid field"},
		)
	}
	eh.Logger.Debug("validated request parameters", requestId)

	// Construct database query parameters
	createNewExpenseDbParams := database.CreateNewExpenseParams{
		Amount:     data.Amount,
		ClientId:   clientId,
		Category:   data.Category,
		CurrencyId: data.CurrencyId,
		ExecutedAt: time.Unix(data.ExecutedAt, 0),
	}
	// Deal with nullable fields - remarks and accountId
	if len(data.Remarks) != 0 {
		createNewExpenseDbParams.Remarks = &data.Remarks
	}
	if len(data.AccountId) != 0 {
		createNewExpenseDbParams.AccountId = &data.AccountId
	}
	eh.Logger.Debug(fmt.Sprintf("constructed parameters for create new expense database query - %#v", createNewExpenseDbParams))

	// Check if there is enough balance to consume when accountId is not empty
	if createNewExpenseDbParams.AccountId != nil {
		accountBalance, getAccountBalanceError := database.GetAccountBalance(eh.Db, *createNewExpenseDbParams.AccountId)
		if getAccountBalanceError != nil {
			eh.Logger.Error(
				fmt.Sprintf(
					"failed to get account balance for %s from database. %s",
					*createNewExpenseDbParams.AccountId,
					getAccountBalanceError.Error(),
				),
				requestId,
			)
			return c.JSON(
				http.StatusInternalServerError,
				LooseJson{"success": false, "error": "Internal server error."},
			)
		}
		balanceInFloat, parseBalanceError := strconv.ParseFloat(accountBalance, 64)
		if parseBalanceError != nil {
			eh.Logger.Error(fmt.Sprintf("failed to parse balance into float. %s", parseBalanceError.Error()))
			return c.JSON(
				http.StatusInternalServerError,
				LooseJson{"success": false, "error": "Internal server error."},
			)
		}
		amountInFloat, parseAmountError := strconv.ParseFloat(data.Amount, 64)
		if parseAmountError != nil {
			eh.Logger.Error(fmt.Sprintf("failed to parse amount into float. %s", parseAmountError.Error()))
			return c.JSON(
				http.StatusInternalServerError,
				LooseJson{"success": false, "error": "Internal server error."},
			)
		}
		balanceComparison := big.NewFloat(amountInFloat).Cmp(big.NewFloat(balanceInFloat))
		if balanceComparison == 1 {
			return c.JSON(
				http.StatusBadRequest,
				LooseJson{"success": false, "error": "Insufficient account balance."},
			)
		}
		eh.Logger.Debug(fmt.Sprintf("sufficient balance in account %s to pay the expense amount", *createNewExpenseDbParams.AccountId))
	}
	eh.Logger.Debug("going to create new expense record in database")

	// Create new expense record in database
	_, createError := database.CreateNewExpense(eh.Db, createNewExpenseDbParams)
	if createError != nil {
		eh.Logger.Error(
			fmt.Sprintf("failed to create new expense record in database. %s", createError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	eh.Logger.Debug("created a new expense record in database", requestId)

	return c.JSON(http.StatusOK, LooseJson{"success": true})
}
