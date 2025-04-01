package repository

import (
	"context"
	"fmt"
	"project-service/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repositorio para manejar proyectos en la base de datos
type ProjectRepository struct {
	DB *pgxpool.Pool
}

// Crear nuevo proyecto
// func (r *ProjectRepository) CreateProject(project models.Project) error {
// 	_, err := r.DB.Exec(context.Background(),
// 		"INSERT INTO Project (projectId, title, description, creatorUserId, isPrivate) VALUES ($1, $2, $3, $4, $5)",
// 		project.ID, project.Title, project.Description, project.CreatorID, project.IsPrivate,
// 	)
// 	return err
// }

// func (r *ProjectRepository) CreateProject(project *models.Project) error {
// 	err := r.DB.QueryRow(context.Background(),
// 		"INSERT INTO Project (title, description, creatorUserId, isPrivate) VALUES ($1, $2, $3, $4) RETURNING projectId",
// 		project.Title, project.Description, project.CreatorID, project.IsPrivate,
// 	).Scan(&project.ID) // Recuperar el ID generado por la BD
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (r *ProjectRepository) CreateProject(project *models.Project) error {
	var id uuid.UUID // Variable temporal para recibir el UUID generado por PostgreSQL
	
	err := r.DB.QueryRow(context.Background(),
		"INSERT INTO Project (title, description, creatorUserId, isPrivate) VALUES ($1, $2, $3, $4) RETURNING projectId",
		project.Title, project.Description, project.CreatorID, project.IsPrivate,
	).Scan(&id)

	if err != nil {
		return err
	}

	project.ID = &id // ✅ Asignar la dirección de memoria del UUID generado
	return nil
}

// func (r *ProjectRepository) CreateProject(project *models.Project) error {
// 	err := r.DB.QueryRow(context.Background(),
// 		`INSERT INTO Project (title, description, creatorUserId, isPrivate) 
// 		 VALUES ($1, $2, $3, COALESCE($4, DEFAULT)) 
// 		 RETURNING projectId, isPrivate`,
// 		project.Title, project.Description, project.CreatorID, project.IsPrivate,
// 	).Scan(&project.ID, project.IsPrivate) // Recuperar el ID y isPrivate generado por la BD

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }




// Obtener todos los proyectos
func (r *ProjectRepository) GetProjects() ([]models.Project, error) {
	rows, err := r.DB.Query(context.Background(),
		"SELECT projectId, title, description, creatorUserId, isPrivate FROM Project")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		if err := rows.Scan(&project.ID, &project.Title, &project.Description, &project.CreatorID, &project.IsPrivate); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	// Retornar un array vacío en lugar de nil
	if len(projects) == 0 {
		return []models.Project{}, nil
	}

	return projects, nil
}

// Actualizar un proyecto
func (r *ProjectRepository) UpdateProject(project *models.Project) error {
	_, err := r.DB.Exec(context.Background(),
		`UPDATE Project 
		 SET title = $1, description = $2, isPrivate = $3, updatedAt = CURRENT_TIMESTAMP
		 WHERE projectId = $4`,
		project.Title, project.Description, project.IsPrivate, project.ID,
	)
	return err
}

// Eliminar un proyecto por ID
func (r *ProjectRepository) DeleteProject(projectID uuid.UUID) error {
	_, err := r.DB.Exec(context.Background(),
		"DELETE FROM Project WHERE projectId = $1", projectID)
	return err
}

// Actualizar parcialmente un proyecto
func (r *ProjectRepository) UpdateProjectPartial(projectID uuid.UUID, updates map[string]interface{}) error {
	// Si no hay campos para actualizar, retornar sin hacer nada
	if len(updates) == 0 {
		return nil
	}

	// Construir la consulta dinámica
	query := "UPDATE Project SET "
	values := []interface{}{}
	i := 1

	for key, value := range updates {
		query += key + " = $" + fmt.Sprintf("%d", i) + ", "
		values = append(values, value)
		i++
	}

	// Eliminar la última coma y agregar la condición WHERE
	query = query[:len(query)-2] + " WHERE projectId = $" + fmt.Sprintf("%d", i)
	values = append(values, projectID)

	// Ejecutar la consulta SQL
	_, err := r.DB.Exec(context.Background(), query, values...)
	return err
}
