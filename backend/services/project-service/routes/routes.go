package routes

import (
	"project-service/controllers"

	"github.com/gofiber/fiber/v2"
)

// Definir rutas
func SetupRoutes(app *fiber.App, projectController *controllers.ProjectController) {
	app.Get("/projects", projectController.GetProjects) // Nueva ruta para obtener proyectos
	app.Post("/projects", projectController.CreateProject)
	app.Put("/projects/:id", projectController.UpdateProject)  // Ruta para actualizar
	app.Patch("/projects/:id", projectController.UpdateProjectPartial) 
	app.Delete("/projects/:id", projectController.DeleteProject) // Ruta para eliminar

}
