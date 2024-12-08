CREATE TABLE orders (
    id uuid.UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid.UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    total_amount DECIMAL(10,2) NOT NULL,
    status INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE order_items (
    id uuid.UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    order_id uuid.UUId NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id uuid.UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);