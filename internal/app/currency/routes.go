package currency

import (
	"github.com/labstack/echo/v4"
)

func (ch *CurrencyHandler) BindRoutes(g *echo.Group) {
	g.GET("/", ch.GetAllCurrencies, ch.AuthMiddleware.New)
}
