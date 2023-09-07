package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nighostchris/everytrack-backend/internal/utils"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	Logger     *zap.Logger
	TokenUtils *utils.TokenUtils
}

type LogMiddleware struct {
	Logger *zap.Logger
}

func (am *AuthMiddleware) New(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		am.Logger.Info("going through auth middleware")
		token, getTokenError := c.Request().Cookie("token")

		if getTokenError != nil {
			am.Logger.Error("token does not exist in cookie")
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success": false})
		}

		isTokenValid, uid := am.TokenUtils.VerifyToken(token.Value)

		if !isTokenValid {
			am.Logger.Error("invalid token")
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success": false})
		}

		c.Set("uid", uid)
		next(c)
		return nil
	}
}

func (lm *LogMiddleware) New(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		lm.Logger.Info(
			"incoming request",
			zap.String("method", c.Request().Method),
			zap.String("path", c.Request().RequestURI),
		)

		next(c)

		lm.Logger.Info(
			"request finished",
			zap.Int("status", c.Response().Status),
		)

		return nil
	}
}