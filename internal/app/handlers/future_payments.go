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
	Income      bool   `json:"income"`
	Rolling     bool   `json:"rolling"`
	Remarks     string `json:"remarks"`
	Frequency   int64  `json:"frequency"`
	AccountId   string `json:"accountId"`
	CurrencyId  string `json:"currencyId"`
	ScheduledAt int64  `json:"scheduledAt"`
}

func (eh *FuturePaymentsHandler) GetAllFuturePayments(c echo.Context) error {
	clientId := c.Get("uid").(string)
	requestId := zap.String("requestId", c.Get("requestId").(string))
	eh.Logger.Info("starts", requestId)

	// Get all future payments from database
	futurePayments, getFuturePaymentsError := database.GetAllFuturePayments(eh.Db, &clientId)
	if getFuturePaymentsError != nil {
		eh.Logger.Error(
			fmt.Sprintf("failed to get all future payment records from database. %s", getFuturePaymentsError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	eh.Logger.Debug("got expense records from database", requestId)

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
	eh.Logger.Debug(fmt.Sprintf("constructed response object - %#v", futurePaymentRecords), requestId)

	return c.JSON(http.StatusOK, LooseJson{"success": true, "data": futurePaymentRecords})
}

func (eh *FuturePaymentsHandler) CreateNewFuturePayment(c echo.Context) error {
	data := new(CreateNewFuturePaymentRequestBody)
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

	// Throw error if the payment is on rolling basis but upstream does not send payment frequency as well
	if data.Rolling && data.Frequency < 1 {
		eh.Logger.Error("missing frequency when payment is on rolling basis.", requestId)
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Missing required field frequency"},
		)
	}

	// Throw error if the scheduled payment time is on today or previous days
	if data.ScheduledAt <= time.Now().Unix() {
		eh.Logger.Error("the payment schedule time is earlier than current date.", requestId)
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Invalid payment schedule time"},
		)
	}

	// Construct database query parameters
	createNewFuturePaymentDbParams := database.CreateNewFuturePaymentParams{
		Name:        data.Name,
		Amount:      data.Amount,
		Income:      data.Income,
		Rolling:     data.Rolling,
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
	eh.Logger.Debug(fmt.Sprintf("constructed parameters for create new future payment database query - %#v", createNewFuturePaymentDbParams))

	// Create new future payment record in database
	_, createError := database.CreateNewFuturePayment(eh.Db, createNewFuturePaymentDbParams)
	if createError != nil {
		eh.Logger.Error(
			fmt.Sprintf("failed to create new future payment record in database. %s", createError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	eh.Logger.Debug("created a new future payment record in database", requestId)

	return c.JSON(http.StatusOK, LooseJson{"success": true})
}

func (eh *FuturePaymentsHandler) DeleteFuturePayment(c echo.Context) error {
	clientId := c.Get("uid").(string)
	requestId := zap.String("requestId", c.Get("requestId").(string))
	eh.Logger.Info("starts", requestId)

	futurePaymentId := c.QueryParam("id")
	if len(futurePaymentId) == 0 {
		eh.Logger.Error("undefined future payment id", requestId)
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Undefined future payment id."},
		)
	}

	// Delete future payment record in database
	_, deleteError := database.DeleteFuturePayment(eh.Db, futurePaymentId, clientId)
	if deleteError != nil {
		eh.Logger.Error(
			fmt.Sprintf("failed to delete future payment in database. %s", deleteError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	eh.Logger.Debug("deleted future payment record in database", requestId)

	return c.JSON(http.StatusOK, LooseJson{"success": true})
}
