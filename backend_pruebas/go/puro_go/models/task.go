package models

import "time"

type Task struct {
	ID            int        `json:"id"`
	ParentID      *int       `json:"parent_id,omitempty"` // Puede ser nulo
	Progress      float64    `json:"progress"`
	EstimatedTime int        `json:"estimated_time"` // en segundos
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
