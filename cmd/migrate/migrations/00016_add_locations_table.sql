-- +goose Up
-- +goose StatementBegin
CREATE TABLE locations (
    id BIGSERIAL PRIMARY KEY,
    street VARCHAR(255),  -- Optional if only city/state is needed
    city VARCHAR(100) NOT NULL,
    state VARCHAR(50),  -- Can be ENUM for US states
    zip_code VARCHAR(10) NOT NULL,
    country VARCHAR(50) NOT NULL DEFAULT 'USA',
    latitude DECIMAL(9,6),  -- For geolocation-based searches
    longitude DECIMAL(9,6),  -- For geolocation-based searches
    is_precise BOOLEAN DEFAULT FALSE,  -- Indicates if the location is precise
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add location_id to users table
-- ALTER TABLE users
--     ADD COLUMN location_id BIGINT;

-- Add foreign key constraint between users.location_id and locations.id
-- ALTER TABLE users
--     ADD CONSTRAINT fk_location FOREIGN KEY (location_id) REFERENCES locations(id);

-- Drop the unused columns (zip_code, city, state)
-- ALTER TABLE users
--     DROP COLUMN zip_code,
--     DROP COLUMN city,
--     DROP COLUMN state;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- DROP TABLE IF EXISTS locations CASCADE;
-- -- Drop the foreign key constraint
-- ALTER TABLE users
--     DROP CONSTRAINT IF EXISTS fk_location;
-- +goose StatementEnd
