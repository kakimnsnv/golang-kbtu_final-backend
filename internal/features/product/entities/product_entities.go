package product_entities

import (
	"final/internal/infrastructure/db_gen"
	"time"

	"github.com/google/uuid"
)

type (
	Product struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Price       float64   `json:"price"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	ProductRequest struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}
)

func FromDBProduct(dbProduct db_gen.Product) Product {
	return Product{
		ID:          dbProduct.ID,
		Name:        dbProduct.Name,
		Description: dbProduct.Description,
		Price:       dbProduct.Price,
		CreatedAt:   dbProduct.CreatedAt,
		UpdatedAt:   dbProduct.UpdatedAt,
	}
}

func FromDBProducts(dbProducts []db_gen.Product) []Product {
	products := make([]Product, 0, len(dbProducts))
	for _, dbProduct := range dbProducts {
		products = append(products, FromDBProduct(dbProduct))
	}
	return products
}
