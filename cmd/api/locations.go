package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	location_package "github.com/michaelhoman/ShotSeek/internal/location"
	"github.com/michaelhoman/ShotSeek/internal/store"
	"github.com/michaelhoman/ShotSeek/internal/utils"
)

type Address struct {
	Postcode    string `json:"postcode"`
	City        string `json:"city"`
	Town        string `json:"town"`
	Village     string `json:"village"`
	County      string `json:"county"`
	State       string `json:"state"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
}

// func (a *application) lookupByZip(ctx context.Context, tx *sql.Tx, zip string) (*store.Location, error) {
// 	fmt.Println("lookupByZip: Starting ZIP lookup for:", zip)
// 	if zip == "" {
// 		fmt.Println("lookupByZip: Empty ZIP code provided")
// 		return nil, fmt.Errorf("empty ZIP code provided")
// 	}

// 	// First, check if the location already exists in the database
// 	location, err := a.store.Locations.GetGeneralLocationByZip(ctx, zip)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			fmt.Println("\n\nErrNotFound - lookupByZip: No record found for ZIP:", zip)

// 			// Construct the URL for the Nominatim API lookup (by ZIP code)
// 			url := fmt.Sprintf("https://nominatim.openstreetmap.org/search?postalcode=%s&country=USA&format=json", zip)
// 			fmt.Println("lookupByZip: Constructed URL:", url)

// 			// Create the HTTP request
// 			req, err := http.NewRequest("GET", url, nil)
// 			if err != nil {
// 				fmt.Println("lookupByZip: Error creating request:", err)
// 				return nil, err
// 			}

// 			// Set a custom User-Agent header
// 			userAgent := "ShotSeek/1.0 (contact@shotseek.app)"
// 			req.Header.Set("User-Agent", userAgent)

// 			// Send the HTTP request
// 			client := &http.Client{}
// 			resp, err := client.Do(req)
// 			if err != nil {
// 				fmt.Println("lookupByZip: HTTP request error:", err)
// 				return nil, err
// 			}
// 			defer resp.Body.Close()

// 			// Handle non-OK status codes
// 			if resp.StatusCode != http.StatusOK {
// 				fmt.Println("lookupByZip: Non-OK HTTP status:", resp.Status)
// 				return nil, fmt.Errorf("nominatim error: %s", resp.Status)
// 			}

// 			// Decode the JSON response into a slice of results
// 			var results []location_package.NominatimResult
// 			if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
// 				fmt.Println("lookupByZip: JSON decode error:", err)
// 				return nil, err
// 			}

// 			// Handle the case when no results are returned
// 			if len(results) == 0 {
// 				fmt.Println("lookupByZip: No results found for ZIP:", zip)
// 				return nil, fmt.Errorf("no results found for ZIP %s", zip)
// 			}

// 			// Extract latitude and longitude from the first result
// 			lat, lon := results[0].Lat, results[0].Lon
// 			fmt.Println("lookupByZip: Found coordinates - Latitude:", lat, "Longitude:", lon)

// 			// Perform reverse geocoding using the latitude and longitude
// 			reverseUrl := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%s&lon=%s&format=json", lat, lon)
// 			fmt.Println("lookupByZip: Constructed reverse geocoding URL:", reverseUrl)

// 			// Create the HTTP request for reverse lookup
// 			req, err = http.NewRequest("GET", reverseUrl, nil)
// 			if err != nil {
// 				fmt.Println("lookupByZip: Error creating reverse request:", err)
// 				return nil, err
// 			}
// 			req.Header.Set("User-Agent", userAgent)

// 			// Send the reverse geocoding request
// 			resp, err = client.Do(req)
// 			if err != nil {
// 				fmt.Println("lookupByZip: Reverse geocoding HTTP request error:", err)
// 				return nil, err
// 			}
// 			defer resp.Body.Close()

// 			// Handle non-OK status codes for reverse geocoding
// 			if resp.StatusCode != http.StatusOK {
// 				fmt.Println("lookupByZip: Non-OK reverse geocoding HTTP status:", resp.Status)
// 				return nil, fmt.Errorf("reverse geocoding error: %s", resp.Status)
// 			}

// 			// Decode the reverse geocoding JSON response
// 			var reverseResult location_package.NominatimResult
// 			if err := json.NewDecoder(resp.Body).Decode(&reverseResult); err != nil {
// 				fmt.Println("lookupByZip: Reverse JSON decode error:", err)
// 				return nil, err
// 			}

// 			// Convert the Nominatim result to Location
// 			loc, err := location_package.LocationFromNominatim(&reverseResult)
// 			if err != nil {
// 				fmt.Println("lookupByZip: Error converting Nominatim result to Location struct:", err)
// 				return nil, err
// 			}

// 			// Store the location in the database
// 			fmt.Println("lookupByZip: Storing location in database")
// 			a.store.Locations.Create(ctx, tx, loc)

// 			// Return the location
// 			fmt.Println("lookupByZip: Returning location:", loc)
// 			return loc, nil
// 		}

// 		// If there's another error fetching from the database
// 		log.Println("lookupByZip: Error getting location from DB:", err)
// 		return nil, err
// 	}

// 	// Return the location if it's found in the database
// 	return &location, nil
// }

func (a *application) lookupByZip(ctx context.Context, tx *sql.Tx, zip string) (*store.Location, error) {
	fmt.Println("lookupByZip: Starting ZIP lookup for:", zip)
	if zip == "" {
		fmt.Println("lookupByZip: Empty ZIP code provided")
		return nil, fmt.Errorf("empty ZIP code provided")
	}

	// First, check if the location already exists in the database
	location, err := a.store.Locations.GetGeneralLocationByZip(ctx, zip)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("\n\nErrNotFound - lookupByZip: No record found for ZIP:", zip)

			// Construct the URL for the Nominatim API lookup (by ZIP code)
			url := fmt.Sprintf("https://nominatim.openstreetmap.org/search?postalcode=%s&country=USA&format=json", zip)
			fmt.Println("lookupByZip: Constructed URL:", url)

			// Create the HTTP request
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println("lookupByZip: Error creating request:", err)
				return nil, err
			}

			// Set a custom User-Agent header
			userAgent := "ShotSeek/1.0 (mthomanmt@gmail.com)"
			req.Header.Set("User-Agent", userAgent)

			// Send the HTTP request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("lookupByZip: HTTP request error:", err)
				return nil, err
			}
			defer resp.Body.Close()

			// Handle non-OK status codes
			if resp.StatusCode != http.StatusOK {
				fmt.Println("lookupByZip: Non-OK HTTP status:", resp.Status)
				return nil, fmt.Errorf("nominatim error: %s", resp.Status)
			}

			// Decode the JSON response into a slice of results
			var results []location_package.NominatimResult
			if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
				fmt.Println("lookupByZip: JSON decode error:", err)
				return nil, err
			}

			// Handle the case when no results are returned
			if len(results) == 0 {
				fmt.Println("lookupByZip: No results found for ZIP:", zip)
				return nil, fmt.Errorf("no results found for ZIP %s", zip)
			}

			// Extract latitude and longitude from the first result
			lat, lon := results[0].Lat, results[0].Lon
			fmt.Println("lookupByZip: Found coordinates - Latitude:", lat, "Longitude:", lon)

			// Perform reverse geocoding using the latitude and longitude
			reverseUrl := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%s&lon=%s&format=json", lat, lon)
			fmt.Println("lookupByZip: Constructed reverse geocoding URL:", reverseUrl)

			// Create the HTTP request for reverse lookup
			req, err = http.NewRequest("GET", reverseUrl, nil)
			if err != nil {
				fmt.Println("lookupByZip: Error creating reverse request:", err)
				return nil, err
			}
			req.Header.Set("User-Agent", userAgent)

			// Send the reverse geocoding request
			resp, err = client.Do(req)
			if err != nil {
				fmt.Println("lookupByZip: Reverse geocoding HTTP request error:", err)
				return nil, err
			}
			defer resp.Body.Close()

			// Handle non-OK status codes for reverse geocoding
			if resp.StatusCode != http.StatusOK {
				fmt.Println("lookupByZip: Non-OK reverse geocoding HTTP status:", resp.Status)
				return nil, fmt.Errorf("reverse geocoding error: %s", resp.Status)
			}

			// Decode the reverse geocoding JSON response
			var reverseResult location_package.NominatimResult
			if err := json.NewDecoder(resp.Body).Decode(&reverseResult); err != nil {
				fmt.Println("lookupByZip: Reverse JSON decode error:", err)
				return nil, err
			}

			// Convert the Nominatim result to Location
			loc, err := location_package.LocationFromNominatim(&reverseResult)
			if err != nil {
				fmt.Println("lookupByZip: Error converting Nominatim result to Location struct:", err)
				return nil, err
			}

			// Store the location in the database
			fmt.Println("lookupByZip: Storing location in database")
			a.store.Locations.Create(ctx, tx, loc)

			// Return the location
			fmt.Println("lookupByZip: Returning location:", loc)
			return loc, nil
		}

		// If there's another error fetching from the database
		log.Println("lookupByZip: Error getting location from DB:", err)
		return nil, err
	}

	// Return the location if it's found in the database
	return &location, nil
}

// LookupByZip godoc
//
//	@Summary		Lookup location by ZIP code
//	@Description	Lookup location by ZIP code
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param			ZIPCode	path		string	true	"ZIP code"
//	@Success		200		{object}	store.Location
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/locations/zip/{ZIPCode} [get]
//
//	@Security		ApiKeyAuth
//	@Security		BearerAuth
func (app *application) zipLookupHandler(w http.ResponseWriter, r *http.Request) {
	zip := chi.URLParam(r, "ZIPCode")
	fmt.Println("\n\ntestZipLookupHandler: ZIP code:", zip)
	if zip == "" {
		http.Error(w, "ZIP code is required", http.StatusBadRequest)
		return
	}
	locationStore, ok := app.store.Locations.(*store.LocationStore)
	if !ok {
		http.Error(w, "LocationStore is not initialized correctly", http.StatusInternalServerError)
		return
	}

	// Now you can access locationStore.db to begin a transaction
	tx, err := locationStore.DB().BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to begin transaction: %v", err), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // Ensure the transaction is rolled back in case of an error.

	// Look up the location by ZIP code
	location, err := app.lookupByZip(r.Context(), tx, zip)
	if err != nil {
		fmt.Println("zipLookupHandler: Error looking up ZIP code:", err)
		http.Error(w, fmt.Sprintf("Error looking up ZIP code: %v", err), http.StatusInternalServerError)
		return
	}

	if location == nil {
		http.Error(w, "No location found", http.StatusNotFound)
		return
	}

	// If everything is successful, commit the transaction
	if err := tx.Commit(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to commit transaction: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the found location as JSON
	utils.JsonResponse(w, http.StatusOK, location)
}

func getLocationFromCtx(r *http.Request) *store.Location {
	location, _ := r.Context().Value(userCtx).(*store.Location)
	return location
}

// func fromNominatimResult(n *NominatimResult) (*store.Location, error) {
// 	lat, err := strconv.ParseFloat(n.Lat, 64)
// 	if err != nil {
// 		return nil, err
// 	}

// 	lon, err := strconv.ParseFloat(n.Lon, 64)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &store.Location{
// 		ZIPCode:     n.Address.Postcode,
// 		City:        n.Address.City,
// 		County:      n.Address.County,
// 		State:       n.Address.State,
// 		Country:     n.Address.Country,
// 		CountryCode: n.Address.CountryCode,
// 		Latitude:    lat,
// 		Longitude:   lon,
// 	}, nil
// }

func (app *application) getNearbyByZip(w http.ResponseWriter, r *http.Request, zip string, miles float64) {
	fmt.Println("getNearbyByZip: ZIP code:", zip)
	if zip == "" {
		http.Error(w, "ZIP code is required", http.StatusBadRequest)
		return
	}

	locationStore, ok := app.store.Locations.(*store.LocationStore)
	if !ok {
		http.Error(w, "LocationStore is not initialized correctly", http.StatusInternalServerError)
		return
	}

	// Start a transaction
	tx, err := locationStore.DB().BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, "failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer func() {
		// Rollback if not already committed
		_ = tx.Rollback()
	}()

	// Call your function with the transaction
	providedLocation, err := app.lookupByZip(r.Context(), tx, zip)
	if err != nil {
		http.Error(w, "failed to look up zip: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// If everything is good, commit the transaction
	if err := tx.Commit(); err != nil {
		http.Error(w, "failed to commit transaction", http.StatusInternalServerError)
		return
	}

	if err != nil {
		fmt.Println("getNearbyByZip: Error looking up ZIP code:", err)
		http.Error(w, fmt.Sprintf("Error looking up ZIP code: %v", err), http.StatusInternalServerError)
		return
	}
	if providedLocation == nil {
		http.Error(w, "No location found", http.StatusNotFound)
		return
	}

	boundingMinLat, boundingMaxLat, boundingMinLon, boundingMaxLon := location_package.GetBoundingBox(providedLocation.Latitude, providedLocation.Longitude, miles)

	locations, err := app.store.Locations.GetLocationsByBoundingBox(r.Context(), boundingMinLat, boundingMaxLat, boundingMinLon, boundingMaxLon)
	if err != nil {
		fmt.Println("getNearbyByZip: Error fetching locations by bounding box:", err)
		http.Error(w, fmt.Sprintf("Error fetching locations: %v", err), http.StatusInternalServerError)
		return
	}
	if err := utils.JsonResponse(w, http.StatusOK, locations); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

}

// LookupByZip godoc
//
//	@Summary		Get nearby locations by ZIP code
//	@Description	Get nearby locations by ZIP code
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param			ZIPCode	path		string	true	"ZIP code"
//	@Param			miles	path		string	true	"Distance in miles"
//	@Success		200		{array}		store.Location
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/locations/zip/nearby/{ZIPCode}/{miles} [get]
func (app *application) getNearbyByZipHandler(w http.ResponseWriter, r *http.Request) { // Debugging - remove
	zip := chi.URLParam(r, "ZIPCode")
	miles := chi.URLParam(r, "miles")
	// Convert miles to float64
	milesFloat, err := strconv.ParseFloat(miles, 64)
	if err != nil {
		http.Error(w, "Invalid miles parameter", http.StatusBadRequest)
		return
	}

	fmt.Println("getNearbyByZipHandler: ZIP code:", zip)
	if zip == "" {
		http.Error(w, "ZIP code is required", http.StatusBadRequest)
		return
	}
	// Default distance in miles
	app.getNearbyByZip(w, r, zip, milesFloat)
}
