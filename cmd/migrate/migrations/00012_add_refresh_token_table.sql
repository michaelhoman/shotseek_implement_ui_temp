-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_email citext NOT NULL,
    token_hash TEXT NOT NULL UNIQUE,  -- Ensure each token is unique
    stored_fp TEXT NOT NULL,  -- Store the fingerprint of the device
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_email) REFERENCES users(email)  -- Link to the users table
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS refresh_tokens;
-- +goose StatementEnd
