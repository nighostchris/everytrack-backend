package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nighostchris/everytrack-backend/internal/config"
	"github.com/nighostchris/everytrack-backend/internal/utils"
	"go.uber.org/zap"
)

type Handlers struct {
	Auth           *AuthHandler
	Cash           *CashHandler
	Stocks         *StocksHandler
	Accounts       *AccountsHandler
	Settings       *SettingsHandler
	Providers      *ProvidersHandler
	Countries      *CountriesHandler
	Currencies     *CurrenciesHandler
	Transactions   *TransactionsHandler
	ExchangeRates  *ExchangeRatesHandler
	FuturePayments *FuturePaymentsHandler
}

type LooseJson map[string]interface{}

func Init(db *pgxpool.Pool, env *config.Config, logger *zap.Logger) *Handlers {
	return &Handlers{
		Cash:           &CashHandler{Db: db, Logger: logger},
		Stocks:         &StocksHandler{Db: db, Logger: logger},
		Settings:       &SettingsHandler{Db: db, Logger: logger},
		Accounts:       &AccountsHandler{Db: db, Logger: logger},
		Providers:      &ProvidersHandler{Db: db, Logger: logger},
		Countries:      &CountriesHandler{Db: db, Logger: logger},
		Currencies:     &CurrenciesHandler{Db: db, Logger: logger},
		Transactions:   &TransactionsHandler{Db: db, Logger: logger},
		ExchangeRates:  &ExchangeRatesHandler{Db: db, Logger: logger},
		FuturePayments: &FuturePaymentsHandler{Db: db, Logger: logger},
		Auth:           &AuthHandler{Db: db, Logger: logger, TokenUtils: &utils.TokenUtils{Env: env, Logger: logger}},
	}
}

func (h *Handlers) BindRoutes(e *echo.Echo) {
	// ============================================================
	// / - Healthcheck endpoint
	// ============================================================
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{"success": true})
	})

	v1 := e.Group("/v1")
	// ============================================================
	// /v1/accounts endpoints
	// ============================================================
	accounts := v1.Group("/accounts")
	accounts.PUT("", h.Accounts.UpdateAccount)
	accounts.DELETE("", h.Accounts.DeleteAccount)
	accounts.POST("", h.Accounts.CreateNewAccount)
	accounts.POST("/transfer", h.Accounts.TransferBetweenAccounts)
	accounts.GET("", h.Accounts.GetAllAccountsByType)
	// ============================================================
	// /v1/auth endpoints
	// ============================================================
	auth := v1.Group("/auth")
	auth.POST("/login", h.Auth.Login)
	auth.POST("/signup", h.Auth.Signup)
	auth.POST("/verify", h.Auth.Verify)
	auth.POST("/logout", h.Auth.Logout)
	auth.POST("/refresh", h.Auth.Refresh)
	// ============================================================
	// /v1/cash endpoints
	// ============================================================
	cash := v1.Group("/cash")
	cash.GET("", h.Cash.GetAllCash)
	cash.DELETE("", h.Cash.DeleteCash)
	cash.PUT("", h.Cash.UpdateCashRecord)
	cash.POST("", h.Cash.CreateNewCashRecord)
	// ============================================================
	// /v1/countries endpoints
	// ============================================================
	countries := v1.Group("/countries")
	countries.GET("", h.Countries.GetAllCountries)
	// ============================================================
	// /v1/currency endpoints
	// ============================================================
	currencies := v1.Group("/currencies")
	currencies.GET("", h.Currencies.GetAllCurrencies)
	// ============================================================
	// /v1/transactions endpoints
	// ============================================================
	transactions := v1.Group("/transactions")
	transactions.GET("", h.Transactions.GetAllTransactions)
	transactions.DELETE("", h.Transactions.DeleteTransaction)
	transactions.POST("", h.Transactions.CreateNewTransaction)
	// ============================================================
	// /v1/exrates endpoints
	// ============================================================
	exchangeRates := v1.Group("/exrates")
	exchangeRates.GET("", h.ExchangeRates.GetAllExchangeRates)
	// ============================================================
	// /v1/fpayments endpoints
	// ============================================================
	futurePayments := v1.Group("/fpayments")
	futurePayments.PUT("", h.FuturePayments.UpdateFuturePayment)
	futurePayments.GET("", h.FuturePayments.GetAllFuturePayments)
	futurePayments.DELETE("", h.FuturePayments.DeleteFuturePayment)
	futurePayments.POST("", h.FuturePayments.CreateNewFuturePayment)
	// ============================================================
	// /v1/providers endpoints
	// ============================================================
	providers := v1.Group("/providers")
	providers.GET("", h.Providers.GetAllProvidersByType)
	// ============================================================
	// /v1/settings endpoints
	// ============================================================
	settings := v1.Group("/settings")
	settings.PUT("", h.Settings.UpdateSettings)
	settings.GET("", h.Settings.GetAllClientSettings)
	// ============================================================
	// /v1/stocks endpoints
	// ============================================================
	stocks := v1.Group("/stocks")
	stocks.GET("", h.Stocks.GetAllStocks)
	stocks.PUT("/holdings", h.Stocks.UpdateStockHolding)
	stocks.GET("/holdings", h.Stocks.GetAllStockHoldings)
	stocks.DELETE("/holdings", h.Stocks.DeleteStockHolding)
	stocks.POST("/holdings", h.Stocks.CreateNewStockHolding)
}
