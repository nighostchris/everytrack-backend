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
		whitelistPaths := []string{"/", "/v1/auth/login", "/v1/auth/logout", "/v1/auth/refresh", "/v1/auth/signup"}

		// Ignore authentication check if the request path is whitelisted
		if slices.Contains(whitelistPaths, c.Request().RequestURI) {
			next(c)
			return nil
		}
		am.Logger.Info("going through auth middleware")
		accessToken, getTokenError := c.Request().Cookie("token")

		// Deal with case where access token does not exist in cookie
		if getTokenError != nil {
			am.Logger.Error("access token does not exist in cookie")
			// Try to extract bearer access token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if len(authHeader) == 0 {
				am.Logger.Error("access token does not exist in authorization header as well")
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success": false})
			}

			// Try to extract access token from bearer access token
			regexExpression := "\\s|Bearer"
			regex := regexp.MustCompile(regexExpression)
			bearerToken := regex.ReplaceAllString(authHeader, "")
			if len(bearerToken) == 0 {
				am.Logger.Error("access token does not exist in authorization header as well")
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success": false})
			}

			// Verify access token
			isAccessTokenValid, uid := am.TokenUtils.VerifyToken(bearerToken, 0)
			if !isAccessTokenValid {
				am.Logger.Error("invalid access token")
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success": false})
			}

			c.Set("uid", uid)
			next(c)
			return nil
		}

		// Verify the access token in cookie
		isTokenValid, uid := am.TokenUtils.VerifyToken(accessToken.Value, 0)

		// Allow request to go through if access token is valid
		if isTokenValid {
			c.Set("uid", uid)
			next(c)
			return nil
		}
		am.Logger.Error("invalid access token")

		return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success": false})
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
			requestId := c.Request().Header.Get(echo.HeaderXRequestID)
			c.Set("requestId", requestId)

			lm.Logger.Info(
				"incoming request",
				zap.String("requestId", requestId),
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
