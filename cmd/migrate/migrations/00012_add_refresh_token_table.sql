-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id uuid NOT NULL,
    token_hash TEXT NOT NULL UNIQUE,  -- Ensure each token is unique
    stored_fp TEXT NOT NULL,  -- Store the fingerprint of the device
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id)  -- Link to the users table
);

-- Add a composite unique constraint on user_id and stored_fp
ALTER TABLE refresh_tokens ADD CONSTRAINT unique_user_device UNIQUE (user_id, stored_fp);
-- Add an index on the token_hash for faster lookups
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop the unique constraint if it exists
ALTER TABLE refresh_tokens DROP CONSTRAINT IF EXISTS unique_user_device;

DROP TABLE IF EXISTS refresh_tokens;
-- Drop the index if it exists
DROP INDEX IF EXISTS idx_refresh_tokens_token_hash;

-- +goose StatementEnd
