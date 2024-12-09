package order_interface

import (
	"context"
	order_entities "final/internal/features/order/entities"
)

// CREATE TABLE orders (
//     id uuid.UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
//     user_id uuid.UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
//     total_amount DECIMAL(10,2) NOT NULL,
//     status INT NOT NULL DEFAULT 0,
//     created_at TIMESTAMP NOT NULL DEFAULT NOW()
// );

// CREATE TABLE order_items (
//
//	id uuid.UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
//	order_id uuid.UUId NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
//	product_id uuid.UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
//	quantity INTEGER NOT NULL,
//	price DECIMAL(10,2) NOT NULL,
//	created_at TIMESTAMP NOT NULL DEFAULT NOW()
//
// );
type (
	OrderRepo interface {
		CreateOrder(ctx context.Context, userID string, totalAmount float64) (string, error)
		CreateOrderItem(ctx context.Context, orderID, productID string, quantity int, price float64) (string, error)
		UpdateOrderStatus(ctx context.Context, orderID string, status order_entities.OrderStatus) error
		DeleteOrder(ctx context.Context, orderID string) error
		GetOrder(ctx context.Context, orderID string) (order_entities.Order, error)
		GetOrdersOfUser(ctx context.Context, userID string) ([]order_entities.Order, error)
		CreateOrderItemsBatch(ctx context.Context, orderItems []order_entities.OrderItem) error
	}

	OrderUsecase interface {
		PlaceOrder(ctx context.Context, userID string) (string, error)
		GetOrder(ctx context.Context, orderID string) (order_entities.Order, error)
		GetOrders(ctx context.Context, userID string) ([]order_entities.Order, error)
		CancelOrder(ctx context.Context, orderID string) error
		ChangeOrderStatus(ctx context.Context, orderID string, status order_entities.OrderStatus) error
	}
)
