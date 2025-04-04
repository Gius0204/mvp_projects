package main

import (
	"go_fiber/config" // Importando configuración
	"go_fiber/controllers"
	"go_fiber/db" // Importando conexión a BD
	"go_fiber/routes"
	"log"

	"go_fiber/ent"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Cargar configuración
	cfg := config.LoadConfig()

	driver := db.GetEntDriver(cfg.DatabaseURL)
	client := ent.NewClient(ent.Driver(driver))
	
	defer client.Close() // Cerrar cliente de ent al salir

	// Crear servidor Fiber
	app := fiber.New()

	// Inicializar controlador y pasar client de ent
	taskController := controllers.NewTaskController(client)

	// Configurar rutas
	routes.SetupRoutes(app, taskController)

	log.Fatal(app.Listen(":" + cfg.Port))
}
