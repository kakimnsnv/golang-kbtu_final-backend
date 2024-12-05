ALTER TABLE products
    ALTER COLUMN created_at SET NOT NULL;

ALTER TABLE products
    ALTER COLUMN updated_at SET NOT NULL;

ALTER TABLE user_carts
    ALTER COLUMN created_at SET NOT NULL;

ALTER TABLE user_carts
    ALTER COLUMN updated_at SET NOT NULL;

ALTER TABLE cart_items
    ALTER COLUMN created_at SET NOT NULL;

ALTER TABLE cart_items 
    ALTER COLUMN updated_at SET NOT NULL;

ALTER TABLE users
    ALTER COLUMN created_at SET NOT NULL;

ALTER TABLE users
    ALTER COLUMN updated_at SET NOT NULL;