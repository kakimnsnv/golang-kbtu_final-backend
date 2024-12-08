package product_usecase

import (
	"context"
	product_entities "final/internal/features/product/entities"
	product_interface "final/internal/features/product/interface"

	"go.uber.org/zap"
)

type ProductUseCaseImpl struct {
	logger *zap.Logger
	repo   product_interface.ProductRepo
}

var _ product_interface.ProductUseCase = (*ProductUseCaseImpl)(nil)

func New(logger *zap.Logger, repo product_interface.ProductRepo) *ProductUseCaseImpl {
	return &ProductUseCaseImpl{
		logger: logger,
		repo:   repo,
	}
}

func (u *ProductUseCaseImpl) GetProductByID(ctx context.Context, id string, userID ...string) (product_entities.Product, error) {
	if len(userID) > 0 {
		return u.repo.GetProductByIDForUser(ctx, id, userID[0])
	}
	return u.repo.GetProductByID(ctx, id)
}

func (u *ProductUseCaseImpl) GetAllProducts(ctx context.Context, userID ...string) ([]product_entities.Product, error) {
	if len(userID) > 0 {
		return u.repo.GetAllProductsForUser(ctx, userID[0])
	}
	return u.repo.GetAllProducts(ctx)
}

func (u *ProductUseCaseImpl) CreateProduct(ctx context.Context, product product_entities.ProductRequest) (product_entities.Product, error) {
	return u.repo.CreateProduct(ctx, product)
}

func (u *ProductUseCaseImpl) UpdateProduct(ctx context.Context, id string, product product_entities.ProductRequest) (product_entities.Product, error) {
	return u.repo.UpdateProduct(ctx, id, product)
}

func (u *ProductUseCaseImpl) DeleteProduct(ctx context.Context, id string) error {
	return u.repo.DeleteProduct(ctx, id)
}

func (u *ProductUseCaseImpl) LikeProduct(ctx context.Context, userID, productID string) error {
	return u.repo.LikeProduct(ctx, userID, productID)
}

func (u *ProductUseCaseImpl) UnlikeProduct(ctx context.Context, userID, productID string) error {
	return u.repo.UnlikeProduct(ctx, userID, productID)
}
