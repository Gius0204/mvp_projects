package routes

import (
	"go_fiber/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, taskController *controllers.TaskController) {
	app.Get("/tasks", taskController.GetTasks)
	app.Post("/tasks", taskController.CreateTask)
	// app.Put("/tasks/:id", taskController.UpdateTask)
	// app.Patch("/tasks/:id", taskController.PatchTask)
	// app.Delete("/tasks/:id", taskController.DeleteTask)
	app.Get("/tasks/:id/calculate", taskController.CalculateTaskHandler)
}
