package config

import (
	"errors"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Estructura de los claims
type JWTClaims struct {
	UserID   string `json:"userId"`
	Email    string `json:"email"`
	UserType string `json:"userType"`
	jwt.RegisteredClaims
}

var JWTSecret = []byte(os.Getenv("JWT_SECRET"))

// Lista negra de tokens (para invalidarlos)
var tokenBlacklist = make(map[string]time.Time)
//var blacklistMutex sync.Mutex

// Generar token JWT
func GenerateJWT(userID string, email string, userType string) (string, error) {
	claims := jwt.MapClaims{
		"userId":   userID,
		"email":    email,
		"userType": userType,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// Invalidar un token (añadirlo a la lista negra)
func BlacklistToken(tokenString string) {
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()

	// Extraer los claims para obtener la expiración
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		println("Error al validar token:", err.Error()) // Log del error
		return
	}

	tokenBlacklist[tokenString] = claims.ExpiresAt.Time
	println("Token invalidado correctamente")
}


// Validar token JWT
// func ValidateJWT(tokenString string) (*JWTClaims, error) {
// 	// Verificar si el token está en la lista negra
// 	blacklistMutex.Lock()
// 	expiration, found := tokenBlacklist[tokenString]
// 	blacklistMutex.Unlock()
// 	if found && expiration.After(time.Now()) {
// 		return nil, errors.New("token invalidado")
// 	}

// 	// Parsear y validar el token
// 	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		return JWTSecret, nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	claims, ok := token.Claims.(*JWTClaims)
// 	if !ok || !token.Valid {
// 		return nil, errors.New("token inválido")
// 	}

// 	// Verificar expiración
// 	if claims.ExpiresAt.Time.Before(time.Now()) {
// 		return nil, errors.New("token expirado")
// 	}

// 	return claims, nil
// }

var blacklistMutex sync.RWMutex // RWMutex permite lectura concurrente sin bloquear

func ValidateJWT(tokenString string) (*JWTClaims, error) {
    println("Validando token:", tokenString) // Log

    // Verificar si el token está en la lista negra
    blacklistMutex.RLock() // Usamos RLock() para lectura sin bloquear escrituras
    expiration, found := tokenBlacklist[tokenString]
    blacklistMutex.RUnlock()

    if found {
        println("El token está en la lista negra") // Log
        if expiration.After(time.Now()) {
            println("Token ya invalidado") // Log
            return nil, errors.New("token invalidado")
        }
    }

    println("Intentando parsear el token") // Log
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return JWTSecret, nil
    })

    if err != nil {
        println("Error al parsear el token:", err.Error()) // Log
        return nil, err
    }

    claims, ok := token.Claims.(*JWTClaims)
    if !ok || !token.Valid {
        println("Token inválido") // Log
        return nil, errors.New("token inválido")
    }

    println("Token parseado correctamente") // Log
    if claims.ExpiresAt.Time.Before(time.Now()) {
        println("Token expirado") // Log
        return nil, errors.New("token expirado")
    }

    println("Token válido") // Log
    return claims, nil
}
