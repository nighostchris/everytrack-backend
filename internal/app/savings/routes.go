package savings

import (
	"github.com/labstack/echo/v4"
)

func (sh *SavingsHandler) BindRoutes(g *echo.Group) {
	g.GET("/", sh.GetAllBankDetails, sh.AuthMiddleware.New)
}
