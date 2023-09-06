package auth

import "github.com/labstack/echo/v4"

func (ah *AuthHandler) BindRoutes(g *echo.Group) {
	g.POST("/login", ah.Login)
	g.POST("/signup", ah.Signup)
	// auth.POST("/refresh")
	// auth.POST("/verify")
}
