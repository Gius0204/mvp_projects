package controllers

import (
	"context"
	"fmt"
	"strconv"
	"task-service/models"
	"task-service/repository"

	"github.com/gofiber/fiber/v2"
)

type TaskController struct {
	Repo *repository.TaskRepository
}

func NewTaskController(repo *repository.TaskRepository) *TaskController {
	return &TaskController{Repo: repo}
}

func (c *TaskController) CreateTask(ctx *fiber.Ctx) error {
	var task models.Task
	if err := ctx.BodyParser(&task); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	id, err := c.Repo.CreateTask(ctx.Context(), &task)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	task.ID = id
	return ctx.JSON(task)
}

func (c *TaskController) GetAllTasks(ctx *fiber.Ctx) error {
	tasks, err := c.Repo.GetAllTasks(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(tasks)
}

func (c *TaskController) GetTaskByID(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	task, err := c.Repo.GetTaskByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tarea no encontrada"})
	}

	return ctx.JSON(task)
}

func (c *TaskController) GetSubtasks(ctx *fiber.Ctx) error {
	parentID, err := strconv.Atoi(ctx.Params("parentID"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	subtasks, err := c.Repo.GetSubtasks(ctx.Context(), parentID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(subtasks)
}

func (c *TaskController) UpdateTask(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var update struct {
		Progress      float64 `json:"progress"`
		EstimatedTime int     `json:"estimated_time"`
	}
	if err := ctx.BodyParser(&update); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err = c.Repo.UpdateTask(ctx.Context(), id, update.Progress, update.EstimatedTime)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (c *TaskController) CalculateTaskHandler(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	taskID, err := strconv.Atoi(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid task ID")
	}

	progress, totalTime, err := c.Repo.CalculateTaskProgressAndTime(context.Background(), taskID)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"progress":       fmt.Sprintf("%.2f", progress),
		"estimated_time": totalTime,
	})
}
