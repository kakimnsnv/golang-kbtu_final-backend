package cart_usecase

import (
	"context"
	cart_entities "final/internal/features/cart/entities"
	cart_interface "final/internal/features/cart/interface"

	"go.uber.org/zap"
)

type CartUsecaseImpl struct {
	logger *zap.Logger
	repo   cart_interface.CartRepo
}

func New(repo cart_interface.CartRepo, logger *zap.Logger) *CartUsecaseImpl {
	return &CartUsecaseImpl{
		logger: logger,
		repo:   repo,
	}
}

var _ cart_interface.CartUsecase = (*CartUsecaseImpl)(nil)

func (u *CartUsecaseImpl) GetCart(ctx context.Context, userID string) (*cart_entities.Cart, error) {
	return u.repo.GetCart(ctx, userID)
}

func (u *CartUsecaseImpl) AddToCart(ctx context.Context, userID, productID string, quantity int) error {
	return u.repo.AddToCart(ctx, userID, productID, quantity)
}

func (u *CartUsecaseImpl) RemoveFromCart(ctx context.Context, userID, productID string) error {
	return u.repo.RemoveFromCart(ctx, userID, productID)
}

func (u *CartUsecaseImpl) UpdateCart(ctx context.Context, userID, productID string, quantity int) error {
	return u.repo.UpdateCart(ctx, userID, productID, quantity)
}

func (u *CartUsecaseImpl) DeleteCart(ctx context.Context, userID string) error {
	return u.repo.DeleteCart(ctx, userID)
}
