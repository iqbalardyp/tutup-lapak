-- Create table activities
CREATE TABLE purchases (
    id BIGSERIAL PRIMARY KEY,
    total_price INT NOT NULL,
    total_transfer INT NOT NULL,
    sender_name VARCHAR(255) NOT NULL,
    sender_contact_type VARCHAR(255) NOT NULL,
    sender_contact_detail VARCHAR(255) NOT NULL,
    paid_at TIMESTAMPTZ
);

-- Create indices
CREATE INDEX idx_purchase_id ON purchases(id);
-- more index otw