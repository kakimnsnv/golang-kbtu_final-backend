package auth_entities

import (
	"final/pkg/auth"
	"time"

	"github.com/google/uuid"
)

type (
	LoginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	RegisterRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	Token string

	User struct {
		ID           uuid.UUID `db:"id" json:"id" binding:"required,uuid"`
		Email        string    `db:"email" json:"email" binding:"required,email"`
		PasswordHash string    `db:"password_hash" json:"-" binding:"required"`
		Role         auth.Role `db:"role" json:"role" binding:"required"`
		CreatedAt    time.Time `db:"created_at" json:"-" binding:"required"`
		UpdatedAt    time.Time `db:"updated_at" json:"-" binding:"required"`
	}
)
