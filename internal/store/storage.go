package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = 5 * time.Second
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		Update(context.Context, *Post) error
		Delete(context.Context, int64) error
	}
	Users interface {
		// create(context.Context, *sql.Tx, *User) error
		Activate(context.Context, string) error
		GetByEmail(context.Context, string) (*User, error)
		GetByID(context.Context, int64) (*User, error)
		Create(context.Context, *sql.Tx, *User) error
		Update(context.Context, *User) error
		Delete(context.Context, int64) error
		CreateAndInvite(context.Context, *User, string, time.Duration) error
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostID(context.Context, int64) ([]Comment, error)
		GetByCommentID(context.Context, int64) (*Comment, error)
		Update(context.Context, *Comment) error
		DeleteByCommentID(context.Context, int64) error
		DeleteByPostID(context.Context, int64) error
	}
	Tokens interface {
		UpdateRefreshToken(ctx context.Context, userEmail, token string, stored_fp string, expiresAt time.Time) error
		GetRefreshTokens(ctx context.Context, userEmail string) ([]*RefreshToken, error)
		GetByRefreshTokenHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentsStore{db},
		Tokens:   &TokenStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()

}
