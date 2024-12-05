package product_repo

import (
	"context"
	product_entities "final/internal/features/product/entities"
	product_interface "final/internal/features/product/interface"
	"final/internal/infrastructure/db_gen"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ProductRepoImpl struct {
	logger *zap.Logger
	db     *db_gen.Queries
}

var _ product_interface.ProductRepo = (*ProductRepoImpl)(nil)

func New(logger *zap.Logger, db *db_gen.Queries) *ProductRepoImpl {
	return &ProductRepoImpl{
		logger: logger,
		db:     db,
	}
}

func (r *ProductRepoImpl) GetProductByID(ctx context.Context, id string) (product_entities.Product, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		r.logger.Error("error while parsing id", zap.Error(err))
		return product_entities.Product{}, err
	}
	dbProduct, err := r.db.GetProductById(ctx, uuid)
	if err != nil {
		r.logger.Error("error while getting product by id", zap.Error(err))
		return product_entities.Product{}, err
	}
	return product_entities.FromDBProduct(dbProduct), nil
}
func (r *ProductRepoImpl) GetAllProducts(ctx context.Context) ([]product_entities.Product, error) {
	dbProducts, err := r.db.GetProducts(ctx)
	if err != nil {
		r.logger.Error("error while getting products", zap.Error(err))
		return nil, err
	}
	return product_entities.FromDBProducts(dbProducts), nil
}
func (r *ProductRepoImpl) CreateProduct(ctx context.Context, product product_entities.ProductRequest) (product_entities.Product, error) {
	dbProduct, err := r.db.CreateProduct(ctx, db_gen.CreateProductParams{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	})
	if err != nil {
		r.logger.Error("error while creating product", zap.Error(err))
		return product_entities.Product{}, err
	}
	return product_entities.FromDBProduct(dbProduct), nil
}
func (r *ProductRepoImpl) UpdateProduct(ctx context.Context, id string, product product_entities.ProductRequest) (product_entities.Product, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		r.logger.Error("error while parsing id", zap.Error(err))
		return product_entities.Product{}, err
	}
	dbProduct, err := r.db.UpdateProduct(ctx, db_gen.UpdateProductParams{
		ID:          uuid,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	})
	if err != nil {
		r.logger.Error("error while updating product", zap.Error(err))
		return product_entities.Product{}, err
	}
	return product_entities.FromDBProduct(dbProduct), nil

}
func (r *ProductRepoImpl) DeleteProduct(ctx context.Context, id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		r.logger.Error("error while parsing id", zap.Error(err))
		return err
	}
	_, err = r.db.DeleteProduct(ctx, uuid)
	if err != nil {
		r.logger.Error("error while deleting product", zap.Error(err))
		return err
	}
	return nil
}
