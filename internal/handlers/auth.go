package handlers

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

type LoginRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (ah *AuthHandler) Login(c echo.Context) error {
	data := new(LoginRequestBody)

	if bindError := c.Bind(data); bindError != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Missing required fields."})
	}

	if validateError := c.Validate(data); validateError != nil {
		var ve validator.ValidationErrors
		if errors.As(validateError, &ve) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": fmt.Sprintf("Invalid field %s", strcase.ToLowerCamel(ve[0].Field()))})
		}
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": validateError.Error()})
	}

	client, getClientError := database.GetClientByEmail(ah.Db, data.Email)

	if getClientError != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"success": false, "error": getClientError.Error()})
	}

	verifyPasswordError := bcrypt.CompareHashAndPassword([]byte(client.Password), []byte(data.Password))

	if verifyPasswordError != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"success": false, "error": verifyPasswordError.Error()})
	}

	expiryTime := time.Now().Add(time.Hour * time.Duration(24))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "everytrack-backend",
		"sub": client.Id,
		"exp": expiryTime.Unix(),
	})

	signedToken, signError := token.SignedString("to-be-replace-by-secret-later")

	if signError != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": signError.Error()})
	}

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
