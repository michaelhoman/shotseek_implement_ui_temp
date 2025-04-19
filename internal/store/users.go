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

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/michaelhoman/ShotSeek/internal/utils"
	"golang.org/x/crypto/bcrypt"
	// "os/user" // Remove this import as it is not needed
)

var (
	ErrDuplicateEmail = errors.New("a user with that email already exists")
	ErrDuplicateUser  = errors.New("a user with that username already exists")
)

type User struct {
	ID         uuid.UUID `db:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Password   password  `json:"-"`
	LocationID int64     `json:"location_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Location   *Location `json:"location,omitempty"`
	IsActive   bool      `json:"is_active"`
	Version    int       `json:"version"`
}

type password struct {
	// text *string
	hash []byte
}

// type Location struct {
// 	ID        int64   `json:"id"`
// 	Street    string  `json:"street"`
// 	City      string  `json:"city"`
// 	State     string  `json:"state"`
// 	ZIPCode   string  `json:"zip_code"`
// 	Country   string  `json:"country"`
// 	Latitude  float64 `json:"latitude"`
// 	Longitude float64 `json:"longitude"`
// }

func (p *password) Set(plain string) error {
	fmt.Println("Setting password") // TODO: Remove Debugging line

	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	fmt.Println("Set Hash is:", hash) // TODO: Remove Debugging line
	if err != nil {
		return err
	}

	p.hash = hash
	return nil
}

func (p *password) Compare(storedPassword string) error {
	fmt.Println("Comparing password")             //TODO: Remove this line
	fmt.Println("p.hash", p.hash)                 //TODO: Remove this line
	fmt.Println("[]byte(p.hash)", storedPassword) //TODO: Remove this line
	return bcrypt.CompareHashAndPassword(p.hash, []byte(storedPassword))
}

type UserStore struct {
	db            *sql.DB
	locationStore *LocationStore
}

