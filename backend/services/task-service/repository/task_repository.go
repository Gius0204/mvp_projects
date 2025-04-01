package repository

import (
	"context"
	"fmt"
	"log"
	"task-service/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	DB *pgxpool.Pool
}

// Crear una nueva tarea
func (r *TaskRepository) CreateTask(task *models.Task) error {
	var id uuid.UUID

	err := r.DB.QueryRow(context.Background(),
		"INSERT INTO Task (projectId, title, description, creatorUserId, statusId, priority, startDate, dueDate, estimatedHours, isMilestone, parentTaskId) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING taskId",
		task.ProjectID, task.Title, task.Description, task.CreatorID, task.StatusID, task.Priority,
		task.StartDate, task.DueDate, task.EstimatedHours, task.IsMilestone, task.ParentTaskID,
	).Scan(&id)

	if err != nil {
		log.Println("Error al crear tarea:", err)
		return err
	}

	task.ID = &id
	return nil
}

// Obtener todas las tareas
func (r *TaskRepository) GetTasks() ([]models.Task, error) {
	rows, err := r.DB.Query(context.Background(),
		"SELECT taskId, projectId, title, description, creatorUserId, statusId, priority, startDate, dueDate, estimatedHours, isMilestone, parentTaskId, createdDate, lastModifiedDate, completedDate FROM Task")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.ProjectID, &task.Title, &task.Description, &task.CreatorID,
			&task.StatusID, &task.Priority, &task.StartDate, &task.DueDate, &task.EstimatedHours, &task.IsMilestone,
			&task.ParentTaskID, &task.CreatedDate, &task.LastModifiedDate, &task.CompletedDate); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if len(tasks) == 0 {
		return []models.Task{}, nil
	}

	return tasks, nil
}

// Actualizar una tarea completa (PUT)
func (r *TaskRepository) UpdateTask(task *models.Task) error {
	_, err := r.DB.Exec(context.Background(),
		`UPDATE Task SET projectId=$1, title=$2, description=$3, creatorUserId=$4, 
		 statusId=$5, priority=$6, startDate=$7, dueDate=$8, estimatedHours=$9, 
		 isMilestone=$10, parentTaskId=$11, lastModifiedDate=NOW()
		 WHERE taskId=$12`,
		task.ProjectID, task.Title, task.Description, task.CreatorID, task.StatusID,
		task.Priority, task.StartDate, task.DueDate, task.EstimatedHours,
		task.IsMilestone, task.ParentTaskID, task.ID,
	)
	return err
}

// Actualizar parcialmente una tarea (PATCH)
func (r *TaskRepository) PatchTask(taskID uuid.UUID, updates map[string]interface{}) error {
	// Si no hay campos para actualizar, retornar sin hacer nada
	if len(updates) == 0 {
		return nil
	}

	// Construir la consulta dinámica
	query := "UPDATE Task SET "
	values := []interface{}{}
	i := 1

	for key, value := range updates {
		query += key + " = $" + fmt.Sprintf("%d", i) + ", "
		values = append(values, value)
		i++
	}

	// Agregar la fecha de última modificación
	query += "lastModifiedDate = NOW(), "

	// Eliminar la última coma y agregar la condición WHERE
	query = query[:len(query)-2] + " WHERE taskId = $" + fmt.Sprintf("%d", i)
	values = append(values, taskID)

	// Ejecutar la consulta SQL
	_, err := r.DB.Exec(context.Background(), query, values...)
	return err
}


// Eliminar una tarea (DELETE)
func (r *TaskRepository) DeleteTask(id uuid.UUID) error {
	_, err := r.DB.Exec(context.Background(), "DELETE FROM Task WHERE taskId=$1", id)
	return err
}
