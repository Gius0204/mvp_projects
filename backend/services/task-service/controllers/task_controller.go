package controllers

import (
	"encoding/json"
	"fmt"
	"task-service/models"
	"task-service/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type TaskController struct {
	Repo *repository.TaskRepository
}

// Crear tarea
// func (t *TaskController) CreateTask(c *fiber.Ctx) error {
// 	var task models.Task
// 	if err := c.BodyParser(&task); err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
// 	}

// 	if err := t.Repo.CreateTask(&task); err != nil {
// 		return c.Status(500).JSON(fiber.Map{"error": "No se pudo crear la tarea"})
// 	}

// 	return c.Status(201).JSON(task)
// }

func (t *TaskController) CreateTask(c *fiber.Ctx) error {
	var input struct {
		ProjectID      uuid.UUID `json:"projectId"`
		Title          string    `json:"title"`
		Description    *string   `json:"description"`
		CreatorID      uuid.UUID `json:"creatorUserId"`
		StatusID       uuid.UUID `json:"statusId"`
		Priority       *int      `json:"priority"`
		StartDate      string    `json:"startDate"`
		DueDate        string    `json:"dueDate"`
		EstimatedHours *float64  `json:"estimatedHours"`
		IsMilestone    *bool     `json:"isMilestone"`
		ParentTaskID   *uuid.UUID `json:"parentTaskId"`
	}

	// Parsear el JSON del body
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	// Convertir startDate y dueDate de string a time.Time
	var startDate, dueDate *time.Time

	if input.StartDate != "" {
		parsedStartDate, err := time.Parse("2006-01-02", input.StartDate) //2006-01-02 es la plantilla de go para parsear xD
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Formato de fecha inválido en startDate, usa YYYY-MM-DD"})
		}
		startDate = &parsedStartDate
	}

	if input.DueDate != "" {
		parsedDueDate, err := time.Parse("2006-01-02", input.DueDate)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Formato de fecha inválido en dueDate, usa YYYY-MM-DD"})
		}
		dueDate = &parsedDueDate
	}

	// Crear la tarea con los datos convertidos
	task := models.Task{
		ProjectID:      input.ProjectID,
		Title:          input.Title,
		Description:    input.Description,
		CreatorID:      input.CreatorID,
		StatusID:       input.StatusID,
		Priority:       input.Priority,
		StartDate:      startDate,
		DueDate:        dueDate,
		EstimatedHours: input.EstimatedHours,
		IsMilestone:    input.IsMilestone,
		ParentTaskID:   input.ParentTaskID,
	}

	// Guardar en la BD
	if err := t.Repo.CreateTask(&task); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "No se pudo crear la tarea"})
	}

	// Enviar mensaje a Kafka
	sendKafkaMessage("notifications", fiber.Map{
		"type":    "new_task",
		"message": "Se ha creado una nueva tarea",
		"data":    task,
	})

	return c.Status(201).JSON(task)
}


// Obtener todas las tareas
func (t *TaskController) GetTasks(c *fiber.Ctx) error {
	tasks, err := t.Repo.GetTasks()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "No se pudieron obtener las tareas"})
	}

	return c.JSON(tasks)
}

// Actualizar una tarea completa (PUT)

