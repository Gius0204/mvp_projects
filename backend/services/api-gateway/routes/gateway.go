package routes

import (
	middleware "api-gateway/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", middleware.AuthMiddleware) // Todas las rutas requieren autenticación

	// Redirigir a project-service
	api.All("/projects/*", func(c *fiber.Ctx) error {
		return proxy.Do(c, "http://localhost:3001"+c.OriginalURL())
	})

	// Redirigir a task-service
	api.All("/tasks/*", func(c *fiber.Ctx) error {
		return proxy.Do(c, "http://localhost:3002"+c.OriginalURL())
	})
}
