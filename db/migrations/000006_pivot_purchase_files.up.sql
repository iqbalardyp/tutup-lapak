-- Create table pivot_purchase_products
CREATE TABLE pivot_purchase_files (
    id BIGSERIAL PRIMARY KEY,
    purchase_id BIGINT NOT NULL,
    file_id BIGINT,
    FOREIGN KEY (purchase_id) REFERENCES purchases(id) ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);


-- Create indices
CREATE INDEX idx_pivot_purchase_files_id ON pivot_purchase_files(id);
CREATE INDEX idx_pivot_purchase_files_purchase_id ON pivot_purchase_files(purchase_id);
CREATE INDEX idx_pivot_purchase_files_file_id ON pivot_purchase_files(file_id);