package product_interface

import (
	"context"
	product_entities "final/internal/features/product/entities"
)

type (
	ProductRepo interface {
		// authorized
		GetProductByIDForUser(ctx context.Context, id string, userID string) (product_entities.Product, error)
		GetAllProductsForUser(ctx context.Context, userID string) ([]product_entities.Product, error)

		GetProductByID(ctx context.Context, id string) (product_entities.Product, error)
		GetAllProducts(ctx context.Context) ([]product_entities.Product, error)
		CreateProduct(ctx context.Context, product product_entities.ProductRequest) (product_entities.Product, error)
		UpdateProduct(ctx context.Context, id string, product product_entities.ProductRequest) (product_entities.Product, error)
		DeleteProduct(ctx context.Context, id string) error
		LikeProduct(ctx context.Context, userID, productID string) error
		UnlikeProduct(ctx context.Context, userID, productID string) error
		// GetLikedProducts(ctx context.Context, userID string) ([]product_entities.Product, error)
	}

	ProductUseCase interface {
		GetProductByID(ctx context.Context, id string, userID ...string) (product_entities.Product, error)
		GetAllProducts(ctx context.Context, userID ...string) ([]product_entities.Product, error)
		CreateProduct(ctx context.Context, product product_entities.ProductRequest) (product_entities.Product, error)
		UpdateProduct(ctx context.Context, id string, product product_entities.ProductRequest) (product_entities.Product, error)
		DeleteProduct(ctx context.Context, id string) error
		LikeProduct(ctx context.Context, userID, productID string) error
		UnlikeProduct(ctx context.Context, userID, productID string) error
		// GetLikedProducts(ctx context.Context, userID string) ([]product_entities.Product, error)
	}
)
