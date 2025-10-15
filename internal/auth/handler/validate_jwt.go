package handler

import (
	"errors"

	"github.com/Sampath942/ecommerce/config"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWTToken(tokenString string) (map[string]any, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(config.AppConfig.JWTSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return map[string]any{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		return map[string]any{}, errors.New("extracting claims failed")
	}
}
