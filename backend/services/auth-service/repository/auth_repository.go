package repository

import (
	"auth-service/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	DB *pgxpool.Pool
}

// Obtener usuario por email
func (r *AuthRepository) GetUserByEmail(email string) (*models.UserAccount, error) {
	var user models.UserAccount
	err := r.DB.QueryRow(context.Background(),
		`SELECT userId, personId, email, passwordHash, nickname, userType, timezone, isActive, createdAt, updatedAt 
		 FROM UserAccount WHERE email=$1`, email).
		Scan(&user.UserID, &user.PersonID, &user.Email, &user.PasswordHash,
			&user.Nickname, &user.UserType, &user.Timezone, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
