-- Create table pivot_purchase_products
CREATE TABLE pivot_purchase_products (
    id BIGSERIAL PRIMARY KEY,
    purchase_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    qty BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (purchase_id) REFERENCES purchases(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- Create indices
CREATE INDEX idx_pivot_purchase_products_id ON pivot_purchase_products(id);
CREATE INDEX idx_pivot_purchase_products_purchase_id ON pivot_purchase_products(purchase_id);
CREATE INDEX idx_pivot_purchase_products_product_id ON pivot_purchase_products(product_id);
CREATE INDEX idx_pivot_purchase_products_created_at ON pivot_purchase_products(created_at);