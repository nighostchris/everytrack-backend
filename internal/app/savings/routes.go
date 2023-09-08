package savings

import (
	"github.com/labstack/echo/v4"
	"github.com/nighostchris/everytrack-backend/internal/connections/server"
)

func (sh *SavingsHandler) BindRoutes(g *echo.Group) {
	middleware := server.AuthMiddleware{Logger: sh.Logger}
	g.GET("/", sh.GetAllBankDetails, middleware.New)
}
