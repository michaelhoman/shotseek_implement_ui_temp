package postgres

import (
	"context"
	"database/sql"
)

type PostgresUserStore struct {
	db *sql.DB
}

func (s *PostgresUserStore) Create(ctx context.Context) error {
	// Implementation for creating a user
	return nil
}
