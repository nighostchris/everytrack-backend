package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/iancoleman/strcase"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nighostchris/everytrack-backend/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Db *pgxpool.Pool
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

func NewHandler(db *pgxpool.Pool) *AuthHandler {
	handler := AuthHandler{Db: db}
	return &handler
}

func (ah *AuthHandler) Signup(c echo.Context) error {
	data := new(SignupRequestBody)

	// Retrieve request body and validate with schema
	if bindError := c.Bind(data); bindError != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Missing required fields."})
	}

	if validateError := c.Validate(data); validateError != nil {
		var ve validator.ValidationErrors
		if errors.As(validateError, &ve) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": fmt.Sprintf("Invalid field %s.", strcase.ToLowerCamel(ve[0].Field()))})
		}
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": validateError.Error()})
	}

	// Check if client exists already in database
	_, getClientError := database.GetClientByEmail(ah.Db, data.Email)

	if getClientError == nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"success": false, "error": "Username already in use."})
	}

	passwordHash, generatePasswordHashError := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

	if generatePasswordHashError != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Internal server error."})
	}

	newClientId, createNewClientError := database.CreateNewClient(
		ah.Db,
		database.CreateNewClientParams{
			Email:    data.Email,
			Username: data.Username,
			Password: string(passwordHash),
		},
	)

	if createNewClientError != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Internal server error."})
	}

	// Construct access token
	expiryTime := time.Now().Add(time.Hour * time.Duration(24))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "everytrack-backend",
		"sub": newClientId,
		"exp": expiryTime.Unix(),
	})

	signedToken, signError := token.SignedString([]byte("to-be-replace-by-secret-later"))

	if signError != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": signError.Error()})
	}

	// Set access token into cookie
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    signedToken,
		Expires:  expiryTime,
		Path:     "/",
		Domain:   "localhost",
		Secure:   true,                  // Forbid cookie from transmitting over simple HTTP
		HttpOnly: true,                  // Blocks access of related cookie from client side
		SameSite: http.SameSiteNoneMode, // SameSite 'none' has to be used together with secure - true
	})

	return c.JSON(http.StatusOK, map[string]interface{}{"success": true})
}

func (ah *AuthHandler) Login(c echo.Context) error {
	data := new(LoginRequestBody)

	// Retrieve request body and validate with schema
	if bindError := c.Bind(data); bindError != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Missing required fields."})
	}

	if validateError := c.Validate(data); validateError != nil {
		var ve validator.ValidationErrors
		if errors.As(validateError, &ve) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": fmt.Sprintf("Invalid field %s.", strcase.ToLowerCamel(ve[0].Field()))})
		}
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": validateError.Error()})
	}

	// Try to get client from database by input email
	client, getClientError := database.GetClientByEmail(ah.Db, data.Email)

	if getClientError != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"success": false, "error": getClientError.Error()})
	}

	// Verify password
	verifyPasswordError := bcrypt.CompareHashAndPassword([]byte(client.Password), []byte(data.Password))

	if verifyPasswordError != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"success": false, "error": verifyPasswordError.Error()})
	}

	// Construct access token
	expiryTime := time.Now().Add(time.Hour * time.Duration(24))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "everytrack-backend",
		"sub": client.Id,
		"exp": expiryTime.Unix(),
	})

	signedToken, signError := token.SignedString([]byte("to-be-replace-by-secret-later"))

	if signError != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": signError.Error()})
	}

	// Set access token into cookie
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    signedToken,
		Expires:  expiryTime,
		Path:     "/",
		Domain:   "localhost",
		Secure:   true,                  // Forbid cookie from transmitting over simple HTTP
		HttpOnly: true,                  // Blocks access of related cookie from client side
		SameSite: http.SameSiteNoneMode, // SameSite 'none' has to be used together with secure - true
	})

	return c.JSON(http.StatusOK, map[string]interface{}{"success": true})
}
