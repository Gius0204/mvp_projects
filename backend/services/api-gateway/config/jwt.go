package config

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// Clave secreta JWT
var JWTSecret = []byte(os.Getenv("JWT_SECRET"))

// Estructura de los claims del token
type JWTClaims struct {
	UserID   string `json:"userId"`
	Email    string `json:"email"`
	UserType string `json:"userType"`
	jwt.RegisteredClaims
}

// Validar JWT
func ValidateJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token inválido")
	}

	return claims, nil
}
