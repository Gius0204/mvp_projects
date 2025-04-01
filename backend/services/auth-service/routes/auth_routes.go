package routes

import (
	"auth-service/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App, authController *controllers.AuthController) {
	auth := app.Group("/auth")
	auth.Post("/login", authController.Login)
	auth.Post("/logout", authController.Logout) // <-- Ahora usa authController.Logout
}
