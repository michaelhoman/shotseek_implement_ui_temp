package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	UserID    string    `json:"user_id"`
	TokenHash string    `json:"token_hash"`
	StoredFP  string    `json:"stored_fp"`
	ExpiresAt time.Time `json:"expires_at"`
}

type TokenStore struct {
	db *sql.DB
}

// func (s *TokenStore) UpdateRefreshToken(ctx context.Context, userEmail, token_hash string, expiresAt time.Time) error {
// 	query := `
//     INSERT INTO refresh_tokens (user_email, token_hash, expires_at)
//     VALUES ($1, $2, $3)
//     ON CONFLICT(token_hash)
//     DO UPDATE SET token_hash = $2, expires_at = $3
//     `
// 	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
// 	defer cancel()

// 	_, err := s.db.ExecContext(ctx, query, userEmail, token_hash, expiresAt)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (s *TokenStore) UpdateRefreshToken(ctx context.Context, userEmail, token_hash, stored_fp string, expiresAt time.Time) error {
// 	query := `
//     INSERT INTO refresh_tokens (user_email, token_hash, stored_fp, expires_at)
//     VALUES ($1, $2, $3, $4)
//     ON CONFLICT(token_hash)
//     DO UPDATE SET token_hash = $2, stored_fp = $3, expires_at = $4
//     `
// 	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
// 	defer cancel()

// 	fmt.Println("Executing query:", query) // Log the query
// 	fmt.Printf("Inserting token for user: %s, token_hash: %s\n", userEmail, token_hash)

// 	_, err := s.db.ExecContext(ctx, query, userEmail, token_hash, stored_fp, expiresAt)
// 	if err != nil {
// 		fmt.Println("Error inserting token:", err) // Log any errors
// 		return err
// 	}
// 	return nil
// }

func (s *TokenStore) UpdateRefreshToken(ctx context.Context, user_id uuid.UUID, token_hash, stored_fp string, expiresAt time.Time) error {
	fmt.Println("******") // Debugging
	fmt.Println("******") // Debugging
	fmt.Println("UpdateRefreshToken called with user_id:", user_id, "token_hash:", token_hash, "stored_fp:", stored_fp, "expiresAt:", expiresAt)
	query := `
    INSERT INTO refresh_tokens (user_id, token_hash, stored_fp, expires_at)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT(user_id, stored_fp) 
    DO UPDATE SET token_hash = $2, stored_fp = $3, expires_at = $4
    `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	fmt.Println("Executing query:", query) // Log the query
	fmt.Printf("Inserting token for user: %s, token_hash: %s\n", user_id, token_hash)

	_, err := s.db.ExecContext(ctx, query, user_id, token_hash, stored_fp, expiresAt)
	if err != nil {
		fmt.Println("Error inserting token:", err) // Log any errors
		return err
	}
	return nil
}

// GetRefreshTokens retrieves all refresh tokens for a user
func (s *TokenStore) GetRefreshTokens(ctx context.Context, userID uuid.UUID) ([]*RefreshToken, error) {
	query := `
	SELECT user_id, token_hash, stored_fp, expires_at
	FROM refresh_tokens
	WHERE user_id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	fmt.Println("1* GetRefreshTokens query: ", query)

	rows, err := s.db.QueryContext(
		ctx,
		query,
		userID,
	)
	if err != nil {
		return nil, err
	}

	fmt.Println("2* GetRefreshTokens rows: ", rows)

	fmt.Println()

	defer rows.Close()

	var refreshTokens []*RefreshToken

	if err := rows.Err(); err != nil {
		fmt.Println("No refesh tokens found")
		return nil, err
	}

	for rows.Next() {
		var token RefreshToken
		if err := rows.Scan(
			&token.UserID,
			&token.TokenHash,
			&token.StoredFP,
			&token.ExpiresAt); err != nil {
			return nil, err
		}
		refreshTokens = append(refreshTokens, &token)
	}

	return refreshTokens, nil
}

// GetByTokenHash retrieves a refresh token by its hash

func (s *TokenStore) GetByRefreshTokenHash(ctx context.Context, tokenHash string) (*RefreshToken, error) {
	query := `
    SELECT user_id, token_hash, stored_fp, expires_at
    FROM refresh_tokens
    WHERE token_hash = $1
    `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration) // Ensure timeout duration is defined
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, tokenHash)

	var token RefreshToken
	// Scan the result
	if err := row.Scan(
		&token.UserID,
		&token.TokenHash,
		&token.StoredFP,
		&token.ExpiresAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("refresh token not found") // Graceful error handling
		}
		return nil, err
	}

	return &token, nil
}
