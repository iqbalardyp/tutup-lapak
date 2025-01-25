-- Drop indexes
DROP INDEX IF EXISTS idx_products_id;
DROP INDEX IF EXISTS idx_products_seller_id;
DROP INDEX IF EXISTS idx_products_file_id;

-- DROP trigger
DROP TRIGGER IF EXISTS set_timestamp_products ON products CASCADE;
DROP FUNCTION IF EXISTS trigger_set_timestamp CASCADE;

-- DROP users
DROP TABLE IF EXISTS products CASCADE;

-- DROP enum
DROP TYPE IF EXISTS enum_product_categories CASCADE;
