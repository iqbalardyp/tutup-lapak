-- Drop indexes
DROP INDEX IF EXISTS idx_pivot_purchase_products_id;
DROP INDEX IF EXISTS idx_pivot_purchase_products_purchase_id;
DROP INDEX IF EXISTS idx_pivot_purchase_products_product_id;

-- DROP purchases
DROP TABLE IF EXISTS pivot_purchase_products CASCADE;