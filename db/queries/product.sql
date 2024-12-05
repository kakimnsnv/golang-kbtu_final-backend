-- name: CreateProduct :one
INSERT INTO products (name, price, description) VALUES ($1, $2, $3) RETURNING *;

-- name: GetProductById :one
SELECT * FROM products WHERE id = $1;

-- name: GetProducts :many
SELECT * FROM products;

-- name: UpdateProduct :one
UPDATE products SET name = $1, price = $2, description = $3 WHERE id = $4 RETURNING *;

-- name: DeleteProduct :one
DELETE FROM products WHERE id = $1 RETURNING *;

