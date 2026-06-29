CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,
    customer_email VARCHAR(255) NOT NULL,
    total_price BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    quantity INT NOT NULL,
    price BIGINT NOT NULL,
    generated_asset VARCHAR(255) NOT NULL
);

-- Index foreign keys to guarantee high retrieval performance under heavy load
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);