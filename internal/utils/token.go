package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nighostchris/everytrack-backend/internal/config"
	"go.uber.org/zap"
)

type TokenUtils struct {
	Env    *config.Config
	Logger *zap.Logger
}

type CustomClaims struct {
	jwt.RegisteredClaims
}

func (tu TokenUtils) GenerateToken(sub string, tokenType int) (string, error) {
	tu.Logger.Info(fmt.Sprintf("starts generating token of type %d for %s", tokenType, sub))
	// Define which secret key to sign token
	var secret []byte
	if tokenType == 0 {
		secret = []byte(tu.Env.AccessTokenSecret)
	} else {
		return "", errors.New("invalid token type for generation")
	}

	// Construct claims as JWT payload
	claims := new(CustomClaims)
	claims.Issuer = "everytrack-backend"
	claims.Subject = sub
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(tu.Env.TokenExpiryInHour)))
	claims.NotBefore = jwt.NewNumericDate(time.Now())
	claims.IssuedAt = jwt.NewNumericDate(time.Now())

	// Construct JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign JWT
	signedJwt, signError := token.SignedString(secret)

	if signError != nil {
		return "", errors.New("token generation failed")
	}

	return signedJwt, nil
}
