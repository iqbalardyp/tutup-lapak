-- Create table activities
CREATE TABLE purchases (
    id BIGSERIAL PRIMARY KEY,
    total_price INT,
    sender_name VARCHAR(255),
    sender_contact_type VARCHAR(255),
    sender_contact_detail VARCHAR(255)
);

-- Create indices
CREATE INDEX idx_purchase_id ON purchases(id);
-- more index otw