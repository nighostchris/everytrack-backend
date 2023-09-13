package exrate

import (
	"github.com/labstack/echo/v4"
)

func (erh *ExchangeRateHandler) BindRoutes(g *echo.Group) {
	g.GET("/", erh.GetAllExchangeRates)
}
