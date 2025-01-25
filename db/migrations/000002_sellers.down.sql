-- Drop indexes
DROP INDEX IF EXISTS idx_sellers_id;
DROP INDEX IF EXISTS idx_sellers_email;
DROP INDEX IF EXISTS idx_sellers_phone_number;

-- DROP users
DROP TABLE IF EXISTS sellers CASCADE;