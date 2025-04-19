-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email citext UNIQUE NOT NULL,
    password bytea NOT NULL,
    location_id BIGINT,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

ALTER TABLE users
ALTER COLUMN id SET DATA TYPE UUID USING id::UUID;

-- For gen_random_uuid():

INSERT INTO users (id, first_name, last_name, email, password, location_id)
VALUES (gen_random_uuid(),'John', 'Doe', 'john.doe@example.com', '\\x1234567890abcdef', '10080800');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE email = 'john.doe@example.com'; 
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
