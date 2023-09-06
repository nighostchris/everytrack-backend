package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nighostchris/everytrack-backend/internal/config"
)

type TokenUtils struct {
	Env *config.Config
}

type CustomClaims struct {
	jwt.RegisteredClaims
}

func (tu TokenUtils) GenerateToken(sub string, tokenType int) (string, error) {
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
