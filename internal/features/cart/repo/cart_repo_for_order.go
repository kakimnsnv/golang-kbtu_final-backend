package cart_repo

import (
	"context"
	"database/sql"
	cart_entities "final/internal/features/cart/entities"
	cart_interface "final/internal/features/cart/interface"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type CartRepoForOrderImpl struct {
	db     *sqlx.DB
	logger *zap.Logger
}

var _ cart_interface.CartRepoForOrderUsecase = (*CartRepoForOrderImpl)(nil)

func NewCartRepoForOrder(db *sqlx.DB, logger *zap.Logger) cart_interface.CartRepoForOrderUsecase {
	return &CartRepoForOrderImpl{
		db:     db,
		logger: logger,
	}
}

func (r *CartRepoForOrderImpl) GetCart(ctx context.Context, userID string) (*cart_entities.Cart, error) {
	const queryCart = `SELECT id FROM user_carts WHERE user_id = $1`
	cart := &cart_entities.Cart{Items: make([]cart_entities.CartItem, 0)}
	if err := r.db.GetContext(ctx, &cart.ID, queryCart, userID); err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrap(err, "failed to check if user already has a cart")
	}

	if cart.ID == uuid.Nil {
		return nil, errors.New("cart not found")
	}

	const queryCartItems = `
		SELECT ci.id, ci.product_id, p.name AS product_name, p.price AS product_price, p.photo AS product_photo, quantity, quantity * p.price AS total_price
		FROM cart_items ci 
		LEFT JOIN products p ON p.id = ci.product_id 
		WHERE cart_id = $1
		`

	if err := r.db.SelectContext(ctx, &cart.Items, queryCartItems, cart.ID); err != nil {
		return nil, errors.Wrap(err, "failed to get cart items")
	}

	if cart.Items == nil {
		cart.Items = []cart_entities.CartItem{}
	}
	cart.Total = calculateTotal(cart.Items)

	return cart, nil
}
func (r *CartRepoForOrderImpl) DeleteCart(ctx context.Context, userID string) error {
	const q = `DELETE FROM user_carts WHERE user_id = $1`

	if _, err := r.db.ExecContext(ctx, q, userID); err != nil {
		r.logger.Error("failed to delete cart", zap.Error(err))
		return err
	}
	return nil
}
