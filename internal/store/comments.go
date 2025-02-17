package store

import (
	"context"
	"database/sql"
	"fmt"
)

type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	User      User   `json:"user"`
}

type CommentsStore struct {
	db *sql.DB
}

func (s *CommentsStore) Create(ctx context.Context, comment *Comment) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := `
INSERT INTO comments (post_id, user_id, content)
VALUES ($1, $2, $3) RETURNING id, created_at, updated_at
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err = tx.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
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

func (s *CommentsStore) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	query := `
	SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, c.updated_at, u.first_name, u.last_name
	FROM comments c
	JOIN users u ON c.user_id = u.id
	WHERE c.post_id = $1
	ORDER BY c.created_at DESC;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := []Comment{}
	for rows.Next() {
		var c Comment
		c.User = User{}
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.User.FirstName, &c.User.LastName)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (s *CommentsStore) GetByCommentID(ctx context.Context, commentID int64) (*Comment, error) {
	query := `
	SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, c.updated_at, u.first_name, u.last_name
	FROM comments c
	JOIN users u ON c.user_id = u.id
	WHERE c.id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	comment := Comment{}
	comment.User = User{}
	err := s.db.QueryRowContext(ctx, query, commentID).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt, &comment.User.FirstName, &comment.User.LastName)
	if err != nil {
		return nil, err
	}
	fmt.Println("KILL RACHEL", &comment)
	return &comment, nil
}

func (s *CommentsStore) Update(ctx context.Context, comment *Comment) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := `
	UPDATE comments
	SET content = $1
	WHERE id = $2
	RETURNING updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err = tx.QueryRowContext(ctx, query, comment.Content, comment.ID).Scan(&comment.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *CommentsStore) DeleteByCommentID(ctx context.Context, commentID int64) error {
	query := `
	DELETE FROM comments
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, commentID)
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

func (s *CommentsStore) DeleteByPostID(ctx context.Context, postID int64) error {
	query := `
	DELETE FROM comments
	WHERE post_id = $1
	`
	_, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}
	return nil
}
