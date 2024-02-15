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

type CashHandler struct {
	Db     *pgxpool.Pool
	Logger *zap.Logger
}

type CashRecord struct {
	Id         string `json:"id"`
	Amount     string `json:"amount"`
	CurrencyId string `json:"currencyId"`
}

type CreateNewCashRecordRequestBody struct {
	Amount     string `json:"amount" validate:"required"`
	CurrencyId string `json:"currencyId" validate:"required"`
}

type UpdateCashRecordRequestBody struct {
	Id         string `json:"id" validate:"required"`
	Amount     string `json:"amount" validate:"required"`
	CurrencyId string `json:"currencyId" validate:"required"`
}

func (ch *CashHandler) GetAllCash(c echo.Context) error {
	clientId := c.Get("uid").(string)
	requestId := zap.String("requestId", c.Get("requestId").(string))
	ch.Logger.Info("starts", requestId)

	// Get all cash records from database
	cash, getCashError := database.GetAllCash(ch.Db, clientId)
	if getCashError != nil {
		ch.Logger.Error(
			fmt.Sprintf("failed to get all cash records from database. %s", getCashError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	ch.Logger.Debug("got cash records from database", requestId)

	// Construct the response object
	cashRecords := []CashRecord{}
	for _, dbCash := range cash {
		record := CashRecord{
			Id:         dbCash.Id,
			Amount:     dbCash.Amount,
			CurrencyId: dbCash.CurrencyId,
		}
		cashRecords = append(cashRecords, record)
	}
	ch.Logger.Debug(fmt.Sprintf("constructed response object - %#v", cashRecords), requestId)

	return c.JSON(http.StatusOK, LooseJson{"success": true, "data": cashRecords})
}

func (ch *CashHandler) CreateNewCashRecord(c echo.Context) error {
	data := new(CreateNewCashRecordRequestBody)
	clientId := c.Get("uid").(string)
	requestId := zap.String("requestId", c.Get("requestId").(string))
	ch.Logger.Info("starts", requestId)

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
		ch.Logger.Error(fmt.Sprintf("invalid field. %s", validateError.Error()), requestId)
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Invalid field"},
		)
	}
	ch.Logger.Debug("validated request parameters", requestId)

	// Create new cash record in database
	_, createError := database.CreateNewCashRecord(ch.Db, database.CreateNewCashRecordParams{
		ClientId:   clientId,
		Amount:     data.Amount,
		CurrencyId: data.CurrencyId,
	})
	if createError != nil {
		ch.Logger.Error(
			fmt.Sprintf("failed to create new cash record in database. %s", createError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	ch.Logger.Debug("created a new cash record in database", requestId)

	return c.JSON(http.StatusOK, LooseJson{"success": true})
}

func (ch *CashHandler) UpdateCashRecord(c echo.Context) error {
	data := new(UpdateCashRecordRequestBody)
	clientId := c.Get("uid").(string)
	requestId := zap.String("requestId", c.Get("requestId").(string))
	ch.Logger.Info("starts", requestId)

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
		ch.Logger.Error(fmt.Sprintf("invalid field. %s", validateError.Error()), requestId)
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Invalid field"},
		)
	}

	// Update account in database
	_, updateError := database.UpdateCashRecord(
		ch.Db,
		database.UpdateCashRecordParams{
			Id:         data.Id,
			Amount:     data.Amount,
			ClientId:   clientId,
			CurrencyId: data.CurrencyId,
		},
	)
	if updateError != nil {
		ch.Logger.Error(
			fmt.Sprintf("failed to update cash in database. %s", updateError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	ch.Logger.Debug("updated cash in database", requestId)

	return c.JSON(http.StatusOK, LooseJson{"success": true})
}

func (ch *CashHandler) DeleteCash(c echo.Context) error {
	clientId := c.Get("uid").(string)
	requestId := zap.String("requestId", c.Get("requestId").(string))
	ch.Logger.Info("starts", requestId)

	cashId := c.QueryParam("id")
	if len(cashId) == 0 {
		ch.Logger.Error("undefined cash id", requestId)
		return c.JSON(
			http.StatusBadRequest,
			LooseJson{"success": false, "error": "Undefined cash record id."},
		)
	}

	// Delete cash record in database
	_, deleteError := database.DeleteCashRecord(ch.Db, cashId, clientId)
	if deleteError != nil {
		ch.Logger.Error(
			fmt.Sprintf("failed to delete cash record in database. %s", deleteError.Error()),
			requestId,
		)
		return c.JSON(
			http.StatusInternalServerError,
			LooseJson{"success": false, "error": "Internal server error."},
		)
	}
	ch.Logger.Debug("deleted cash record in database", requestId)

	return c.JSON(http.StatusOK, LooseJson{"success": true})
}
