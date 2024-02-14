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
	"go.uber.org/zap"
)

type FuturePaymentsHandler struct {
	Db     *pgxpool.Pool
	Logger *zap.Logger
}

type FuturePaymentRecord struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Amount      string `json:"amount"`
	Income      bool   `json:"income"`
	Rolling     bool   `json:"rolling"`
	Remarks     string `json:"remarks"`
	Frequency   *int64 `json:"frequency"`
	AccountId   string `json:"accountId"`
	CurrencyId  string `json:"currencyId"`
	ScheduledAt int64  `json:"scheduledAt"`
}

type CreateNewFuturePaymentRequestBody struct {
	Name        string `json:"name"`
	Amount      string `json:"amount"`
	Income      string `json:"income"`
	Rolling     string `json:"rolling"`
	Remarks     string `json:"remarks"`
	Frequency   int64  `json:"frequency"`
	AccountId   string `json:"accountId"`
	CurrencyId  string `json:"currencyId"`
	ScheduledAt int64  `json:"scheduledAt"`
}

type UpdateFuturePaymentRequestBody struct {
	Id          string `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Income      string `json:"income" validate:"required"`
	Amount      string `json:"amount" validate:"required"`
	Remarks     string `json:"remarks"`
	Rolling     string `json:"rolling" validate:"required"`
	Frequency   int64  `json:"frequency"`
	AccountId   string `json:"accountId" validate:"required"`
	CurrencyId  string `json:"currencyId" validate:"required"`
	ScheduledAt int64  `json:"scheduledAt" validate:"required"`
}

func (fph *FuturePaymentsHandler) GetAllFuturePayments(c echo.Context) error {
	clientId := c.Get("uid").(string)
	requestId := zap.String("requestId", c.Get("requestId").(string))
	fph.Logger.Info("starts", requestId)

	// Get all future payments from database
	futurePayments, getFuturePaymentsError := database.GetAllFuturePaymentsByClientId(fph.Db, clientId)
	if getFuturePaymentsError != nil {
		fph.Logger.Error(
			fmt.Sprintf("failed to get all future payment records from database. %s", getFuturePaymentsError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	fph.Logger.Debug("got expense records from database", requestId)

	// Construct the response object
	var futurePaymentRecords []FuturePaymentRecord
	for _, futurePayment := range futurePayments {
		record := FuturePaymentRecord{
			Id:          futurePayment.Id,
			Name:        futurePayment.Name,
			Amount:      futurePayment.Amount,
			Income:      futurePayment.Income,
			Rolling:     futurePayment.Rolling,
			AccountId:   futurePayment.AccountId,
			CurrencyId:  futurePayment.CurrencyId,
			Remarks:     futurePayment.Remarks.String,
			ScheduledAt: futurePayment.ScheduledAt.Unix(),
		}
		if futurePayment.Frequency.Valid {
			record.Frequency = &futurePayment.Frequency.Int64
		}
		futurePaymentRecords = append(futurePaymentRecords, record)
	}
	fph.Logger.Debug(fmt.Sprintf("constructed response object - %#v", futurePaymentRecords), requestId)

	return c.JSON(http.StatusOK, LooseJson{"success": true, "data": futurePaymentRecords})
}

func (fph *FuturePaymentsHandler) CreateNewFuturePayment(c echo.Context) error {
	data := new(CreateNewFuturePaymentRequestBody)
	clientId := c.Get("uid").(string)
	requestId := zap.String("requestId", c.Get("requestId").(string))
	fph.Logger.Info("starts", requestId)

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
		fph.Logger.Error(fmt.Sprintf("invalid field. %s", validateError.Error()), requestId)
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

	rolling, parseRollingError := strconv.ParseBool(data.Rolling)
	if parseRollingError != nil {
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Invalid field rolling"},
		)
	}

	fph.Logger.Debug("validated request parameters", requestId)

	// Throw error if the payment is on rolling basis but upstream does not send payment frequency as well
	if rolling && data.Frequency < 1 {
		fph.Logger.Error("missing frequency when payment is on rolling basis.", requestId)
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Missing required field frequency"},
		)
	}

	// Throw error if the scheduled payment time is on today or previous days
	if data.ScheduledAt <= time.Now().Unix() {
		fph.Logger.Error("the payment schedule time is earlier than current date.", requestId)
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Invalid payment schedule time"},
		)
	}

	// Construct database query parameters
	createNewFuturePaymentDbParams := database.CreateNewFuturePaymentParams{
		Name:        data.Name,
		Amount:      data.Amount,
		Income:      income,
		Rolling:     rolling,
		ClientId:    clientId,
		AccountId:   data.AccountId,
		CurrencyId:  data.CurrencyId,
		ScheduledAt: time.Unix(data.ScheduledAt, 0),
	}
	// Deal with nullable fields - remarks and frequency
	if len(data.Remarks) != 0 {
		createNewFuturePaymentDbParams.Remarks = &data.Remarks
	}
	if data.Frequency > 0 {
		createNewFuturePaymentDbParams.Frequency = &data.Frequency
	}
	fph.Logger.Debug(fmt.Sprintf("constructed parameters for create new future payment database query - %#v", createNewFuturePaymentDbParams))

	// Create new future payment record in database
	_, createError := database.CreateNewFuturePayment(fph.Db, createNewFuturePaymentDbParams)
	if createError != nil {
		fph.Logger.Error(
			fmt.Sprintf("failed to create new future payment record in database. %s", createError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	fph.Logger.Debug("created a new future payment record in database", requestId)

	return c.JSON(http.StatusOK, LooseJson{"success": true})
}

func (fph *FuturePaymentsHandler) UpdateFuturePayment(c echo.Context) error {
	data := new(UpdateFuturePaymentRequestBody)
	requestId := zap.String("requestId", c.Get("requestId").(string))
	fph.Logger.Info("starts", requestId)

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
		fph.Logger.Error(fmt.Sprintf("invalid field. %s", validateError.Error()), requestId)
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

	rolling, parseRollingError := strconv.ParseBool(data.Rolling)
	if parseRollingError != nil {
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Invalid field rolling"},
		)
	}

	fph.Logger.Debug("validated request parameters", requestId)

	// Update account in database
	_, updateError := database.UpdateFuturePayment(
		fph.Db,
		database.UpdateFuturePaymentParams{
			Id:          data.Id,
			Name:        data.Name,
			Income:      income,
			Amount:      data.Amount,
			Remarks:     &data.Remarks,
			Rolling:     rolling,
			Frequency:   &data.Frequency,
			AccountId:   data.AccountId,
			CurrencyId:  data.CurrencyId,
			ScheduledAt: time.Unix(data.ScheduledAt, 0),
		},
	)
	if updateError != nil {
		fph.Logger.Error(
			fmt.Sprintf("failed to update future payment in database. %s", updateError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	fph.Logger.Debug("updated future payment in database", requestId)

	return c.JSON(http.StatusOK, LooseJson{"success": true})
}

func (fph *FuturePaymentsHandler) DeleteFuturePayment(c echo.Context) error {
	clientId := c.Get("uid").(string)
	requestId := zap.String("requestId", c.Get("requestId").(string))
	fph.Logger.Info("starts", requestId)

	futurePaymentId := c.QueryParam("id")
	if len(futurePaymentId) == 0 {
		fph.Logger.Error("undefined future payment id", requestId)
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Undefined future payment id."},
		)
	}

	// Delete future payment record in database
	_, deleteError := database.DeleteFuturePayment(fph.Db, futurePaymentId, clientId)
	if deleteError != nil {
		fph.Logger.Error(
			fmt.Sprintf("failed to delete future payment in database. %s", deleteError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	fph.Logger.Debug("deleted future payment record in database", requestId)

	return c.JSON(http.StatusOK, LooseJson{"success": true})
}
