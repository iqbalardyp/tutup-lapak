-- Drop indexes
DROP INDEX IF EXISTS idx_pivot_purchase_files_id;
DROP INDEX IF EXISTS idx_pivot_purchase_files_purchase_id;
DROP INDEX IF EXISTS idx_pivot_purchase_files_file_id;

-- DROP purchases
DROP TABLE IF EXISTS pivot_purchase_files CASCADE;