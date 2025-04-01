package main

import (
	"auth-service/config"
	"auth-service/controllers"
	"auth-service/db"
	"auth-service/repository"
	"auth-service/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {

	cfg := config.LoadConfig()
	database := db.ConnectDB(cfg.DatabaseURL)
	defer database.Close()

	

	// Crear repositorios y controladores
	authRepo := &repository.AuthRepository{DB: database}
	authController := &controllers.AuthController{Repo: authRepo}

	// Agregar controlador de usuario y rutas protegidas
	userController := controllers.NewUserController()

	// Configurar Fiber
	app := fiber.New()
	routes.SetupAuthRoutes(app, authController)

	
	routes.SetupUserRoutes(app, userController)


	log.Fatal(app.Listen(":" + cfg.Port))

}
