package cart_interface

import (
	"context"
	cart_entities "final/internal/features/cart/entities"
)

type (
	CartRepo interface {
		GetCart(ctx context.Context, userID string) (*cart_entities.Cart, error)
		AddToCart(ctx context.Context, userID, productID string, quantity int) error
		RemoveFromCart(ctx context.Context, userID, productID string) error
		UpdateCart(ctx context.Context, userID, productID string, quantity int) error
		DeleteCart(ctx context.Context, userID string) error
	}

	CartUsecase interface {
		GetCart(ctx context.Context, userID string) (*cart_entities.Cart, error)
		AddToCart(ctx context.Context, userID, productID string, quantity int) error
		RemoveFromCart(ctx context.Context, userID, productID string) error
		UpdateCart(ctx context.Context, userID, productID string, quantity int) error
		DeleteCart(ctx context.Context, userID string) error
	}

	CartRepoForOrderUsecase interface {
		GetCart(ctx context.Context, userID string) (*cart_entities.Cart, error)
		DeleteCart(ctx context.Context, userID string) error
	}
)
