package order_entities

import (
	"time"

	"github.com/google/uuid"
)

type (
	Order struct {
		ID          uuid.UUID   `json:"id" db:"id" binding:"required"`
		UserID      uuid.UUID   `json:"user_id" db:"user_id" binding:"required"`
		TotalAmount float64     `json:"total_amount" db:"total_amount" binding:"required"`
		Status      OrderStatus `json:"status" db:"status" binding:"required"`
		CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	}
	OrderStatus int

	OrderItem struct {
		ID        uuid.UUID `json:"id" db:"id" binding:"required"`
		OrderID   uuid.UUID `json:"order_id" db:"order_id" binding:"required"`
		ProductID uuid.UUID `json:"product_id" db:"product_id" binding:"required"`
		Quantity  int       `json:"quantity" db:"quantity" binding:"required"`
		Price     float64   `json:"price" db:"price" binding:"required"`
		CreatedAt time.Time `json:"created_at" db:"created_at"`
	}
)

const (
	OrderStatusPending OrderStatus = iota
	OrderStatusProcessing
	OrderStatusShipping
	OrderStatusDelivered
	OrderStatusCancelled
)
