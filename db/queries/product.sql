-- name: CreateProduct :one
INSERT INTO products (name, price, description) VALUES ($1, $2, $3) RETURNING *;

-- name: GetProductById :one
SELECT * FROM products WHERE id = $1;

-- name: GetProducts :many
SELECT * FROM products;