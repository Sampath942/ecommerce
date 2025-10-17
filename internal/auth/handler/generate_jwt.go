package handler

import (
	"time"

	"github.com/Sampath942/ecommerce/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID  int  `json:"user_id"`
	IsAdmin bool `json:"is_admin"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(uid int, isAdmin bool) (string, error) {
	mySigningKey := []byte(config.AppConfig.JWTSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
			Audience:  []string{"user-service"},
		},
	})
	ss, err := token.SignedString(mySigningKey)
	return ss, err
}
