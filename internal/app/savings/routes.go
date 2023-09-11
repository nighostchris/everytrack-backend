package savings

import (
	"github.com/labstack/echo/v4"
)

func (sh *SavingsHandler) BindRoutes(g *echo.Group) {
	g.GET("/", sh.GetAllBankDetails)
	g.GET("/account", sh.GetAllBankAccounts, sh.AuthMiddleware.New)
	g.POST("/account", sh.CreateNewAccount, sh.AuthMiddleware.New)
}
