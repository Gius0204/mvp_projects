package routes

import (
	"task-service/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, taskController *controllers.TaskController) {
	app.Post("/tasks", taskController.CreateTask)
	app.Get("/tasks", taskController.GetAllTasks)
	app.Get("/tasks/:id", taskController.GetTaskByID)
	app.Get("/tasks/:parentID/subtasks", taskController.GetSubtasks)
	app.Put("/tasks/:id", taskController.UpdateTask)
	app.Get("/tasks/:id/calculate", taskController.CalculateTaskHandler)

}
