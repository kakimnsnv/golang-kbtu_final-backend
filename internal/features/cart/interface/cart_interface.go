package cart_interface

import (
	"context"
	cart_entities "final/internal/features/cart/entities"

	"github.com/google/uuid"
)

type (
	CartRepo interface {
		CreateCart(ctx context.Context, userID uuid.UUID) (*cart_entities.Cart, error)
		GetCart(ctx context.Context, userID uuid.UUID) (*cart_entities.Cart, error)
		AddToCart(ctx context.Context, userID uuid.UUID, productID, quantity int) error
		RemoveFromCart(ctx context.Context, userID uuid.UUID, productID int) error
		UpdateCart(ctx context.Context, userID uuid.UUID, productID, quantity int) error
	}

	CartUsecase interface {
		GetCart(ctx context.Context, userID int) (*cart_entities.Cart, error)
		AddToCart(ctx context.Context, userID, productID, quantity int) error
		RemoveFromCart(ctx context.Context, userID, productID int) error
		UpdateCart(ctx context.Context, userID, productID, quantity int) error
	}
)
