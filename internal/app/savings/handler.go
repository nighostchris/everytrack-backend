package savings

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nighostchris/everytrack-backend/internal/connections/server"
	"github.com/nighostchris/everytrack-backend/internal/database"
	"go.uber.org/zap"
)

type SavingsHandler struct {
	Db             *pgxpool.Pool
	Logger         *zap.Logger
	AuthMiddleware *server.AuthMiddleware
}

type AccountType struct {
	Id   string
	Name string
}

type GetAllBankDetailsResponseData struct {
	Name        string
	AccountType []AccountType
}

func NewHandler(db *pgxpool.Pool, l *zap.Logger, am *server.AuthMiddleware) *SavingsHandler {
	handler := SavingsHandler{Db: db, Logger: l, AuthMiddleware: am}
	return &handler
}

func (sh *SavingsHandler) GetAllBankDetails(c echo.Context) error {
	sh.Logger.Info("starts")

	// Get all bank details from database
	bankDetails, getBankDetailsError := database.GetAllBankDetails(sh.Db)

	if getBankDetailsError != nil {
		sh.Logger.Error(fmt.Sprintf("failed to get bank details from database. %s", getBankDetailsError.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Internal server error."})
	}
	sh.Logger.Debug("got bank details from database")

	// Construct the response object
	bankDetailsMap := make(map[string][]AccountType)
	for _, bankDetail := range bankDetails {
		accountType := AccountType{Id: bankDetail.AccountTypeId, Name: bankDetail.AccountTypeName}
		bankDetailsMap[bankDetail.Name] = append(bankDetailsMap[bankDetail.Name], accountType)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"success": true, "data": bankDetailsMap})
}
