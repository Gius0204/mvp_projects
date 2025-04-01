package main

import (
	"log"
	"task-service/config" // Importando configuración
	"task-service/controllers"
	"task-service/db" // Importando conexión a BD
	"task-service/repository"
	"task-service/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Cargar configuración
	cfg := config.LoadConfig()

	// Conectar a la base de datos
	database := db.ConnectDB(cfg.DatabaseURL)
	defer database.Close() // Cerrar conexión al salir

	// Inicializar repositorios y controladores
	taskRepo := &repository.TaskRepository{DB: database}
	taskController := &controllers.TaskController{Repo: taskRepo}

	// Crear servidor Fiber
	app := fiber.New()

	// Configurar rutas
	routes.SetupRoutes(app, taskController)

	log.Fatal(app.Listen(":" + cfg.Port))
}
