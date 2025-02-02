package store

import (
	"context"
	"database/sql"
	"errors"
	"log"
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
	Comments  []Comment `json:"comments"`
}
type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
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

func (s *PostStore) DeleteByID(ctx context.Context, postID int64) error {
	query := `
	DELETE FROM posts
	WHERE id = $1
	`
	result, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostStore) GetByID(ctx context.Context, postID int64) (*Post, error) {
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
			log.Printf("Error fetching post: %v", err) //TODO: CHECK LOGGING PROCEDURE Or use structured logging
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
	UPDATE posts
	SET title = $1, content = $2, tags = $3
	WHERE id = $4
	RETURNING updated_at
	`

	_, err := s.db.ExecContext(ctx, query, post.Title, post.Content, pq.Array(post.Tags), post.ID)
	if err != nil {
		return err
	}
	return nil

}
