package cron

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nighostchris/everytrack-backend/internal/config"
	"github.com/nighostchris/everytrack-backend/internal/tools"
	"go.uber.org/zap"
)

type CronJob struct {
	Logger *zap.Logger
	Db     *pgxpool.Pool
	Env    *config.Config
}

func Init(db *pgxpool.Pool, env *config.Config, logger *zap.Logger) *CronJob {
	return &CronJob{Db: db, Env: env, Logger: logger}
}

// Fetch exchange rates from Github Currency API every day
func (cj *CronJob) SubscribeExchangeRates() {
	tool := tools.GithubCurrencyApi{Db: cj.Db, Logger: cj.Logger}
	go func() {
		for {
			tool.FetchLatestExchangeRates()
			time.Sleep(24 * time.Hour)
		}
	}()
}

func (cj *CronJob) SubscribeTwelveDataFinancialData() {
	tool := tools.TwelveDataFinancialDataApi{Db: cj.Db, Logger: cj.Logger, Env: cj.Env}
	go func() {
		for {
			tool.FetchLatestUSStockPrice()
			time.Sleep(30 * time.Minute)
		}
	}()
}
