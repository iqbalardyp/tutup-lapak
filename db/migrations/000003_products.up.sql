-- Create enum
CREATE TYPE enum_product_categories as ENUM (
    'Food',
    'Beverage',
    'Clothes',
    'Furniture',
    'Tools'
);

-- Create table users
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    seller_id BIGINT NOT NULL,
    file_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    category enum_product_categories NOT NULL,
    qty INT NOT NULL,
    price INT NOT NULL,
    sku VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_sold_at TIMESTAMPTZ,
    FOREIGN KEY (seller_id) REFERENCES sellers(id) ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers
CREATE TRIGGER set_timestamp_products
    BEFORE UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

-- Create indexes
CREATE INDEX idx_products_id ON products(id);
CREATE INDEX idx_products_seller_id ON products(seller_id);
CREATE INDEX idx_products_file_id ON products(file_id);
CREATE INDEX idx_products_newest ON products(GREATEST(created_at, updated_at));
-- more index otw