-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email citext UNIQUE NOT NULL,
    password bytea NOT NULL,
    zip_code VARCHAR(12) NOT NULL,
    city VARCHAR(255) NOT NULL,
    state VARCHAR(255) NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

INSERT INTO users (first_name, last_name, email, password, zip_code, city, state)
VALUES ('John', 'Doe', 'john.doe@example.com', '\\x1234567890abcdef', '12345', 'Anytown', 'Anystate');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE email = 'john.doe@example.com'; 
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
