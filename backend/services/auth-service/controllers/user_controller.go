package controllers

import (
	"github.com/gofiber/fiber/v2"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

// Obtener perfil del usuario autenticado
func (uc *UserController) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	email := c.Locals("email").(string)
	userType := c.Locals("userType").(string)

	return c.JSON(fiber.Map{
		//"mensaje": "Hola aqui en user profile",
		"userId":   userID,
		"email":    email,
		"userType": userType,
	})
}
