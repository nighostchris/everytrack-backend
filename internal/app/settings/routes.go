package settings

import (
	"github.com/labstack/echo/v4"
)

func (sh *SettingsHandler) BindRoutes(g *echo.Group) {
	g.GET("/", sh.GetAllClientSettings, sh.AuthMiddleware.New)
	g.PUT("/", sh.UpdateSettings, sh.AuthMiddleware.New)
}
