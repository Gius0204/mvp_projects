package main

import (
	"context"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool
var validate = validator.New()
var clients = make(map[*websocket.Conn]bool)

type Task struct {
	ID    int    `json:"id_task"`
	Title string `json:"title" validate:"required"`
	Type  string `json:"type_view" validate:"required,oneof=list kanban gantt"`
}

func main() {
	var err error
	db, err = pgxpool.New(context.Background(), "postgres://postgres:postgresql123@localhost:5432/mvp123")
	if err != nil {
		log.Fatal("Error al conectar con PostgreSQL:", err)
	}
	defer db.Close()

	app := fiber.New()

	//Middlewares
	app.Use(logger.New()) //para ver las peticiones por consola

	// Middleware CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",                          // Permite cualquier origen
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS", // Métodos permitidos
		AllowHeaders: "Content-Type, Authorization", // Headers permitidos
	}))


	// Ruta WebSocket para la comunicación en tiempo real
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		clients[c] = true
		defer func() {
			delete(clients, c)
			c.Close()
		}()

		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				break
			}
		}
	}))


	// Rutas CRUD
	app.Get("/tasks", getTasks)
	app.Post("/tasks", createTask)

	log.Fatal(app.Listen(":3000"))
}

func getTasks(c *fiber.Ctx) error {
	rows, err := db.Query(context.Background(), "SELECT id_task, title, type_view FROM task")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Type); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		tasks = append(tasks, task)
	}

	return c.JSON(tasks)
}

func createTask(c *fiber.Ctx) error {
	var task Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validate.Struct(task); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := db.Exec(context.Background(), "INSERT INTO task (title, type_view) VALUES ($1, $2)", task.Title, task.Type)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Notificar a todos los clientes WebSocket
	for client := range clients {
		client.WriteJSON(task)
	}

	return c.Status(201).JSON(task)
}
