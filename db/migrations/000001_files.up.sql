-- Create table files
CREATE TABLE files (
    id BIGSERIAL PRIMARY KEY,
    uri VARCHAR(255),
    thumbnail_uri VARCHAR(255)
);

-- Create indexes
CREATE INDEX idx_files_id ON files(id);