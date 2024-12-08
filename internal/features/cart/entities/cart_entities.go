package cart_entities

import "github.com/google/uuid"

type (
	Cart struct {
		ID     uuid.UUID  `db:"id" json:"id" binding:"required,uuid"`
		UserID uuid.UUID  `db:"user_id" json:"-" binding:"required,uuid"`
		Items  []CartItem `db:"-" json:"items" binding:"required"`
		Total  float64    `db:"-" json:"total" binding:"required"`
	}

	CartItem struct {
		ID           uuid.UUID `db:"id" json:"id" binding:"required,uuid"`
		ProductID    uuid.UUID `db:"product_id" json:"product_id" binding:"required,uuid"`
		ProductName  string    `db:"product_name" json:"product_name" binding:"required"`
		ProductPrice float64   `db:"product_price" json:"product_price" binding:"required"`
		Quantity     int       `db:"quantity" json:"quantity" binding:"required"`
		TotalPrice   float64   `db:"total_price" json:"total_price" binding:"required"`
	}
)
