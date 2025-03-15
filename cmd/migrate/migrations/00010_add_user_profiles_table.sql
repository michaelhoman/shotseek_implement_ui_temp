-- +goose Up
-- +goose StatementBegin
CREATE TABLE images (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE, -- Owner of the image
    file_url VARCHAR NOT NULL, -- URL of the stored image
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE images;
-- +goose StatementEnd
