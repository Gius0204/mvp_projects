// package main

// import (
// 	"api-gateway/routes"

// 	"github.com/gofiber/fiber/v2"
// )

// func main() {
// 	app := fiber.New()

// 	// Configurar rutas
// 	routes.SetupRoutes(app)

// 	// Iniciar servidor en el puerto 3000
// 	app.Listen(":3000")
// }

package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func main() {
	app := fiber.New()

	// Middleware de autenticación (Valida JWT antes de redirigir)
	app.Use(func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if !strings.HasPrefix(c.Path(), "/auth") { // No validar auth-service
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "No autorizado, se requiere token",
				})
			}
		}

		return c.Next()
	})

	// **Redirigir a project-service**
	app.All("/projects*", func(c *fiber.Ctx) error {
		targetURL := "http://localhost:3001" + c.OriginalURL()
		return proxy.Do(c, targetURL)
	})

	// **Redirigir a task-service**
	app.All("/tasks*", func(c *fiber.Ctx) error {
		targetURL := "http://localhost:3002" + c.OriginalURL()
		return proxy.Do(c, targetURL)
	})

	// **Redirigir a auth-service (login, registro)**
	app.All("/auth*", func(c *fiber.Ctx) error {
		targetURL := "http://localhost:3003" + c.OriginalURL()
		return proxy.Do(c, targetURL)
	})

	log.Println("🚀 API Gateway corriendo en http://localhost:8080")
	log.Fatal(app.Listen(":3000"))
}
