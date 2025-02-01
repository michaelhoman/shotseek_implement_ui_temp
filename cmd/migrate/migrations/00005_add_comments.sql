-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS comments (
    id bigserial PRIMARY KEY,
    post_id BIGINT NOT NULL, 
    user_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
ALTER TABLE comments
    ADD CONSTRAINT fk_post FOREIGN KEY (post_id) REFERENCES posts(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE comments
    DROP CONSTRAINT fk_post;
DROP TABLE IF EXISTS comments;


-- +goose StatementEnd
