package routes

import (
	"auth-service/controllers"
	"auth-service/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, userController *controllers.UserController) {
	user := app.Group("/user")

	// Proteger esta ruta con el middleware de autenticación
	user.Get("/profile", middlewares.AuthMiddleware, userController.GetProfile)
}
