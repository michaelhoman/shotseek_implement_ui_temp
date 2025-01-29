package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	Tags      []string  `json:"tags"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type PostgresPostStore struct {
	db *sql.DB
}

func (s *PostgresPostStore) Create(ctx context.Context, post *Post) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := `
INSERT INTO posts (content, title, tags, user_id)
VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
`
	if post.Tags == nil {
		post.Tags = []string{}
	}
	err = tx.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		pq.Array(&post.Tags),
		post.UserID,
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *PostgresPostStore) GetByID(ctx context.Context, postID int64) (*Post, error) {
	query := `
	SELECT id, content, title, tags, user_id, created_at, updated_at 
	FROM posts 
	WHERE id = $1 
	LIMIT 1`
	post := Post{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		postID,
	).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		pq.Array(&post.Tags),
		&post.UserID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}
