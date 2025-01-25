-- Create table users
CREATE TABLE sellers (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone_number VARCHAR(255) UNIQUE NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    file_id BIGINT,
    bank_account_name VARCHAR(255),
    bank_account_holder VARCHAR(255),
    bank_account_number VARCHAR(255),
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_sellers_id ON sellers(id);
CREATE INDEX idx_sellers_email ON sellers(email);
CREATE INDEX idx_sellers_phone_number ON sellers(phone_number);
