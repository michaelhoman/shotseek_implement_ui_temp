-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_images (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    image_url VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_images CASCADE;
-- +goose StatementEnd
