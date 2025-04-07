package repository

import (
	"context"
	"task-service/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	DB *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{DB: db}
}

// Crear tarea
func (r *TaskRepository) CreateTask(ctx context.Context, task *models.Task) (int, error) {
	query := `
		INSERT INTO tasks (parent_id, progress, estimated_time)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var newID int
	err := r.DB.QueryRow(ctx, query, task.ParentID, task.Progress, task.EstimatedTime).Scan(&newID)
	return newID, err
}

// Obtener todas las tareas
func (r *TaskRepository) GetAllTasks(ctx context.Context) ([]*models.Task, error) {
	query := `SELECT id, parent_id, progress, estimated_time, created_at, updated_at FROM tasks`
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var t models.Task
		err := rows.Scan(&t.ID, &t.ParentID, &t.Progress, &t.EstimatedTime, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, nil
}

// Obtener subtareas por parent_id
func (r *TaskRepository) GetSubtasks(ctx context.Context, parentID int) ([]*models.Task, error) {
	query := `SELECT id, parent_id, progress, estimated_time, created_at, updated_at FROM tasks WHERE parent_id = $1`
	rows, err := r.DB.Query(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subtasks []*models.Task
	for rows.Next() {
		var t models.Task
		err := rows.Scan(&t.ID, &t.ParentID, &t.Progress, &t.EstimatedTime, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		subtasks = append(subtasks, &t)
	}
	return subtasks, nil
}

// Actualizar progreso y tiempo estimado
func (r *TaskRepository) UpdateTask(ctx context.Context, id int, progress float64, estimatedTime int) error {
	query := `
		UPDATE tasks
		SET progress = $1, estimated_time = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.DB.Exec(ctx, query, progress, estimatedTime, time.Now(), id)
	return err
}

// Obtener tarea por ID
func (r *TaskRepository) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
	query := `SELECT id, parent_id, progress, estimated_time, created_at, updated_at FROM tasks WHERE id = $1`
	var t models.Task
	err := r.DB.QueryRow(ctx, query).Scan(&t.ID, &t.ParentID, &t.Progress, &t.EstimatedTime, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepository) CalculateTaskProgressAndTime(ctx context.Context, taskID int) (float64, int, error) {
	// Verificar si la tarea existe
	var exists bool
	err := r.DB.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM tasks WHERE id = $1)`, taskID).Scan(&exists)
	if err != nil {
		return 0, 0, err
	}
	if !exists {
		return 0, 0, fiber.NewError(fiber.StatusNotFound, "Task not found")
	}

	// Progreso ponderado de subtareas directas
	var progress float64
	err = r.DB.QueryRow(ctx, `
		SELECT 
			CASE 
				WHEN SUM(estimated_time) > 0 THEN SUM(progress * estimated_time) / SUM(estimated_time)
				ELSE 0 
			END AS weighted_progress
		FROM tasks 
		WHERE parent_id = $1
	`, taskID).Scan(&progress)
	if err != nil {
		return 0, 0, err
	}

	// Tiempo estimado hasta 10 niveles de profundidad
	var totalEstimatedTime int
	err = r.DB.QueryRow(ctx, `
		WITH RECURSIVE subtasks AS (
			SELECT id, estimated_time, parent_id, 1 as level
			FROM tasks
			WHERE parent_id = $1
			UNION ALL
			SELECT t.id, t.estimated_time, t.parent_id, st.level + 1
			FROM tasks t
			INNER JOIN subtasks st ON t.parent_id = st.id
			WHERE st.level < 10
		)
		SELECT COALESCE(SUM(estimated_time), 0) FROM subtasks
	`, taskID).Scan(&totalEstimatedTime)
	if err != nil {
		return 0, 0, err
	}

	return progress, totalEstimatedTime, nil
}
