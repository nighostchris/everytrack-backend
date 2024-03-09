package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nighostchris/everytrack-backend/internal/database"
	"github.com/nighostchris/everytrack-backend/internal/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Db         *pgxpool.Pool
	Logger     *zap.Logger
	TokenUtils *utils.TokenUtils
}

type SignupRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,max=20"`
	Password string `json:"password" validate:"required"`
}

type LoginRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (ah *AuthHandler) Signup(c echo.Context) error {
	data := new(SignupRequestBody)

	// Retrieve request body and validate with schema
	if bindError := c.Bind(data); bindError != nil {
		return c.JSON(http.StatusBadRequest, LooseJson{"success": false, "error": "Missing required fields."})
	}

	if validateError := c.Validate(data); validateError != nil {
		var ve validator.ValidationErrors
		if errors.As(validateError, &ve) {
			return c.JSON(http.StatusBadRequest, LooseJson{"success": false, "error": fmt.Sprintf("Invalid field %s.", strcase.ToLowerCamel(ve[0].Field()))})
		}
		return c.JSON(http.StatusBadRequest, LooseJson{"success": false, "error": validateError.Error()})
	}

	// Check if client exists already in database
	_, getClientError := database.GetClientByEmail(ah.Db, data.Email)

	if getClientError == nil {
		return c.JSON(http.StatusNotFound, LooseJson{"success": false, "error": "Username already in use."})
	}

	// Get default currency from database
	defaultCurrencyId, getDefaultCurrencyIdError := database.GetDefaultCurrency(ah.Db)
	if getDefaultCurrencyIdError != nil {
		ah.Logger.Error(fmt.Sprintf("failed to get default currency id. %s", getDefaultCurrencyIdError.Error()))
		return c.JSON(http.StatusNotFound, LooseJson{"success": false, "error": "Internal server error."})
	}

	passwordHash, generatePasswordHashError := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

	if generatePasswordHashError != nil {
		return c.JSON(http.StatusInternalServerError, LooseJson{"success": false, "error": "Internal server error."})
	}

	newClientId, createNewClientError := database.CreateNewClient(
		ah.Db,
		database.CreateNewClientParams{
			Email:      data.Email,
			Username:   data.Username,
			Password:   string(passwordHash),
			CurrencyId: defaultCurrencyId,
		},
	)

	if createNewClientError != nil {
		return c.JSON(http.StatusInternalServerError, LooseJson{"success": false, "error": "Internal server error."})
	}

	// Construct access token
	accessToken, generateAccessTokenError := ah.TokenUtils.GenerateToken(newClientId, 0)
	if generateAccessTokenError != nil {
		return c.JSON(http.StatusInternalServerError, LooseJson{"success": false, "error": generateAccessTokenError.Error()})
	}

	// Set access token into cookie
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    accessToken,
		Expires:  ah.TokenUtils.GetAccessTokenExpiry(),
		Path:     "/",
		Secure:   true,                  // Forbid cookie from transmitting over simple HTTP
		HttpOnly: true,                  // Blocks access of related cookie from client side
		SameSite: http.SameSiteNoneMode, // SameSite 'none' has to be used together with secure - true
	})

	// Construct refresh token
	refreshToken, generateRefreshTokenError := ah.TokenUtils.GenerateToken(newClientId, 1)
	if generateRefreshTokenError != nil {
		return c.JSON(http.StatusInternalServerError, LooseJson{"success": false, "error": generateRefreshTokenError.Error()})
	}
	ah.Logger.Debug("generated refresh token")

	return c.JSON(http.StatusOK, LooseJson{"success": true, "data": LooseJson{"token": accessToken, "refresh": refreshToken}})
}

