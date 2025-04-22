package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Location struct {
	ID          int64   `json:"id"`
	Street      string  `json:"street"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	County      string  `json:"county"`
	ZIPCode     string  `json:"zip_code"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type LocationStore struct {
	db *sql.DB
}

func NewLocationStore(db *sql.DB) *LocationStore {
	return &LocationStore{db: db}
}

func (s *LocationStore) Create(ctx context.Context, tx *sql.Tx, location *Location) (Location, error) {
	fmt.Println("Create Location Started")    // Debugging line
	fmt.Println("Location values:", location) // Debugging line

	var query string
	if location.Street == "" {
		fmt.Println("Query is: Insert INTO locations", location.Street, location.City, location.State, location.County, location.ZIPCode, location) // Debugging line
		query = `
        INSERT INTO locations (street, city, state, county, zip_code, country, country_code, latitude, longitude, is_precise)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, false)
        RETURNING id
        `
	} else {
		query = `
        INSERT INTO locations (street, city, state, county, zip_code, country, country_code, latitude, longitude, is_precise)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, true)
        RETURNING id
        `
	}

	location.Normalize()                                 // 👍 perfect place to normalize
	fmt.Println("Normalized Location values:", location) // Debugging line

	var id int64
	err := tx.QueryRowContext(ctx, query,
		location.Street,
		location.City,
		location.State,
		location.County,
		location.ZIPCode,
		location.Country,
		location.CountryCode,
		location.Latitude,
		location.Longitude,
	).Scan(&id)
	if err != nil {
		fmt.Println("Error executing query:", err) // Debugging line
		return Location{}, fmt.Errorf("inserting location: %w", err)
	}

	fmt.Println("Location inserted with ID:", id) // Debugging line
	location.ID = id
	return *location, nil
}

func (s *LocationStore) Get(ctx context.Context, id int64) (Location, error) {
	query := `
		SELECT id, street, city, state, zip_code, country, latitude, longitude
		FROM locations
		WHERE id = $1
	`

	fmt.Println("Get Location Started") // Debugging line

	var location Location
	err := s.db.QueryRowContext(ctx, query, id).Scan(
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
		if err == sql.ErrNoRows {
			return Location{}, fmt.Errorf("location not found")
		}
		return Location{}, fmt.Errorf("getting location: %w", err)
	}
	return location, nil
}

// func (s *LocationStore) Update(ctx context.Context, tx *sql.Tx, location *Location) error {
// 	countQuery := `
// 	SELECT COUNT(*) FROM users WHERE location_id = $1;
// 	`
// 	var count int
// 	err := tx.QueryRowContext(ctx, countQuery, location.ID).Scan(&count)
// 	if err != nil {
// 		return fmt.Errorf("counting users with location_id %d: %w", location.ID, err)
// 	}
// 	if count > 1 {
// 		//location is shared by more than one user

// 		// Create a new location
// 		newLoc, err := s.Create(ctx, tx, location)
// 		if err != nil {
// 			return fmt.Errorf("creating new location: %w", err)
// 		}
// 		// You now need to update the user’s location_id in their row
// 		updateUserLocQuery := `UPDATE users SET location_id = $1 WHERE id = $2`
// 		_, err = tx.ExecContext(ctx, updateUserLocQuery, newLoc.ID, userID) // You'll need userID passed in
// 		if err != nil {
// 			return fmt.Errorf("updating user location_id: %w", err)
// 		}

// 		return nil

// 	} else {
// 		//location is not shared by another user, update it

// 		updateQuery := `
// 			UPDATE locations
// 			SET street = $1, city = $2, state = $3, zip_code = $4, country = $5, latitude = $6, longitude = $7
// 			WHERE id = $8
// 		`

// 		location.Normalize() // 👍 perfect place to normalize

// 		_, err := tx.ExecContext(ctx, updateQuery,
// 			location.Street,
// 			location.City,
// 			location.State,
// 			location.ZIPCode,
// 			location.Country,
// 			location.Latitude,
// 			location.Longitude,
// 			location.ID,
// 		)
// 		if err != nil {
// 			return fmt.Errorf("updating location: %w", err)
// 		}
// 		return nil
// 	}
// }

func (s *LocationStore) GetByLocationPrecise(ctx context.Context, location *Location) (Location, error) {
	fmt.Println("Get Location By Location Started")             // Debugging line
	fmt.Println("Querying for location with values:", location) // Debugging line
	query := `
		SELECT id, street, city, state, zip_code, country, latitude, longitude
		FROM locations
		WHERE is_precise=true AND street = $1 AND city = $2 AND state = $3 AND zip_code = $4 AND country = $5
	`
	var loc Location
	err := s.db.QueryRowContext(ctx, query,
		location.Street,
		location.City,
		location.State,
		location.ZIPCode,
		location.Country,
	).Scan(
		&loc.ID,
		&loc.Street,
		&loc.City,
		&loc.State,
		&loc.ZIPCode,
		&loc.Country,
		&loc.Latitude,
		&loc.Longitude,
	)
	fmt.Println("Query executed")     // Debugging line
	fmt.Println("Query result:", loc) // Debugging line
	if err != nil {
		if err == sql.ErrNoRows {
			return Location{}, fmt.Errorf("location not found")
		}
		return Location{}, fmt.Errorf("getting location: %w", err)
	}
	fmt.Println("Location found:", loc) // Debugging line
	return loc, nil
}

func (s *LocationStore) GetByLocation(ctx context.Context, location *Location) (Location, error) {
	fmt.Println("Get Location By Location Started")             // Debugging line
	fmt.Println("Querying for location with values:", location) // Debugging line
	query := `
		SELECT id, street, city, state, zip_code, country, latitude, longitude
		FROM locations
		WHERE is_precise = false AND street = $1 AND city = $2 AND state = $3 AND zip_code = $4 AND country = $5
	`
	var loc Location
	err := s.db.QueryRowContext(ctx, query,
		location.Street,
		location.City,
		location.State,
		location.ZIPCode,
		location.Country,
	).Scan(
		&loc.ID,
		&loc.Street,
		&loc.City,
		&loc.State,
		&loc.ZIPCode,
		&loc.Country,
		&loc.Latitude,
		&loc.Longitude,
	)
	fmt.Println("Query executed")     // Debugging line
	fmt.Println("Query result:", loc) // Debugging line
	if err != nil {
		if err == sql.ErrNoRows {
			return Location{}, fmt.Errorf("location not found")
		}
		return Location{}, fmt.Errorf("getting location: %w", err)
	}
	fmt.Println("Location found:", loc) // Debugging line
	return loc, nil
}

func (l *Location) Normalize() {
	l.Street = strings.ToUpper(l.Street)
	l.City = strings.ToUpper(l.City)
	l.State = strings.ToUpper(l.State)
	l.ZIPCode = strings.ToUpper(l.ZIPCode)
	l.Country = strings.ToUpper(l.Country)
}

func (s *LocationStore) GetGeneralLocationByZip(ctx context.Context, zipCode string) (Location, error) {
	// Start of the function, print the input
	fmt.Println("GetGeneralLocationByZip called with zipCode:", zipCode)

	query := `
		SELECT id, street, city, state, county, zip_code, country, country_code, latitude, longitude
		FROM locations
		WHERE zip_code = $1 AND is_precise = false
	`
	// Log the query for debugging
	fmt.Println("Executing query:", query)

	var loc Location
	// Execute the query and scan the result into the Location struct
	err := s.db.QueryRowContext(ctx, query, zipCode).Scan(
		&loc.ID,
		&loc.Street,
		&loc.City,
		&loc.State,
		&loc.County,
		&loc.ZIPCode,
		&loc.Country,
		&loc.CountryCode,
		&loc.Latitude,
		&loc.Longitude,
	)

	if err != nil {
		// Log the error in case of failure
		fmt.Println("Error during QueryRowContext execution:", err)

		// Specific handling for no rows found
		if err == sql.ErrNoRows {
			fmt.Println("\n\n--GetGeneralLocationByZip--\n\nNo location found for zipCode:", zipCode)
			return Location{}, sql.ErrNoRows
		}

		// Log the error for other types of database failures
		fmt.Println("Error executing query:", err)
		return Location{}, fmt.Errorf("getting location: %w", err)
	}

	// Log the result for debugging
	fmt.Println("Location found:", loc)

	return loc, nil
}

// func (s *LocationStore) GetGeneralLocationByZip(ctx context.Context, zipCode string) (Location, error) {

// 	query := `
// 		SELECT id, street, city, state, county, zip_code, country, country_code, latitude, longitude
// 		FROM locations
// 		WHERE zip_code = $1 AND is_precise = false
// 	`

// 	var loc Location
// 	err := s.db.QueryRowContext(ctx, query, zipCode).Scan(
// 		&loc.ID,
// 		&loc.Street,
// 		&loc.City,
// 		&loc.State,
// 		&loc.County,
// 		&loc.ZIPCode,
// 		&loc.Country,
// 		&loc.CountryCode,
// 		&loc.Latitude,
// 		&loc.Longitude,
// 	)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return Location{}, fmt.Errorf("location not found")
// 		}
// 		return Location{}, fmt.Errorf("getting location: %w", err)
// 	}
// 	return loc, nil
// }

func (s *LocationStore) DB() *sql.DB {
	return s.db
}

func (s *LocationStore) GetLocationsByBoundingBox(ctx context.Context, minLat, maxLat, minLon, maxLon float64) ([]Location, error) {
	query := `
		SELECT id, street, city, state, county, zip_code, country, country_code, latitude, longitude
		FROM locations
		WHERE latitude BETWEEN $1 AND $2 AND longitude BETWEEN $3 AND $4
	`

	rows, err := s.db.QueryContext(ctx, query, minLat, maxLat, minLon, maxLon)
	if err != nil {
		return nil, fmt.Errorf("getting locations: %w", err)
	}
	defer rows.Close()

	var locations []Location
	for rows.Next() {
		var loc Location
		if err := rows.Scan(
			&loc.ID,
			&loc.Street,
			&loc.City,
			&loc.State,
			&loc.County,
			&loc.ZIPCode,
			&loc.Country,
			&loc.CountryCode,
			&loc.Latitude,
			&loc.Longitude,
		); err != nil {
			return nil, fmt.Errorf("scanning location: %w", err)
		}
		locations = append(locations, loc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over rows: %w", err)
	}

	return locations, nil
}
