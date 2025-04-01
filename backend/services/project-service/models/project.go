package models

import "github.com/google/uuid"

type Project struct {
	ID          *uuid.UUID `json:"id_project"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	CreatorID   uuid.UUID `json:"creator_user_id" validate:"required"`
	IsPrivate   *bool      `json:"is_private,omitempty"` // Permite que sea opcional en JSON"
}
