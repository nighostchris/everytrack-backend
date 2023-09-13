package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nighostchris/everytrack-backend/internal/config"
	"github.com/nighostchris/everytrack-backend/internal/utils"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if error := cv.validator.Struct(i); error != nil {
		return error
	}
	return nil
}

func New(domainWhitelist []string, logger *zap.Logger, env *config.Config) *echo.Echo {
	logger.Info("initializing web server")

	e := echo.New()
	// Hiding framework promotional banner
	e.HideBanner = true
	// Bind go-playground validator to the server
	e.Validator = &CustomValidator{validator: validator.New()}
	// Middleware - CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			// TO Fix: origin is https://localhost:3000, need to update function below to handle without wildcard
			if slices.Contains(domainWhitelist, "*") || slices.Contains(domainWhitelist, origin) {
				return true, nil
			} else {
				return false, nil
			}
		},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowCredentials: true,
	}))
	// Middleware - Auth
	authMiddleware := AuthMiddleware{Logger: logger, TokenUtils: &utils.TokenUtils{Env: env, Logger: logger}}
	e.Use(authMiddleware.New)
	// Middleware - Log
	logMiddleware := LogMiddleware{Logger: logger}
	e.Use(logMiddleware.New)

	return e
}
