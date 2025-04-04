package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"go_fiber/ent" // Asegúrate de importar el paquete generado de ent
	// Importa el paquete de task de ent
	"context"
	"log"
	"time"

	"go_fiber/ent/task"
)

type TaskController struct {
	Client *ent.Client
}

type Benchmark struct {
	ProgressCalcMS      float64 `json:"progress_calc_ms"`
	EstimatedTimeMS     float64 `json:"estimated_time_ms"`
	DBUpdateMS          float64 `json:"db_update_ms"`
	TotalHandlerTimeMS  float64 `json:"total_handler_time_ms"`
}

type TaskCalculateResponse struct {
	ID            int        `json:"id"`
	ParentID      *int       `json:"parent_id"`
	Progress      float64    `json:"progress"`
	EstimatedTime int        `json:"estimated_time"`
	CreatedAt     string     `json:"created_at"`
	UpdatedAt     string     `json:"updated_at"`
	Benchmark     Benchmark  `json:"benchmark"`
}


func NewTaskController(client *ent.Client) *TaskController {
	return &TaskController{Client: client}
}

// GET /tasks
func (tc *TaskController) GetTasks(c *fiber.Ctx) error {
	tasks, err := tc.Client.Task.
		Query().
		WithParent().
		WithSubtasks().
		All(context.Background())

	if err != nil {
		log.Printf("Error al obtener tareas: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al obtener tareas"})
	}

	return c.JSON(tasks)
}

