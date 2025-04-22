package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
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
		GetByEmailWithPassword(context.Context, string) (*User, error)
		GetByID(context.Context, uuid.UUID) (*User, error)
		Create(context.Context, *sql.Tx, *User, *Location) error
		update(context.Context, *sql.Tx, *User) error
		Update(context.Context, *User, *Location) error
		Delete(context.Context, uuid.UUID) error
		CreateAndInvite(context.Context, *User, *Location, string, time.Duration) error
		GetHashedPassword(context.Context, string) (string, error)
		LocationStore() *LocationStore
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
		UpdateRefreshToken(ctx context.Context, userID uuid.UUID, token string, stored_fp string, expiresAt time.Time) error
		GetRefreshTokens(ctx context.Context, userID uuid.UUID) ([]*RefreshToken, error)
		GetByRefreshTokenHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	}
	Locations interface {
		Create(context.Context, *sql.Tx, *Location) (Location, error)
		Get(context.Context, int64) (Location, error)
		GetByLocation(context.Context, *Location) (Location, error)
		GetGeneralLocationByZip(ctx context.Context, zipCode string) (Location, error)
		GetLocationsByBoundingBox(ctx context.Context, minLat, maxLat, minLon, maxLon float64) ([]Location, error)
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db, NewLocationStore(db)},
		Comments:  &CommentsStore{db},
		Tokens:    &TokenStore{db},
		Locations: &LocationStore{db},
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
