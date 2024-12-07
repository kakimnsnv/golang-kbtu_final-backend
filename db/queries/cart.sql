-- name: GetCart :one
SELECT 
    uc.id AS cart_id, 
    uc.user_id, 
    json_agg(
        json_build_object(
            'id', ci.id,
            'product_id', ci.product_id,
            'quantity', ci.quantity,
            'product_name', p.name,
            'product_price', p.price
        )
    ) AS items
FROM user_carts uc
LEFT JOIN cart_items ci ON uc.id = ci.cart_id
LEFT JOIN products p ON ci.product_id = p.id
WHERE uc.user_id = @user_id
GROUP BY uc.id, uc.user_id;

-- name: AddToCart :one
INSERT INTO cart_items (cart_id, product_id, quantity)
VALUES (
    (SELECT id FROM user_carts WHERE user_id = @user_id),
    @product_id, 
    @quantity
)
ON CONFLICT (cart_id, product_id) DO UPDATE 
SET quantity = cart_items.quantity + @quantity, 
    updated_at = NOW()
RETURNING *;

-- name: RemoveFromCart :exec
DELETE FROM cart_items 
WHERE cart_id = (SELECT id FROM user_carts WHERE user_id = @user_id) 
AND product_id = @product_id;

-- name: UpdateCartItemQuantity :exec
UPDATE cart_items 
SET quantity = @new_quantity, 
    updated_at = NOW()
WHERE cart_id = (SELECT id FROM user_carts WHERE user_id = @user_id)
AND product_id = @product_id;