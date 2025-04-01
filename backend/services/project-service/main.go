package main

import (
	"log"

	"project-service/controllers"
	"project-service/repository"
	"project-service/routes"

	"project-service/config"
	"project-service/db"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.LoadConfig()
	database := db.ConnectDB(cfg.DatabaseURL)
	defer database.Close()

	// Crear repositorios y controladores
	projectRepo := &repository.ProjectRepository{DB: database}
	projectController := &controllers.ProjectController{Repo: projectRepo}

	app := fiber.New()

	// Configurar rutas
	routes.SetupRoutes(app, projectController)

	log.Fatal(app.Listen(":" + cfg.Port))
}
