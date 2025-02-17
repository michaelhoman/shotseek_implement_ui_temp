package store

import (
	"context"
	"database/sql"
	"errors"
	"log"
	// "os/user" // Remove this import as it is not needed
)

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Zipcode   string `json:"zip_code"`
	City      string `json:"city"`
	State     string `json:"state"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Version   int    `json:"version"`
}
type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
INSERT INTO users ( email, password, first_name, last_name, zip_code, city, state) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at	
`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Zipcode,
		user.City,
		user.State,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		return err
	} else {
		return nil

	}
}

func (s *UserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
SELECT id, email, first_name, last_name, zip_code, city, state, created_at, updated_at, version
FROM users
WHERE id = $1
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := User{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Zipcode,
		&user.City,
		&user.State,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
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
	return &user, nil
}

func (s *UserStore) Update(ctx context.Context, user *User) error {
	query := `
UPDATE users
SET email = $1, password = $2, first_name = $3, last_name = $4, zip_code = $5, city = $6, state = $7, version = version + 1, updated_at = NOW()
WHERE id = $8 AND version = $9
RETURNING version
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Zipcode,
		user.City,
		user.State,
		user.ID,
		user.Version,
	).Scan(&user.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}
	return nil
}

func (s *UserStore) Delete(ctx context.Context, id int64) error {
	query := `
DELETE FROM users
WHERE id = $1
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, id)
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