// Actualizar una tarea completa (PUT)
func (t *TaskController) UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	taskID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	var input struct {
		ProjectID      uuid.UUID  `json:"projectId"`
		Title          string     `json:"title"`
		Description    *string    `json:"description"`
		CreatorID      uuid.UUID  `json:"creatorUserId"`
		StatusID       uuid.UUID  `json:"statusId"`
		Priority       *int       `json:"priority"`
		StartDate      *string    `json:"startDate"`
		DueDate        *string    `json:"dueDate"`
		EstimatedHours *float64   `json:"estimatedHours"`
		IsMilestone    *bool      `json:"isMilestone"`
		ParentTaskID   *uuid.UUID `json:"parentTaskId"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	// Convertir fechas de string a time.Time
	var startDate, dueDate *time.Time
	if input.StartDate != nil {
		parsedStartDate, err := time.Parse("2006-01-02", *input.StartDate)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Formato de fecha inválido en startDate, usa YYYY-MM-DD"})
		}
		startDate = &parsedStartDate
	}
	if input.DueDate != nil {
		parsedDueDate, err := time.Parse("2006-01-02", *input.DueDate)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Formato de fecha inválido en dueDate, usa YYYY-MM-DD"})
		}
		dueDate = &parsedDueDate
	}

	// Construir la tarea
	task := models.Task{
		ID:             &taskID,
		ProjectID:      input.ProjectID,
		Title:          input.Title,
		Description:    input.Description,
		CreatorID:      input.CreatorID,
		StatusID:       input.StatusID,
		Priority:       input.Priority,
		StartDate:      startDate,
		DueDate:        dueDate,
		EstimatedHours: input.EstimatedHours,
		IsMilestone:    input.IsMilestone,
		ParentTaskID:   input.ParentTaskID,
	}

	if err := t.Repo.UpdateTask(&task); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "No se pudo actualizar la tarea"})
	}

	return c.JSON(fiber.Map{"message": "Tarea actualizada con éxito"})
}


// func (t *TaskController) UpdateTask(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	taskID, err := uuid.Parse(id)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
// 	}

// 	var task models.Task
// 	if err := c.BodyParser(&task); err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
// 	}
// 	task.ID = &taskID

// 	if err := t.Repo.UpdateTask(&task); err != nil {
// 		return c.Status(500).JSON(fiber.Map{"error": "No se pudo actualizar la tarea"})
// 	}

// 	return c.JSON(fiber.Map{"message": "Tarea actualizada con éxito"})
// }

// Actualizar parcialmente una tarea (PATCH)
func (t *TaskController) PatchTask(c *fiber.Ctx) error {
	id := c.Params("id")
	taskID, err := uuid.Parse(id)
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

	// Validar campos permitidos
	validFields := map[string]bool{
		"projectId":      true,
		"title":          true,
		"description":    true,
		"creatorUserId":  true,
		"statusId":       true,
		"priority":       true,
		"startDate":      true,
		"dueDate":        true,
		"estimatedHours": true,
		"isMilestone":    true,
		"parentTaskId":   true,
	}

	for field := range updates {
		if !validFields[field] {
			return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("Campo '%s' no permitido", field)})
		}
	}

	// Convertir fechas si existen
	if startDateStr, ok := updates["startDate"].(string); ok {
		parsedStartDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Formato de fecha inválido en startDate, usa YYYY-MM-DD"})
		}
		updates["startDate"] = parsedStartDate
	}

	if dueDateStr, ok := updates["dueDate"].(string); ok {
		parsedDueDate, err := time.Parse("2006-01-02", dueDateStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Formato de fecha inválido en dueDate, usa YYYY-MM-DD"})
		}
		updates["dueDate"] = parsedDueDate
	}

	// Llamar al repositorio para actualizar los campos específicos
	if err := t.Repo.PatchTask(taskID, updates); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "No se pudo actualizar la tarea parcialmente"})
	}

	return c.JSON(fiber.Map{"message": "Tarea actualizada parcialmente con éxito"})
}



// func (t *TaskController) PatchTask(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	taskID, err := uuid.Parse(id)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
// 	}

// 	var task models.Task
// 	if err := c.BodyParser(&task); err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
// 	}
// 	task.ID = &taskID

// 	if err := t.Repo.PatchTask(&task); err != nil {
// 		return c.Status(500).JSON(fiber.Map{"error": "No se pudo actualizar la tarea parcialmente"})
// 	}

// 	return c.JSON(fiber.Map{"message": "Tarea actualizada parcialmente con éxito"})
// }

// Eliminar una tarea (DELETE)
func (t *TaskController) DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	taskID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	if err := t.Repo.DeleteTask(taskID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "No se pudo eliminar la tarea"})
	}

	return c.JSON(fiber.Map{"message": "Tarea eliminada con éxito"})
}

//tratara de separar en el futuro a KAFKA
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