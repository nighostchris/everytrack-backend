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

func (tu TokenUtils) VerifyToken(token string) (bool, string) {
	tu.Logger.Info(fmt.Sprintf("starts verifying token %s", token))

	t, parseTokenError := jwt.ParseWithClaims(token, &CustomClaims{}, func(tk *jwt.Token) (interface{}, error) {
		if _, isValidSigningMethod := tk.Method.(*jwt.SigningMethodHMAC); !isValidSigningMethod {
			return nil, fmt.Errorf("unexpected signing method %s", tk.Header["alg"])
		}
		return tu.Env.AccessTokenSecret, nil
	})

	if parseTokenError != nil || t.Valid {
		tu.Logger.Error(fmt.Sprintf("invalid token %s", token))
		return false, ""
	}

	claim := t.Claims.(*CustomClaims)
	tu.Logger.Info(fmt.Sprintf("claim in token %s - %#v", token, claim))

	currentTime := time.Now().Unix()
	if claim.ExpiresAt.Time.Unix() <= currentTime {
		tu.Logger.Error("token expired")
		return false, ""
	}
	sub := claim.Subject

	return true, sub
}
