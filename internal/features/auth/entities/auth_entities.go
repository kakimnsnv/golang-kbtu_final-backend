package auth_entities

import (
	"final/internal/infrastructure/db_gen"
	"final/pkg/auth"
	"time"

	"github.com/google/uuid"
)

type (
	LoginRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	RegisterRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	Token string

	User struct {
		ID           uuid.UUID `json:"id"`
		Email        string    `json:"email"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		PasswordHash string    `json:"-"`
		Role         auth.Role `json:"role"`
	}
)

// MARK: Mappers
func FromDBUser(dbUser db_gen.User) User {
	return User{
		ID:           dbUser.ID,
		Email:        dbUser.Email,
		CreatedAt:    dbUser.CreatedAt.Time,
		UpdatedAt:    dbUser.UpdatedAt.Time,
		PasswordHash: dbUser.PasswordHash,
		Role:         auth.Role(dbUser.Role),
	}
}
