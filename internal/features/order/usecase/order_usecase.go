package order_usecase

import (
	"context"
	cart_interface "final/internal/features/cart/interface"
	order_entities "final/internal/features/order/entities"
	order_interface "final/internal/features/order/interface"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type OrderUsecaseImpl struct {
	logger    *zap.Logger
	orderRepo order_interface.OrderRepo
	cartRepo  cart_interface.CartRepoForOrderUsecase
}

var _ order_interface.OrderUsecase = (*OrderUsecaseImpl)(nil)

func New(logger *zap.Logger, orderRepo order_interface.OrderRepo, cartRepo cart_interface.CartRepoForOrderUsecase) order_interface.OrderUsecase {
	return &OrderUsecaseImpl{
		logger:    logger,
		orderRepo: orderRepo,
		cartRepo:  cartRepo,
	}
}

func (u *OrderUsecaseImpl) PlaceOrder(ctx context.Context, userID string) (string, error) {
	cart, err := u.cartRepo.GetCart(ctx, userID)
	if err != nil {
		u.logger.Error("failed to get cart", zap.Error(err))
		return "", err
	}
	orderID, err := u.orderRepo.CreateOrder(ctx, userID, cart.Total)
	if err != nil {
		u.logger.Error("failed to create order", zap.Error(err))
		return "", err
	}
	orderUUID, err := uuid.ParseBytes([]byte(orderID))
	if err != nil {
		u.logger.Error("failed to parse order id", zap.Error(err))
		return "", err
	}

	orderItems := make([]order_entities.OrderItem, 0, len(cart.Items))
	for _, item := range cart.Items {
		orderItems = append(orderItems, order_entities.OrderItem{
			OrderID:   orderUUID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.TotalPrice,
		})
	}
	if err := u.orderRepo.CreateOrderItemsBatch(ctx, orderItems); err != nil {
		u.logger.Error("failed to create order items", zap.Error(err))
		return "", err
	}
	if err := u.cartRepo.DeleteCart(ctx, userID); err != nil {
		u.logger.Error("failed to delete cart", zap.Error(err))
		return "", err
	}
	return orderID, nil
}

func (u *OrderUsecaseImpl) GetOrder(ctx context.Context, orderID string) (order_entities.Order, error) {
	return u.orderRepo.GetOrder(ctx, orderID)
}

func (u *OrderUsecaseImpl) GetOrders(ctx context.Context, userID string) ([]order_entities.Order, error) {
	return u.orderRepo.GetOrdersOfUser(ctx, userID)
}

func (u *OrderUsecaseImpl) CancelOrder(ctx context.Context, orderID string) error {
	return u.orderRepo.UpdateOrderStatus(ctx, orderID, order_entities.OrderStatusCancelled)
}

func (u *OrderUsecaseImpl) ChangeOrderStatus(ctx context.Context, orderID string, status order_entities.OrderStatus) error {
	return u.orderRepo.UpdateOrderStatus(ctx, orderID, status)
}
