package cart_repo

import (
	"context"
	"database/sql"
	pg_errors "final/common/kerrors"
	cart_entities "final/internal/features/cart/entities"
	cart_interface "final/internal/features/cart/interface"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type CartRepoImpl struct {
	logger *zap.Logger
	db     *sqlx.DB
}

var _ cart_interface.CartRepo = (*CartRepoImpl)(nil)

func New(logger *zap.Logger, db *sqlx.DB) cart_interface.CartRepo {
	return &CartRepoImpl{
		logger: logger,
		db:     db,
	}
}

func (r *CartRepoImpl) GetCart(ctx context.Context, userID string) (*cart_entities.Cart, error) {
	// check if user already has a cart
	const queryCart = `SELECT id FROM user_carts WHERE user_id = $1`
	cart := &cart_entities.Cart{Items: make([]cart_entities.CartItem, 0)}
	if err := r.db.GetContext(ctx, &cart.ID, queryCart, userID); err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrap(err, "failed to check if user already has a cart")
	}

	if cart.ID != uuid.Nil {
		const queryCartItems = `
		SELECT ci.id, ci.product_id, p.name AS product_name, p.price AS product_price, quantity, quantity * p.price AS total_price
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
	} else {
		const queryInsertCart = `INSERT INTO user_carts (user_id) VALUES ($1) RETURNING id`
		if err := r.db.GetContext(ctx, &cart.ID, queryInsertCart, userID); err != nil {
			return nil, errors.Wrap(err, "failed to insert cart")
		}

		cart.Items = []cart_entities.CartItem{}

		return cart, nil
	}
}

func calculateTotal(items []cart_entities.CartItem) float64 {
	var total float64
	for _, item := range items {
		total += item.ProductPrice * float64(item.Quantity)
	}

	return total
}

func (r *CartRepoImpl) AddToCart(ctx context.Context, userID, productID string, quantity int) error {
	// get cart id
	cart, err := r.GetCart(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "failed to get cart")
	}

	// check if product already in cart
	const queryCheckProduct = `SELECT id FROM cart_items WHERE cart_id = $1 AND product_id = $2`
	var cartItem cart_entities.CartItem
	if err := r.db.GetContext(ctx, &cartItem.ID, queryCheckProduct, cart.ID, productID); err != nil && err != sql.ErrNoRows {
		return errors.Wrap(err, "failed to check if product already in cart")
	}

	if cartItem.ID != uuid.Nil {
		// if product already in cart, update quantity
		const queryUpdateQuantity = `UPDATE cart_items SET quantity = quantity + $1 WHERE id = $2`
		if _, err := r.db.ExecContext(ctx, queryUpdateQuantity, quantity, cartItem.ID); err != nil {
			return errors.Wrap(err, "failed to update quantity")
		}
		return nil
	} else {
		// if product not in cart, insert product
		const queryInsertProduct = `
		INSERT INTO cart_items (cart_id, product_id, quantity) 
		VALUES ($1, $2, $3)`
		if _, err := r.db.ExecContext(ctx, queryInsertProduct, cart.ID, productID, quantity); err != nil {
			return errors.Wrap(err, "failed to insert product")
		}

		return nil
	}
}

func (r *CartRepoImpl) RemoveFromCart(ctx context.Context, userID, productID string) error {
	const q = `
	DELETE FROM cart_items
	WHERE cart_id = (SELECT id FROM user_carts WHERE user_id = $1) AND product_id = $2 
	RETURNING (SELECT count(*) FROM cart_items WHERE cart_id = (SELECT id FROM user_carts WHERE user_id = $1))
	`

	var count int
	if err := r.db.GetContext(ctx, &count, q, userID, productID); err != nil {
		return errors.Wrap(err, "failed to remove product from cart")
	}

	if count == 0 {
		const q = `DELETE FROM user_carts WHERE user_id = $1`
		if _, err := r.db.ExecContext(ctx, q, userID); err != nil {
			return errors.Wrap(err, "failed to remove cart")
		}
		r.logger.Info("cart removed")
	}

	return nil
}

func (r *CartRepoImpl) UpdateCart(ctx context.Context, userID, productID string, quantity int) error {
	if quantity == 0 {
		return r.RemoveFromCart(ctx, userID, productID)
	}

	const q = `
	UPDATE cart_items SET quantity = $3 
	WHERE cart_id = (SELECT id FROM user_carts WHERE user_id = $1) AND product_id = $2`

	res, err := r.db.ExecContext(ctx, q, userID, productID, quantity)
	if err != nil {
		return errors.Wrap(err, "failed to update cart")
	}

	rows, err := res.RowsAffected()
	if rows < 1 {
		return pg_errors.NewPgError(pg_errors.PgErrorNoRowsAffected)
	}

	return nil
}

func (r *CartRepoImpl) DeleteCart(ctx context.Context, userID string) error {
	const queryCart = `SELECT id FROM user_carts WHERE user_id = $1`

	var cartID string

	if err := r.db.GetContext(ctx, &cartID, queryCart, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return errors.Wrap(err, "failed to check if user already has a cart")
	}

	const q = `DELETE FROM user_carts WHERE user_id = $1`
	if _, err := r.db.ExecContext(ctx, q, userID); err != nil {
		return errors.Wrap(err, "failed to delete cart")
	}

	return nil
}
