package product_entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type (
	Product struct {
		ID          uuid.UUID    `db:"id" json:"id" binding:"required,uuid"`
		Name        string       `db:"name" json:"name" binding:"required"`
		Photo       string       `db:"photo" json:"photo" binding:"required"`
		Description string       `db:"description" json:"description" binding:"required"`
		Price       float64      `db:"price" json:"price" binding:"required"`
		IsLiked     bool         `db:"is_liked" json:"is_liked" binding:"required"`
		CreatedAt   time.Time    `db:"created_at" json:"-" binding:"required"`
		UpdatedAt   time.Time    `db:"updated_at" json:"-" binding:"required"`
		DeletedAt   sql.NullTime `db:"deleted_at" json:"-" binding:"required"`
	}

	ProductRequest struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description" binding:"required"`
		Price       float64 `json:"price" binding:"required"`
		Photo       string  `json:"photo" binding:"required"`
	}
)
