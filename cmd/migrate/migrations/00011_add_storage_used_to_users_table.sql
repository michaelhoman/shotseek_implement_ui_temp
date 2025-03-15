-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN storage_used BIGINT DEFAULT 0;
ALTER TABLE users ADD COLUMN storage_limit BIGINT DEFAULT 104857600; -- 100MB default

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN storage_used;
ALTER TABLE users DROP COLUMN storage_limit;
-- +goose StatementEnd
