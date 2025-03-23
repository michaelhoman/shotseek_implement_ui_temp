package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type RefreshToken struct {
	UserEmail string    `json:"user_email"`
	TokenHash string    `json:"token_hash"`
	StoredFP  string    `json:"stored_fp"`
	ExpiresAt time.Time `json:"expires_at"`
}

type TokenStore struct {
	db *sql.DB
}

func (s *TokenStore) UpdateRefreshToken(ctx context.Context, userEmail, token_hash string, expiresAt time.Time) error {
	query := `
    INSERT INTO refresh_tokens (user_email, token_hash, expires_at)
    VALUES ($1, $2, $3)
    ON CONFLICT(token) 
    DO UPDATE SET token = $2, expires_at = $3
    `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userEmail, token_hash, expiresAt)
	if err != nil {
		return err
	}
	return nil
}

// GetRefreshTokens retrieves all refresh tokens for a user
func (s *TokenStore) GetRefreshTokens(ctx context.Context, userEmail string) ([]RefreshToken, error) {
	query := `
	SELECT user_email, token_hash, stored_fp, expires_at
	FROM refresh_tokens
	WHERE user_email = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	fmt.Println("1* GetRefreshTokens query: ", query)

	rows, err := s.db.QueryContext(
		ctx,
		query,
		userEmail,
	)
	if err != nil {
		return nil, err
	}

	fmt.Println("2* GetRefreshTokens rows: ", rows)

	fmt.Println()

	defer rows.Close()

	var tokens []RefreshToken

	if err := rows.Err(); err != nil {
		fmt.Println("No refesh tokens found")
		return nil, err
	}

	for rows.Next() {
		var token RefreshToken
		if err := rows.Scan(
			&token.UserEmail,
			&token.TokenHash,
			&token.StoredFP,
			&token.ExpiresAt); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}
