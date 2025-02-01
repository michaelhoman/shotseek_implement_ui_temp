package store

import (
	"context"
	"database/sql"
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
}
type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
INSERT INTO users ( email, password, first_name, last_name, zip_code, city, state) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at	
`
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
