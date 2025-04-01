package controllers

import (
	"encoding/json"
	"fmt"
	"project-service/models"
	"project-service/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

// Controlador de proyectos
type ProjectController struct {
	Repo *repository.ProjectRepository
}

// Crear proyecto
func (p *ProjectController) CreateProject(c *fiber.Ctx) error {
	var project models.Project
	if err := c.BodyParser(&project); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	// Si is_private no se especifica, asignar el valor por defecto `false`
	if project.IsPrivate == nil {
		defaultPrivate := false
		project.IsPrivate = &defaultPrivate
	}

	// Generar UUID
	// project.ID = uuid.New()

	if err := p.Repo.CreateProject(&project); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "No se pudo crear el proyecto"})
	}


	// Enviar mensaje a Kafka
	sendKafkaMessage("notifications", fiber.Map{
		"type":    "new_project",
		"message": "Se ha creado un nuevo proyecto",
		"data":    project,
	})
	
	return c.Status(201).JSON(project)
}
// Obtener todos los proyectos
func (p *ProjectController) GetProjects(c *fiber.Ctx) error {
	projects, err := p.Repo.GetProjects()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "No se pudieron obtener los proyectos"})
	}

	return c.JSON(projects)
}

// Actualizar un proyecto
func (p *ProjectController) UpdateProject(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	var project models.Project
	if err := c.BodyParser(&project); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	project.ID = &projectID // Asignamos el ID desde la URL

	if err := p.Repo.UpdateProject(&project); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "No se pudo actualizar el proyecto"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Proyecto actualizado correctamente"})
}

// Eliminar un proyecto
func (p *ProjectController) DeleteProject(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	if err := p.Repo.DeleteProject(projectID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "No se pudo eliminar el proyecto"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Proyecto eliminado correctamente"})
}

// Actualizar parcialmente un proyecto
func (p *ProjectController) UpdateProjectPartial(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	// Leer el cuerpo de la solicitud en un mapa genérico
	var updates map[string]interface{}
	if err := json.Unmarshal(c.Body(), &updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	// Verificar que al menos un campo se quiere actualizar
	if len(updates) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "No hay campos para actualizar"})
	}

	// Validar que los campos sean válidos
	validFields := map[string]bool{
		"title":       true,
		"description": true,
		"isPrivate":   true,
	}

	for field := range updates {
		if !validFields[field] {
			return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("Campo '%s' no permitido", field)})
		}
	}

	// Llamar al repositorio para actualizar los campos específicos
	if err := p.Repo.UpdateProjectPartial(projectID, updates); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "No se pudo actualizar el proyecto"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Proyecto actualizado correctamente"})
}

//adelante tratare de separar el kafka
func sendKafkaMessage(topic string, message interface{}) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, 0)
	if err != nil {
		log.Println("Error al conectar con Kafka:", err)
		return
	}
	defer conn.Close()

	msg, _ := json.Marshal(message)
	_, err = conn.WriteMessages(kafka.Message{Value: msg})
	if err != nil {
		log.Println("Error al enviar mensaje a Kafka:", err)
	}
}