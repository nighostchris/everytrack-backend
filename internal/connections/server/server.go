package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func New(domainWhitelist []string, logger *zap.Logger) *echo.Echo {
	logger.Info("initializing web server")

	e := echo.New()
	// Hiding framework promotional banner
	e.HideBanner = true
	// Bind go-playground validator to the server
	e.Validator = &CustomValidator{validator: validator.New()}
	// CORS setting
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			// TO Fix: origin is https://localhost:3000, need to update function below to handle without wildcard
			if slices.Contains(domainWhitelist, "*") || slices.Contains(domainWhitelist, origin) {
				return true, nil
			} else {
				return false, nil
			}
		},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowCredentials: true,
	}))
	// Logger setting
	logMiddleware := LogMiddleware{Logger: logger}
	e.Use(logMiddleware.New)

	return e
}
