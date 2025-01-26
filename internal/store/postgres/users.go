package postgres

import (
	"context"
	"database/sql"
	"os/user"

)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}
	type PostgresUserStore struct {
	db *sql.DB
}

func (s *PostgresUserStore) Create(ctx context.Context) error {
	query := `
INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at	
`
err := s.db.QueryRowContext(
	ctx,
	query,
	user.Username,
	user.Password,
	user.Email,
).Scan(
	&user.ID,
	&user.CreatedAt,
)

if err != nil {
	return err
} else {
	return nil
}
