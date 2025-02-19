-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN is_active;
-- +goose StatementEnd
