package server

import (
	"bytes"
	"io"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/nighostchris/everytrack-backend/internal/utils"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
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
		whitelistPaths := []string{"/v1/auth/login", "/v1/auth/signup"}

		if !slices.Contains(whitelistPaths, c.Request().RequestURI) {
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
		}

		next(c)
		return nil
	}
}

func (lm *LogMiddleware) New(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		blacklistPaths := []string{"/v1/auth/login", "/v1/auth/signup"}

		if !slices.Contains(blacklistPaths, c.Request().RequestURI) {
			var data string = ""
			if c.Request().Header.Get("Content-Type") == echo.MIMEApplicationJSON {
				rawBody, _ := io.ReadAll(c.Request().Body)
				regexExpression := "\\s|\n|\t"
				regex := regexp.MustCompile(regexExpression)
				data = regex.ReplaceAllString(string(rawBody), "")
				c.Request().Body = io.NopCloser(bytes.NewBuffer(rawBody))
			}

			lm.Logger.Info(
				"incoming request",
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().RequestURI),
				zap.String("data", data),
			)
		}

		next(c)

		lm.Logger.Info(
			"request finished",
			zap.Int("status", c.Response().Status),
		)

		return nil
	}
}