func (ah *AuthHandler) Login(c echo.Context) error {
	data := new(LoginRequestBody)
	ah.Logger.Info("starts")

	// Retrieve request body and validate with schema
	if bindError := c.Bind(data); bindError != nil {
		return c.JSON(http.StatusBadRequest, LooseJson{"success": false, "error": "Missing required fields"})
	}

	if validateError := c.Validate(data); validateError != nil {
		var ve validator.ValidationErrors
		if errors.As(validateError, &ve) {
			return c.JSON(http.StatusBadRequest, LooseJson{"success": false, "error": fmt.Sprintf("Invalid field %s", strcase.ToLowerCamel(ve[0].Field()))})
		}
		ah.Logger.Error(fmt.Sprintf("invalid field. %s", validateError.Error()))
		return c.JSON(http.StatusBadRequest, LooseJson{"success": false, "error": "Invalid field"})
	}
	ah.Logger.Debug("validated request parameters")

	// Try to get client from database by input email
	client, getClientError := database.GetClientByEmail(ah.Db, data.Email)
	if getClientError != nil {
		ah.Logger.Error(fmt.Sprintf("failed to get client from database by email - %s. %s", data.Email, getClientError.Error()))
		return c.JSON(http.StatusNotFound, LooseJson{"success": false, "error": "Invalid user"})
	}
	ah.Logger.Debug(fmt.Sprintf("got client from database with email - %s", client.Email))

	// Verify password
	verifyPasswordError := bcrypt.CompareHashAndPassword([]byte(client.Password), []byte(data.Password))
	if verifyPasswordError != nil {
		ah.Logger.Error(fmt.Sprintf("password verification failed. %s", verifyPasswordError.Error()))
		return c.JSON(http.StatusNotFound, LooseJson{"success": false, "error": "Incorrect password"})
	}
	ah.Logger.Debug("verified password")

	// Construct access token
	accessToken, generateAccessTokenError := ah.TokenUtils.GenerateToken(client.Id, 0)
	if generateAccessTokenError != nil {
		ah.Logger.Error(fmt.Sprintf("access token generation failed. %s", generateAccessTokenError.Error()))
		return c.JSON(http.StatusInternalServerError, LooseJson{"success": false, "error": "Internal server error"})
	}
	ah.Logger.Debug("generated access token")

	// Set access token into cookie
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    accessToken,
		Expires:  ah.TokenUtils.GetAccessTokenExpiry(),
		Path:     "/",
		Secure:   true,                  // Forbid cookie from transmitting over simple HTTP
		HttpOnly: true,                  // Blocks access of related cookie from client side
		SameSite: http.SameSiteNoneMode, // SameSite 'none' has to be used together with secure - true
	})
	ah.Logger.Debug("finished setting access token to response cookie")

	// Construct refresh token
	refreshToken, generateRefreshTokenError := ah.TokenUtils.GenerateToken(client.Id, 1)
	if generateRefreshTokenError != nil {
		ah.Logger.Error(fmt.Sprintf("refresh token generation failed. %s", generateRefreshTokenError.Error()))
		return c.JSON(http.StatusInternalServerError, LooseJson{"success": false, "error": "Internal server error"})
	}
	ah.Logger.Debug("generated refresh token")

	return c.JSON(http.StatusOK, LooseJson{"success": true, "data": LooseJson{"token": accessToken, "refresh": refreshToken}})
}

func (ah *AuthHandler) Logout(c echo.Context) error {
	ah.Logger.Info("starts")

	// Void access token in cookie
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    "",
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
	ah.Logger.Debug("finished voiding access token in cookie")

	return c.JSON(http.StatusOK, LooseJson{"success": true})
}

func (ah *AuthHandler) Refresh(c echo.Context) error {
	ah.Logger.Info("starts")

	// Try to extract bearer refresh token from Authorization header
	authHeader := c.Request().Header.Get("Authorization")
	if len(authHeader) == 0 {
		ah.Logger.Error("refresh token does not exist in authorization header")
		return c.JSON(http.StatusUnauthorized, LooseJson{"success": false})
	}
	ah.Logger.Debug("extracted bearer refresh token from authorization header")

	// Try to extract refresh token from bearer refresh token
	regexExpression := "\\s|Bearer"
	regex := regexp.MustCompile(regexExpression)
	bearerToken := regex.ReplaceAllString(authHeader, "")
	if len(bearerToken) == 0 {
		ah.Logger.Error("refresh token does not exist in authorization header")
		return c.JSON(http.StatusUnauthorized, LooseJson{"success": false})
	}
	ah.Logger.Debug("extracted refresh token from bearer refresh token")

	// Verify refresh token
	isRefreshTokenValid, uid := ah.TokenUtils.VerifyToken(bearerToken, 1)
	if !isRefreshTokenValid {
		ah.Logger.Error("invalid refresh token")
		return c.JSON(http.StatusUnauthorized, LooseJson{"success": false})
	}

	// Construct new access token
	accessToken, generateAccessTokenError := ah.TokenUtils.GenerateToken(uid, 0)
	if generateAccessTokenError != nil {
		ah.Logger.Error(fmt.Sprintf("access token generation failed. %s", generateAccessTokenError.Error()))
		return c.JSON(http.StatusInternalServerError, LooseJson{"success": false, "error": "Internal server error"})
	}
	ah.Logger.Debug("generated new access token")

	// Set access token into cookie
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    accessToken,
		Expires:  ah.TokenUtils.GetAccessTokenExpiry(),
		Path:     "/",
		Secure:   true,                  // Forbid cookie from transmitting over simple HTTP
		HttpOnly: true,                  // Blocks access of related cookie from client side
		SameSite: http.SameSiteNoneMode, // SameSite 'none' has to be used together with secure - true
	})
	ah.Logger.Debug("finished setting access token to response cookie")

	// Construct new refresh token
	refreshToken, generateRefreshTokenError := ah.TokenUtils.GenerateToken(uid, 1)
	if generateRefreshTokenError != nil {
		ah.Logger.Error(fmt.Sprintf("refresh token generation failed. %s", generateRefreshTokenError.Error()))
		return c.JSON(http.StatusInternalServerError, LooseJson{"success": false, "error": "Internal server error"})
	}
	ah.Logger.Debug("generated new refresh token")

	return c.JSON(http.StatusOK, LooseJson{"success": true, "data": LooseJson{"token": accessToken, "refresh": refreshToken}})
}

func (ah *AuthHandler) Verify(c echo.Context) error {
	return c.JSON(http.StatusAccepted, LooseJson{"success": true})
}
