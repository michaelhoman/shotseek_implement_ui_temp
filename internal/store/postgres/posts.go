package postgres

import (
	"context"
	"database/sql"
)

type PostgresPostStore struct {
	db *sql.DB
}

func (s *PostgresPostStore) Create(ctx context.Context) error {
	// Implementation for creating a post
	return nil
}