// Satisfy the Users interface
func (s *UserStore) LocationStore() *LocationStore {
	return s.locationStore
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User, location *Location) error {
	var locationID int64 // Pointer so we can insert NULL if needed

	fmt.Println()
	fmt.Println()
	fmt.Println("Location is:", location)

	if location != nil {
		if !location.IsValid() {
			return fmt.Errorf("location provided but missing required fields (city, state, or zip code)")
		}
		// continue...
		fmt.Println("Location is valid, proceeding with location check")
		submittedLocation, err := s.locationStore.GetByLocation(ctx, location)

		if err != nil {
			utils.Logger.Info("Error querying location: %v", err)
		}
		if submittedLocation.ID != 0 {
			fmt.Println("Location already exists, using existing location ID")
			locationID = submittedLocation.ID
		} else {
			// Location does not exist, first query Nominatim API and then insert new location into locations table

			fmt.Println("Location does not exist, inserting new location")
			// locationID = submittedLocation.ID

			// Debugging logs
			log.Printf("Attempting to find location with email: %s", user.Email)
			log.Printf("Location query with: Street=%s, City=%s, State=%s, ZipCode=%s", location.Street, location.City, location.State, location.ZIPCode)

			// fmt.Printf(err)
			//print err
			fmt.Println("Error 12 is:", err)

			newLocation, err := s.locationStore.Create(ctx, tx, location)
			if err != nil {
				log.Printf("Error inserting location: %v", err)
				return fmt.Errorf("inserting location: %w", err)
			}
			locationID = newLocation.ID
			log.Printf("Inserted new location with ID: %d", locationID)
		}

		// Debugging log before inserting user
		log.Printf("Inserting user: %s", user.Email)
		log.Printf("User info - First Name: %s, Last Name: %s", user.FirstName, user.LastName)
		if locationID != 0 {
			log.Printf("Location ID to insert: %d", locationID)
		} else {
			log.Printf("No location ID provided, setting to NULL")
		}

		newUserID := uuid.New()

		// Debugging log for new user ID
		log.Printf("New user ID generated: %s", newUserID)
		log.Printf("Inserting user: id=%s, first_name=%s, last_name=%s, email=%s, password=%s, location_id=%d", newUserID, user.FirstName, user.LastName, user.Email, user.Password.hash, locationID)

		// Insert the user with correct location_id handling
		// userInsertQuery := `
		// 	INSERT INTO users (id, email, password, first_name, last_name, location_id)
		// 	VALUES ( $1, $2, $3, $4, $5, $6)
		// 	RETURNING id, created_at
		// 	`
		userInsertQuery := `
		INSERT INTO users (id, email, password, first_name, last_name, location_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

		ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
		defer cancel()

		// Check if locationID is valid or nil before inserting
		err = tx.QueryRowContext(
			ctx,
			userInsertQuery,
			newUserID,
			user.Email,
			user.Password.hash, // assuming a method, or use []byte directly if applicable
			user.FirstName,
			user.LastName,
			locationID, // Use NULL if no location exists, or the valid locationID
		).Scan(&user.ID, &user.CreatedAt)

		if err != nil {
			var pgErr *pq.Error
			if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.Constraint == "users_email_key" {
				log.Printf("Duplicate email found: %s", user.Email)
				return ErrDuplicateEmail
			}
			log.Printf("Error inserting user: %v", err)
			return fmt.Errorf("inserting user: %w", err)
		}

		// Success debugging
		log.Printf("User successfully created with ID: %s at %v", user.ID, user.CreatedAt)

		return nil
	}
	return nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	fmt.Println("GetByEmail called with email:", email) // Debugging line
	// Check if email is empty
	query := `
SELECT
    u.id, u.email, u.first_name, u.last_name, u.created_at, u.updated_at, u.version, u.is_active,
    l.id AS location_id, l.street, l.city, l.state, l.zip_code, l.country, l.latitude, l.longitude
FROM users u
LEFT JOIN locations l ON u.location_id = l.id
WHERE u.email = $1
`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := User{}
	location := Location{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
		&user.IsActive,
		// Scan location fields
		&location.ID,
		&location.Street,
		&location.City,
		&location.State,
		&location.ZIPCode,
		&location.Country,
		&location.Latitude,
		&location.Longitude,
	)

	if err != nil {
		fmt.Println("Error occurred while querying user:", err) // Debugging line
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	fmt.Println("User found:", user) // Debugging line
	user.LocationID = location.ID
	return &user, nil
}

func (s *UserStore) GetByEmailWithPassword(ctx context.Context, email string) (*User, error) {
	fmt.Println("GetByEmail called with email:", email) // Debugging line
	// Check if email is empty
	query := `
SELECT
    u.id, u.email, u.password, u.first_name, u.last_name, u.created_at, u.updated_at, u.version, u.is_active,
    l.id AS location_id, l.street, l.city, l.state, l.zip_code, l.country, l.latitude, l.longitude
FROM users u
LEFT JOIN locations l ON u.location_id = l.id
WHERE u.email = $1
`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := User{}
	location := Location{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password.hash,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
		&user.IsActive,
		// Scan location fields
		&location.ID,
		&location.Street,
		&location.City,
		&location.State,
		&location.ZIPCode,
		&location.Country,
		&location.Latitude,
		&location.Longitude,
	)

	if err != nil {
		fmt.Println("Error occurred while querying user:", err) // Debugging line
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	fmt.Println("User found:", user) // Debugging line
	user.LocationID = location.ID
	return &user, nil
}

// func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
// 	query := `
// SELECT id, email, password, first_name, last_name, created_at, updated_at, version, is_active
// FROM users
// WHERE email = $1
// `
// 	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
// 	defer cancel()

// 	user := User{}

// 	err := s.db.QueryRowContext(ctx, query, email).Scan(
// 		&user.ID,
// 		&user.Email,
// 		&user.Password.hash,
// 		&user.FirstName,
// 		&user.LastName,
// 		&user.CreatedAt,
// 		&user.UpdatedAt,
// 		&user.Version,
// 		&user.IsActive,
// 	)

// 	if err != nil {
// 		switch {
// 		case errors.Is(err, sql.ErrNoRows):
// 			return nil, ErrNotFound
// 		default:
// 			return nil, err
// 		}
// 	}
// 	return &user, nil
// }

func (s *UserStore) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `
SELECT id, email, first_name, last_name, created_at, updated_at, version
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

func (s *UserStore) Update(ctx context.Context, user *User, location *Location) error {
	// Begin a transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // Make sure to rollback in case of an error

	// If location is provided (indicating it might be updated), update the location table
	if location != nil {
		// Check if the location exists or needs to be created/updated
		updateLocationQuery := `
		UPDATE locations
		SET street = $1, city = $2, state = $3, zip_code = $4, country = $5, latitude = $6, longitude = $7
		WHERE id = $8
		RETURNING id
		`
		err := tx.QueryRowContext(
			ctx,
			updateLocationQuery,
			location.Street,
			location.City,
			location.State,
			location.ZIPCode,
			location.Country,
			location.Latitude,
			location.Longitude,
			location.ID, // Assuming the location ID is already known
		).Scan(&location.ID)

		if err != nil {
			// If location doesn't exist, insert a new one
			if errors.Is(err, sql.ErrNoRows) {
				insertLocationQuery := `
				INSERT INTO locations (street, city, state, zip_code, country, latitude, longitude)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING id
				`
				err := tx.QueryRowContext(
					ctx,
					insertLocationQuery,
					location.Street,
					location.City,
					location.State,
					location.ZIPCode,
					location.Country,
					location.Latitude,
					location.Longitude,
				).Scan(&location.ID)

				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	// Now update the user table with the new or unchanged location ID
	updateUserQuery := `
	UPDATE users
	SET email = $1, password = $2, first_name = $3, last_name = $4, location_id = $5, version = version + 1, updated_at = NOW()
	WHERE id = $6 AND version = $7
	RETURNING version
	`

	err = tx.QueryRowContext(
		ctx,
		updateUserQuery,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		location.ID, // Update location_id in the users table
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

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *UserStore) Delete(ctx context.Context, id uuid.UUID) error {
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

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, location *Location, token string, invitationExp time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Create(ctx, tx, user, location); err != nil {
			return err
		}
		err := s.createUserInvitation(ctx, tx, user.ID, invitationExp, token)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, userID uuid.UUID, invitationExp time.Duration, token string) error {
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
	SELECT u.id, u.email, u.first_name, u.last_name, u.location_id, u.created_at, u.is_active
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
		&user.LocationID,
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
	fmt.Println("GetUserFromInvitation User found:", user) // Debugging line
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

func (s *UserStore) GetHashedPassword(ctx context.Context, email string) (string, error) {
	fmt.Println("GetHashedPassword called with email:", email) // Debugging line

	query := `
SELECT password
FROM users
WHERE email = $1
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	var passwordHash string
	err := s.db.QueryRowContext(ctx, query, email).Scan(&passwordHash)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			fmt.Println("Error occurred while querying password:", err) // Debugging line
			return "", ErrNotFound
		default:
			return "", err
		}
	}
	fmt.Println("GetHashedPassword: Hashed Password:", passwordHash) // Debugging line
	return passwordHash, nil
}

func (s *UserStore) deleteInvitation(ctx context.Context, tx *sql.Tx, userID uuid.UUID) error {
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

func (l *Location) IsValid() bool {
	return l != nil && l.City != "" && l.State != "" && l.ZIPCode != ""
}
