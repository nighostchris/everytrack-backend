package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
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
		AllowOrigins:     domainWhitelist,
		AllowCredentials: true,
	}))
	// Logger setting
	logMiddleware := LogMiddleware{Logger: logger}
	e.Use(logMiddleware.New)

	return e
}
