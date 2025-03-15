package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	// "os/user" // Remove this import as it is not needed
)

var (
	ErrDuplicateEmail = errors.New("a user with that email already exists")
	ErrDuplicateUser  = errors.New("a user with that username already exists")
)

type User struct {
	ID        int64    `json:"id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	Zipcode   string   `json:"zip_code"`
	City      string   `json:"city"`
	State     string   `json:"state"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	IsActive  bool     `json:"is_active"`
	Version   int      `json:"version"`
}

type password struct {
	// text *string
	hash []byte
}

func (p *password) Set(plain string) error {
	fmt.Println("Setting password") // TODO: Remove Debugging line
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.hash = hash
	return nil
}

func (p *password) Compare(plain string) error {
	fmt.Println("Comparing password")            //TODO: Remove this line
	fmt.Println("p.hash", p.hash)                //TODO: Remove this line
	fmt.Println("[]byte(p.hash)", []byte(plain)) //TODO: Remove this line
	return bcrypt.CompareHashAndPassword(p.hash, []byte(plain))
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
INSERT INTO users ( email, password, first_name, last_name, zip_code, city, state) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at	
`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.Password.hash,
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
		switch {
		case err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"":
			return ErrDuplicateEmail
		default:
			return err
		}
	} else {
		return nil

	}
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
SELECT id, email, password, first_name, last_name, zip_code, city, state, created_at, updated_at, version, is_active
FROM users
WHERE email = $1
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := User{}

	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password.hash,
		&user.FirstName,
		&user.LastName,
		&user.Zipcode,
		&user.City,
		&user.State,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
		&user.IsActive,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
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

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}
		err := s.createUserInvitation(ctx, tx, user.ID, invitationExp, token)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, userID int64, invitationExp time.Duration, token string) error {
	query := `
INSERT INTO user_invitations (user_id, token, expires_at) VALUES ($1, $2,  $3)
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID, token, time.Now().Add(invitationExp))
	if err != nil {
		return err
	}
	return err
}

func (s *UserStore) Activate(ctx context.Context, token string) error {
	// find the user that this token corresponds to
	// check if the token is expired
	// if expired return an error
	// if not expired
	// activate the user
	// delete the invitation
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}
		// Update user
		user.IsActive = true
		if err := s.update(ctx, tx, user); err != nil {
			return err
		}
		// Clean Invitations

		if err := s.deleteInvitation(ctx, tx, user.ID); err != nil {
			return err
		}
		return nil
	})

}

func (s *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
	SELECT u.id, u.email, u.first_name, u.last_name, u.zip_code, u.city, u.state, u.created_at, u.is_active
	FROM users u
	JOIN user_invitations ui ON u.id = ui.user_id
	WHERE ui.token = $1 AND ui.expires_at > $2
	`
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Zipcode,
		&user.City,
		&user.State,
		&user.CreatedAt,
		&user.IsActive,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (s *UserStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
	UPDATE users
	SET is_active = $1
	WHERE id = $2
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.IsActive, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) deleteInvitation(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `
	DELETE FROM user_invitations
	WHERE user_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}
