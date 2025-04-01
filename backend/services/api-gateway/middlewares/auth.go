package middleware

import (
	"api-gateway/config"

	"github.com/gofiber/fiber/v2"
)

// Middleware para validar JWT
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token no proporcionado"})
	}

	// Extraer token
	tokenString := authHeader[7:]

	// Validar token
	claims, err := config.ValidateJWT(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token inválido"})
	}

	// Guardar información del usuario en el contexto
	c.Locals("userID", claims.UserID)
	c.Locals("email", claims.Email)
	c.Locals("userType", claims.UserType)

	return c.Next()
}
