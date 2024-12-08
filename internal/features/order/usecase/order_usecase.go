package order_usecase

import (
	"context"
	cart_interface "final/internal/features/cart/interface"
	order_interface "final/internal/features/order/interface"

	"go.uber.org/zap"
)

type OrderUsecaseImpl struct {
	logger    *zap.Logger
	orderRepo order_interface.OrderRepo
	cartRepo  cart_interface.CartRepo
}

var _ order_interface.OrderUsecase = (*OrderUsecaseImpl)(nil)

func New(logger *zap.Logger, orderRepo order_interface.OrderRepo, cartRepo cart_interface.CartRepo) order_interface.OrderUsecase {
	return &OrderUsecaseImpl{
		logger:    logger,
		orderRepo: orderRepo,
		cartRepo:  cartRepo,
	}
}

func (u *OrderUsecaseImpl) PlaceOrder(ctx context.Context, userID string, cartID string) (string, error) {
	cart, err := u.cartRepo.GetCart(ctx, userID)
	if err != nil {
		u.logger.Error("failed to get cart", zap.Error(err))
		return "", err
	}
	u.orderRepo.CreateOrder(ctx, userID, cart.Total)

}

func (u *OrderUsecaseImpl) GetOrder(ctx context.Context, orderID string) (order_entities.Order, error) {

}

func (u *OrderUsecaseImpl) GetOrders(ctx context.Context, userID string) ([]order_entities.Order, error) {

}

func (u *OrderUsecaseImpl) CancelOrder(ctx context.Context, orderID string) error {

}

func (u *OrderUsecaseImpl) ChangeOrderStatus(ctx context.Context, orderID string, status int) error {

}
