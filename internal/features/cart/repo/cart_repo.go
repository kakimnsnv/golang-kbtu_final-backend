package cart_repo

import (
	"context"
	"encoding/json"
	cart_entities "final/internal/features/cart/entities"
	cart_interface "final/internal/features/cart/interface"
	"final/internal/infrastructure/db_gen"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CartRepoImpl struct {
	logger *zap.Logger
	db     *db_gen.Queries
}

var _ cart_interface.CartRepo = (*CartRepoImpl)(nil)

func New(logger *zap.Logger, db *db_gen.Queries) cart_interface.CartRepo {
	return &CartRepoImpl{
		logger: logger,
		db:     db,
	}
}

func (r *CartRepoImpl) CreateCart(ctx context.Context, userID uuid.UUID) (*cart_entities.Cart, error) {
	existingCart, err := r.db.GetCart(ctx, userID)
	if err == nil && existingCart.CartID != uuid.Nil {
		cart := &cart_entities.Cart{
			ID:     existingCart.CartID,
			UserID: existingCart.UserID,
			Items:  []cart_entities.CartItem{},
		}

	}
}

func (r *CartRepoImpl) GetCart(ctx context.Context, userID uuid.UUID) (*cart_entities.Cart, error) {
	result, err := r.db.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	cart := &cart_entities.Cart{
		ID:     result.CartID,
		UserID: result.UserID,
	}

	if err := json.Unmarshal(result.Items, &cart.Items); err != nil {
		return nil, err
	}

	cart.Total = calculateTotal(cart.Items)

	return cart, nil
}

func calculateTotal(items []cart_entities.CartItem) float64 {
	var total float64
	for _, item := range items {
		total += item.ProductPrice * float64(item.Quantity)
	}

	return total
}

func (r *CartRepoImpl) AddToCart(ctx context.Context, userID uuid.UUID, productID, quantity int) error {

}

func (r *CartRepoImpl) RemoveFromCart(ctx context.Context, userID uuid.UUID, productID int) error {

}

func (r *CartRepoImpl) UpdateCart(ctx context.Context, userID uuid.UUID, productID, quantity int) error {

}