// POST /tasks
func (tc *TaskController) CreateTask(c *fiber.Ctx) error {
	// Estructura para recibir el JSON del body
	type TaskInput struct {
		ParentID      *int    `json:"parent_id"`       // opcional
		Progress      float64 `json:"progress"`
		EstimatedTime int     `json:"estimated_time"`
	}

	var input TaskInput

	// Parsear JSON
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "JSON inválido",
		})
	}

	// Crear la tarea en la base de datos
	taskCreate := tc.Client.Task.
		Create().
		SetProgress(input.Progress).
		SetEstimatedTime(input.EstimatedTime)

	if input.ParentID != nil {
		taskCreate.SetParentID(*input.ParentID)
	}

	task, err := taskCreate.Save(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudo crear la tarea",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

// GET /tasks/:id/calculate
func (tc *TaskController) CalculateTaskHandler(c *fiber.Ctx) error {
	startTotal := time.Now()

	ctx := context.Background()
	idParam := c.Params("id")
	taskID, err := strconv.Atoi(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "ID inválido")
	}

	// Obtener la tarea padre
	parent, err := tc.Client.Task.Get(ctx, taskID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Tarea no encontrada")
	}

	// ───────────── Calcular Progreso ─────────────
	startProgress := time.Now()

	// Obtener subtareas directas
	subtasks, err := tc.Client.Task.
		Query().
		Where(task.ParentID(taskID)).
		All(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error consultando subtareas")
	}

	// Calcular progreso solo si tiene subtareas directas
	progress := 0.0
	if len(subtasks) > 0 {
		var weightedSum float64
		var totalWeight float64

		for _, st := range subtasks {
			weightedSum += st.Progress * float64(st.EstimatedTime)
			totalWeight += float64(st.EstimatedTime)
		}

		if totalWeight > 0 {
			progress = weightedSum / totalWeight
		}
	}
	elapsedProgress := time.Since(startProgress)

	// ───────────── Calcular Estimated Time ─────────────
	startTimeCalc := time.Now()
	// Calcular estimated_time anidado hasta 10 niveles
	totalEstimatedTime := tc.sumEstimatedTimeRecursive(ctx, taskID, 10)
	elapsedTimeCalc := time.Since(startTimeCalc)

	// ───────────── Actualizar BD ─────────────
	var updated *ent.Task
	startDB := time.Now()

	if len(subtasks) > 0 || totalEstimatedTime > 0 {
		updated, err = parent.Update().
			SetProgress(progress).
			SetEstimatedTime(totalEstimatedTime).
			Save(ctx)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error actualizando la tarea")
		}
	} else {
		updated = parent
	}
	elapsedDB := time.Since(startDB)

	// ───────────── Tiempo total ─────────────
	elapsedTotal := time.Since(startTotal)

	// ───────────── Respuesta ─────────────
	// return c.JSON(fiber.Map{
	// 	"id":                  updated.ID,
	// 	"parent_id":           updated.ParentID,
	// 	"progress":            updated.Progress,
	// 	"estimated_time":      updated.EstimatedTime,
	// 	"created_at":          updated.CreatedAt.Format(time.RFC3339),
	// 	"updated_at":          updated.UpdatedAt.Format(time.RFC3339),
	// 	"benchmark": fiber.Map{
	// 		"progress_calc_ms":       elapsedProgress.Milliseconds(),
	// 		"estimated_time_ms":      elapsedTimeCalc.Milliseconds(),
	// 		"db_update_ms":           elapsedDB.Milliseconds(),
	// 		"total_handler_time_ms":  elapsedTotal.Milliseconds(),
	// 	},
	// })

	response := TaskCalculateResponse{
		ID:            updated.ID,
		ParentID:      updated.ParentID,
		Progress:      updated.Progress,
		EstimatedTime: updated.EstimatedTime,
		CreatedAt:     updated.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     updated.UpdatedAt.Format(time.RFC3339),
		Benchmark: Benchmark{
			ProgressCalcMS:     elapsedProgress.Seconds() * 1000, // elapsedProgress.Milliseconds(),
			EstimatedTimeMS:    elapsedTimeCalc.Seconds() * 1000,
			DBUpdateMS:         elapsedDB.Seconds() * 1000,
			TotalHandlerTimeMS: elapsedTotal.Seconds() * 1000,
		},
	}
	
	return c.JSON(response)

	// return c.JSON(fiber.Map{
	// 	"id":             parent.ID,
	// 	"parent_id":      parent.ParentID,
	// 	"progress":       progress,
	// 	"estimated_time": totalEstimatedTime,
	// 	"created_at":     parent.CreatedAt.Format(time.RFC3339),
	// 	"updated_at":     time.Now().Format(time.RFC3339), // si quieres, puedes usar parent.UpdatedAt si lo tienes
	// })
}
func (tc *TaskController) sumEstimatedTimeRecursive(ctx context.Context, parentID int, depth int) int {
	if depth == 0 {
		return 0
	}

	subtasks, err := tc.Client.Task.
		Query().
		Where(task.ParentID(parentID)).
		All(ctx)
	if err != nil {
		return 0
	}

	sum := 0
	for _, st := range subtasks {
		sum += st.EstimatedTime
		sum += tc.sumEstimatedTimeRecursive(ctx, st.ID, depth-1)
	}
	return sum
}



// // PUT /tasks/:id
// func (tc *TaskController) UpdateTask(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	var input struct {
// 			Progress      float64 `json:"progress"`
// 			EstimatedTime int     `json:"estimated_time"`
// 	}
// 	if err := c.BodyParser(&input); err != nil {
// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
// 	}

// 	task, err := tc.Client.Task.UpdateOneID(id).
// 			SetProgress(input.Progress).
// 			SetEstimatedTime(input.EstimatedTime).
// 			Save(context.Background())
// 	if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al actualizar tarea"})
// 	}
// 	//return c.JSON(task)
// 	return c.JSON(fiber.Map{"message": "Tarea actualizada con éxito"})
// }

// // DELETE /tasks/:id
// func (tc *TaskController) DeleteTask(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	err := tc.Client.Task.DeleteOneID(id).Exec(context.Background())
// 	if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al eliminar tarea"})
// 	}
// 	return c.SendStatus(fiber.StatusNoContent) //.JSON(fiber.Map{"message": "Tarea eliminada con éxito"})
// }