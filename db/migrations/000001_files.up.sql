-- Create table files
CREATE TABLE files (
    id BIGSERIAL PRIMARY KEY,
    uri VARCHAR(255) NOT NULL,
    thumbnail_uri VARCHAR(255) NOT NULL
);

-- Create indexes
CREATE INDEX idx_files_id ON files(id);