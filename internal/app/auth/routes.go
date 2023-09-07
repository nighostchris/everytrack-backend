package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/nighostchris/everytrack-backend/internal/connections/server"
)

func (ah *AuthHandler) BindRoutes(g *echo.Group) {
	middleware := server.AuthMiddleware{Logger: ah.Logger, TokenUtils: ah.TokenUtils}
	g.POST("/login", ah.Login)
	g.POST("/signup", ah.Signup)
	g.POST("/verify", ah.Verify, middleware.New)
	// auth.POST("/refresh")
}
