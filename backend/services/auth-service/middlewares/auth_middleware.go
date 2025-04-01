package middlewares

import (
	"auth-service/config"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Middleware para proteger rutas con JWT
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	// Verificar que el header Authorization esté presente y tenga el formato correcto
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token no proporcionado"})
	}

	// Extraer el token eliminando "Bearer "
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Validar el token
	claims, err := config.ValidateJWT(tokenString)
	if err != nil {
		fmt.Println("Error en validación JWT:", err.Error()) // 🔍 Para depuración
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token inválido o expirado"})
	}

	// Guardar los datos del usuario en el contexto para que estén accesibles en los controladores
	c.Locals("userID", claims.UserID)
	c.Locals("email", claims.Email)
	c.Locals("userType", claims.UserType)

	// Continuar con la solicitud
	return c.Next()
}
