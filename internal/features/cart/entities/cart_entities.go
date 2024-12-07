package cart_entities

import "github.com/google/uuid"

type (
	Cart struct {
		ID     uuid.UUID  `json:"id"`
		UserID uuid.UUID  `json:"user_id"`
		Items  []CartItem `json:"items"`
		Total  float64    `json:"total"`
	}

	CartItem struct {
		ID           uuid.UUID `json:"id"`
		ProductID    uuid.UUID `json:"product_id"`
		ProductName  string    `json:"product_name"`
		ProductPrice float64   `json:"product_price"`
		Quantity     int       `json:"quantity"`
		TotalPrice   float64   `json:"total_price"`
	}
)
