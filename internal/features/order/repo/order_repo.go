package order_repo

import (
	"context"
	"database/sql"
	order_entities "final/internal/features/order/entities"
	order_interface "final/internal/features/order/interface"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type OrderRepoImpl struct {
	logger *zap.Logger
	db     *sqlx.DB
}

var _ order_interface.OrderRepo = (*OrderRepoImpl)(nil)

func New(logger *zap.Logger, db *sqlx.DB) order_interface.OrderRepo {
	return &OrderRepoImpl{
		logger: logger,
		db:     db,
	}
}

func (r *OrderRepoImpl) CreateOrder(ctx context.Context, userID string, totalAmount float64) (string, error) {
	const q = `INSERT INTO orders (user_id, total_amount) VALUES ($1, $2) RETURNING id`

	var orderID string
	if err := r.db.GetContext(ctx, &orderID, q, userID, totalAmount); err != nil {
		r.logger.Error("failed to create order", zap.Error(err))
		return "", err
	}
	return orderID, nil
}

func (r *OrderRepoImpl) CreateOrderItem(ctx context.Context, orderID, productID string, quantity int, price float64) (string, error) {
	const q = `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4) RETURNING id`

	var orderItemID string
	if err := r.db.GetContext(ctx, &orderItemID, q, orderID, productID, quantity, price); err != nil {
		r.logger.Error("failed to create order item", zap.Error(err))
		return "", err
	}
	return orderItemID, nil
}

func (r *OrderRepoImpl) CreateOrderItemsBatch(ctx context.Context, orderItems []order_entities.OrderItem) error {
	const q = `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (:order_id, :product_id, :quantity, :price)`

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.,
	})
	if err != nil {
		r.logger.Error("failed to begin transaction", zap.Error(err))
		return err
	}

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		r.logger.Error("failed to prepare statement", zap.Error(err))
		return err
	}
	defer stmt.Close()

	for _, item := range orderItems {
		if _, err := stmt.Exec(item); err != nil {
			tx.Rollback()
			r.logger.Error("failed to execute statement", zap.Error(err))
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (r *OrderRepoImpl) UpdateOrderStatus(ctx context.Context, orderID string, status int) error {
	const q = `UPDATE orders SET status = $1 WHERE id = $2`
	if _, err := r.db.ExecContext(ctx, q, status, orderID); err != nil {
		r.logger.Error("failed to update order status", zap.Error(err))
		return err
	}
	return nil
}

func (r *OrderRepoImpl) DeleteOrder(ctx context.Context, orderID string) error {
	const q = `DELETE FROM orders WHERE id = $1`
	if _, err := r.db.ExecContext(ctx, q, orderID); err != nil {
		r.logger.Error("failed to delete order", zap.Error(err))
		return err
	}
	return nil
}

func (r *OrderRepoImpl) GetOrder(ctx context.Context, orderID string) (order_entities.Order, error) {
	const q = `SELECT * FROM orders WHERE id = $1`

	var order order_entities.Order
	if err := r.db.GetContext(ctx, &order, q, orderID); err != nil {
		r.logger.Error("failed to get order", zap.Error(err))
		return order, err
	}
	return order, nil
}

func (r *OrderRepoImpl) GetOrdersOfUser(ctx context.Context, userID string) ([]order_entities.Order, error) {
	const q = `SELECT * FROM orders WHERE user_id = $1`

	var orders []order_entities.Order
	if err := r.db.SelectContext(ctx, &orders, q, userID); err != nil {
		r.logger.Error("failed to get orders of user", zap.Error(err))
		return nil, err
	}
	return orders, nil
}
