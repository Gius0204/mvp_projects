package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID               *uuid.UUID `json:"id_task"`
	ProjectID        uuid.UUID  `json:"projectId" validate:"required"`
	Title            string     `json:"title" validate:"required"`
	Description      *string    `json:"description,omitempty"`
	CreatorID        uuid.UUID  `json:"creatorUserId" validate:"required"`
	StatusID         uuid.UUID  `json:"statusId" validate:"required"`
	Priority         *int       `json:"priority,omitempty"`
	StartDate        *time.Time   `json:"startDate,omitempty"`
	DueDate          *time.Time   `json:"dueDate,omitempty"`
	EstimatedHours   *float64   `json:"estimatedHours,omitempty"`
	IsMilestone      *bool      `json:"isMilestone,omitempty"`
	ParentTaskID     *uuid.UUID `json:"parentTaskId,omitempty"`
	CreatedDate      *time.Time   `json:"created_date,omitempty"`
	LastModifiedDate *time.Time   `json:"last_modified_date,omitempty"`
	CompletedDate    *time.Time   `json:"completed_date,omitempty"`
}





//quizas luego cambiar startdate, duedate a *time.Time
//ojo los json q se colocan es como se muestran, se envian con POST,PUT,...
// y deben respetarse, xq sino no da xq mas q todo este bien
//ejm: si coloco en el json q voy a enviar project-id: "..." aunque este bien el id
// aqui ya indicado q la forma para ProjectID es json: "projectId" no "project-id"

//cuidado con colocar lo de la base de datos fechas aqui en *string, deben ser *time.Time

//este juego de cambiar string y time afecta a get y post xd, asi q mejor se crea una str