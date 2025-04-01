package controllers

import (
	"auth-service/config"
	"auth-service/repository"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	Repo *repository.AuthRepository
}

// func NewAuthController(repo *repository.AuthRepository) *AuthController {
// 	return &AuthController{Repo: repo}
// }

func (a *AuthController) Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	user, err := a.Repo.GetUserByEmail(input.Email)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Usuario no encontrado"})
	}

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Credenciales incorrectas"})
	}

	// Generar JWT
	token, err := config.GenerateJWT(user.UserID.String(), user.Email, user.UserType)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "No se pudo generar el token"})
	}

	return c.JSON(fiber.Map{"token": token})
}

// //Logout
// func (a *AuthController) Logout(c *fiber.Ctx) error {
// 	authHeader := c.Get("Authorization")
// 	if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Token no proporcionado"})
// 	}

// 	// Extraer el token
// 	tokenString := authHeader[7:]

// 	// Invalidarlo
// 	config.BlacklistToken(tokenString)

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logout exitoso, token invalidado"})
// }

func (a *AuthController) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Token no proporcionado"})
	}

	// Extraer el token
	tokenString := authHeader[7:]
	println("Intentando invalidar token:", tokenString) // Log para ver el token recibido

	// Invalidarlo
	config.BlacklistToken(tokenString)
	println("Token agregado a la lista negra")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logout exitoso, token invalidado"})
}

