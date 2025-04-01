package models

import (
	"time"

	"github.com/google/uuid"
)

type UserAccount struct {
	UserID       uuid.UUID `json:"userId"`
	PersonID     uuid.UUID `json:"personId"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Nickname     *string   `json:"nickname,omitempty"`
	UserType     string    `json:"userType"`
	Timezone     string    `json:"timezone"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
