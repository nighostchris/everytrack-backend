package handlers

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nighostchris/everytrack-backend/internal/config"
	"github.com/nighostchris/everytrack-backend/internal/utils"
	"go.uber.org/zap"
)

type Handlers struct {
	Auth         *AuthHandler
	Savings      *SavingsHandler
	Settings     *SettingsHandler
	Currencies   *CurrenciesHandler
	ExchangeRate *ExchangeRateHandler
}

func Init(db *pgxpool.Pool, env *config.Config, logger *zap.Logger) *Handlers {
	return &Handlers{
		Savings:      &SavingsHandler{Db: db, Logger: logger},
		Settings:     &SettingsHandler{Db: db, Logger: logger},
		Currencies:   &CurrenciesHandler{Db: db, Logger: logger},
		ExchangeRate: &ExchangeRateHandler{Db: db, Logger: logger},
		Auth:         &AuthHandler{Db: db, Logger: logger, TokenUtils: &utils.TokenUtils{Env: env, Logger: logger}},
	}
}

func (h *Handlers) BindRoutes(e *echo.Echo) {
	v1 := e.Group("/v1")
	// ============================================================
	// /v1/auth endpoints
	// ============================================================
	auth := v1.Group("/auth")
	auth.POST("/login", h.Auth.Login)
	auth.POST("/signup", h.Auth.Signup)
	auth.POST("/verify", h.Auth.Verify)
	// ============================================================
	// /v1/currency endpoints
	// ============================================================
	currency := v1.Group("/currencies")
	currency.GET("", h.Currencies.GetAllCurrencies)
	// ============================================================
	// /v1/exrates endpoints
	// ============================================================
	exchangeRates := v1.Group("/exrates")
	exchangeRates.GET("", h.ExchangeRate.GetAllExchangeRates)
	// ============================================================
	// /v1/savings endpoints
	// ============================================================
	savings := v1.Group("/savings")
	savings.GET("", h.Savings.GetAllBankDetails)
	savings.PUT("/account", h.Savings.UpdateAccount)
	savings.POST("/account", h.Savings.CreateNewAccount)
	savings.GET("/account", h.Savings.GetAllBankAccounts)
	// ============================================================
	// /v1/settings endpoints
	// ============================================================
	settings := v1.Group("/settings")
	settings.PUT("", h.Settings.UpdateSettings)
	settings.GET("", h.Settings.GetAllClientSettings)
}
