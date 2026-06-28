CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    price_in_cents BIGINT NOT NULL,
    type VARCHAR(20) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

-- Seed initial data
INSERT INTO products (id, sku, name, price_in_cents, type, stock, is_active) 
VALUES
    ('11111111-1111-1111-1111-111111111111', 'BOOK-001', 'E-Commerce done right', 2999, 'BOOK', 100, true),
    ('22222222-2222-2222-2222-222222222222', 'LIC-001', 'E-Commerce simulator', 9900, 'LICENSE', 0, true),
    ('33333333-3333-3333-3333-333333333333', 'VCH-001', 'Gift Voucher', 5000, 'VOUCHER', 0, true)
ON CONFLICT (id) DO NOTHING;