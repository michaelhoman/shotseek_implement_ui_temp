-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN version INT NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN version;
ALTER TABLE users DROP COLUMN updated_at;
-- +goose StatementEnd
