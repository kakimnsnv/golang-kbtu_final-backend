// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: product.sql

package db_gen

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createProduct = `-- name: CreateProduct :one
INSERT INTO products (name, price, description) VALUES ($1, $2, $3) RETURNING id, name, description, price, created_at, updated_at
`

type CreateProductParams struct {
	Name        string         `json:"name"`
	Price       pgtype.Numeric `json:"price"`
	Description pgtype.Text    `json:"description"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error) {
	row := q.db.QueryRow(ctx, createProduct, arg.Name, arg.Price, arg.Description)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProductById = `-- name: GetProductById :one
SELECT id, name, description, price, created_at, updated_at FROM products WHERE id = $1
`

func (q *Queries) GetProductById(ctx context.Context, id uuid.UUID) (Product, error) {
	row := q.db.QueryRow(ctx, getProductById, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProducts = `-- name: GetProducts :many
SELECT id, name, description, price, created_at, updated_at FROM products
`

func (q *Queries) GetProducts(ctx context.Context) ([]Product, error) {
	rows, err := q.db.Query(ctx, getProducts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}