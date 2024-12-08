package product_repo

import (
	"context"
	"encoding/json"
	product_entities "final/internal/features/product/entities"
	product_interface "final/internal/features/product/interface"
	"time"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type ProductRepoImpl struct {
	logger *zap.Logger
	db     *sqlx.DB
	cache  *redis.Client
}

var _ product_interface.ProductRepo = (*ProductRepoImpl)(nil)

func New(logger *zap.Logger, db *sqlx.DB, cache *redis.Client) *ProductRepoImpl {
	return &ProductRepoImpl{
		logger: logger,
		db:     db,
		cache:  cache,
	}
}

func (r *ProductRepoImpl) LikeProduct(ctx context.Context, userID, productID string) error {
	const q = `INSERT INTO user_saveds (user_id, product_id) VALUES ($1, $2)`
	if _, err := r.db.ExecContext(ctx, q, userID, productID); err != nil {
		r.logger.Error("error while liking product", zap.Error(err))
		return err
	}
	return nil
}

func (r *ProductRepoImpl) UnlikeProduct(ctx context.Context, userID, productID string) error {
	const q = `DELETE FROM user_saveds WHERE user_id = $1 AND product_id = $2`
	if _, err := r.db.ExecContext(ctx, q, userID, productID); err != nil {
		r.logger.Error("error while unliking product", zap.Error(err))
		return err
	}
	return nil
}

func (r *ProductRepoImpl) GetProductByIDForUser(ctx context.Context, id string, userID string) (product_entities.Product, error) {
	var dbProduct product_entities.Product

	const q = `
	SELECT p.id, p.name, p.description, p.price, p.created_at, p.updated_at, p.deleted_at, (SELECT id FROM user_saveds WHERE user_id = $2 AND product_id = $1) IS NOT NULL AS is_liked
	FROM products p 
	LEFT JOIN user_saveds s ON s.product_id = p.id 
	LEFT JOIN users u ON u.id = s.user_id
	WHERE p.id = $1`

	if err := r.db.GetContext(ctx, &dbProduct, q, id, userID); err != nil {
		r.logger.Error("error while getting product by id for user", zap.Error(err))
		return dbProduct, err
	}

	return dbProduct, nil
}

func (r *ProductRepoImpl) GetAllProductsForUser(ctx context.Context, userID string) ([]product_entities.Product, error) {
	dbProducts := []product_entities.Product{}

	const q = `
		SELECT p.id, p.name, p.description, p.price, p.created_at, p.updated_at, p.deleted_at, (SELECT id FROM user_saveds WHERE user_id = $1 AND product_id = p.id) IS NOT NULL AS is_liked
		FROM products p 
		LEFT JOIN user_saveds s ON s.product_id = p.id 
		LEFT JOIN users u ON u.id = s.user_id
		WHERE deleted_at IS NULL`

	if err := r.db.SelectContext(ctx, &dbProducts, q, userID); err != nil {
		r.logger.Error("error while getting products for user", zap.Error(err))
		return nil, err
	}

	return dbProducts, nil

}

func (r *ProductRepoImpl) GetProductByID(ctx context.Context, id string) (product_entities.Product, error) {

	productJSON, cacheErr := r.cache.Get("product-" + id).Result()
	if cacheErr == nil {
		var product product_entities.Product
		if err := json.Unmarshal([]byte(productJSON), &product); err != nil {
			r.logger.Error("error while unmarshalling product", zap.Error(err))
			return product, err
		}
		return product, nil
	}

	var dbProduct product_entities.Product

	const q = `
	SELECT id, name, description, price, created_at, updated_at, deleted_at
	FROM products p 
	WHERE p.id = $1`

	if err := r.db.GetContext(ctx, &dbProduct, q, id); err != nil {
		r.logger.Error("error while getting product by id", zap.Error(err))
		return dbProduct, err
	}

	if cacheErr == redis.Nil {
		productJSON, err := json.Marshal(dbProduct)
		if err == nil {
			r.cache.Set("product-"+dbProduct.ID.String(), productJSON, time.Minute*5)
		}
	}

	return dbProduct, nil
}

func (r *ProductRepoImpl) GetAllProducts(ctx context.Context) ([]product_entities.Product, error) {
	dbProducts := []product_entities.Product{}

	const q = `
	SELECT id, name, description, price, created_at, updated_at, deleted_at
	FROM products 
	WHERE deleted_at IS NULL`

	if err := r.db.SelectContext(ctx, &dbProducts, q); err != nil {
		r.logger.Error("error while getting products", zap.Error(err))
		return nil, err
	}
	return dbProducts, nil
}

func (r *ProductRepoImpl) CreateProduct(ctx context.Context, product product_entities.ProductRequest) (product_entities.Product, error) {
	const q = `INSERT INTO products (name, description, price) VALUES ($1, $2, $3) RETURNING *`
	var dbProduct product_entities.Product

	if err := r.db.GetContext(ctx, &dbProduct, q, product.Name, product.Description, product.Price); err != nil {
		r.logger.Error("error while creating product", zap.Error(err))
		return dbProduct, err
	}

	productJSON, err := json.Marshal(dbProduct)
	if err == nil {
		err := r.cache.Set("product-"+dbProduct.ID.String(), productJSON, time.Minute*5).Err()
		if err != nil {
			r.logger.Error("error while caching product", zap.Error(err))
		}
	}

	return dbProduct, nil
}
func (r *ProductRepoImpl) UpdateProduct(ctx context.Context, id string, product product_entities.ProductRequest) (product_entities.Product, error) {
	var dbProduct product_entities.Product
	const q = `UPDATE products SET name = $1, description = $2, price = $3 WHERE id = $4 RETURNING *`

	if err := r.db.GetContext(ctx, &dbProduct, q, product.Name, product.Description, product.Price, id); err != nil {
		r.logger.Error("error while updating product", zap.Error(err))
		return dbProduct, err
	}
	return dbProduct, nil

}
func (r *ProductRepoImpl) DeleteProduct(ctx context.Context, id string) error {
	const q = `UPDATE products SET deleted_at = NOW() WHERE id = $1`
	if _, err := r.db.ExecContext(ctx, q, id); err != nil {
		r.logger.Error("error while deleting product", zap.Error(err))
		return err
	}
	return nil
}
