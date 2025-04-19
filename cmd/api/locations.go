package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

func (a *application) lookupByZip(ctx context.Context, tx *sql.Tx, zip string) (*store.Location, error) {
	fmt.Println("lookupByZip: Starting ZIP lookup for:", zip)
	if zip == "" {
		fmt.Println("lookupByZip: Empty ZIP code provided")
		return nil, fmt.Errorf("empty ZIP code provided")
	}
	location, err := a.store.Locations.GetGeneralLocationByZip(ctx, zip)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("\n\nErrNotFound - lookupByZip: No record found for ZIP:", zip) // Debugging line

			// Construct the URL for the Nominatim API lookup
			url := fmt.Sprintf("https://nominatim.openstreetmap.org/search?postalcode=%s&country=USA&format=json", zip)
			fmt.Println("lookupByZip: Constructed URL:", url)

			// Create the HTTP request
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println("lookupByZip: Error creating request:", err)
				return nil, err
			}

			// Set a custom User-Agent header
			userAgent := "ShotSeek/1.0 (contact@shotseek.app)"
			req.Header.Set("User-Agent", userAgent)
			fmt.Println("lookupByZip: Set User-Agent:", userAgent)

			// Send the HTTP request
			client := &http.Client{}
			fmt.Println("lookupByZip: Sending HTTP request...")
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("lookupByZip: HTTP request error:", err)
				return nil, err
			}
			defer func() {
				fmt.Println("lookupByZip: Closing response body")
				resp.Body.Close()
			}()

			// Handle non-OK status codes
			fmt.Println("lookupByZip: Received response with status:", resp.Status)
			if resp.StatusCode != http.StatusOK {
				fmt.Println("lookupByZip: Non-OK HTTP status:", resp.Status)
				return nil, fmt.Errorf("nominatim error: %s", resp.Status)
			}

			// Decode the JSON response into a slice of results
			var results []location_package.NominatimResult
			fmt.Println("lookupByZip: Decoding JSON response...")
			if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
				fmt.Println("lookupByZip: JSON decode error:", err)
				return nil, err
			}

			// Handle the case when no results are returned
			fmt.Println("lookupByZip: Number of results received:", len(results))
			if len(results) == 0 {
				fmt.Println("lookupByZip: No results found for ZIP:", zip)
				return nil, fmt.Errorf("no results found for ZIP %s", zip)
			}

			// Use the first  (if there are multiple			fmt.Println("lookupByZip: Returning first result:", results[0])

			loc, err := location_package.LocationFromNominatim(&results[0])

			if err != nil {
				fmt.Println("lookupByZip: Error converting Nominatim result to Location struct:", err)
				return nil, err
			}
			// Return the location struct
			fmt.Println("API lookup successful, location found:", loc)
			a.store.Locations.Create(ctx, tx, loc)
			return loc, nil
		}
		log.Println("lookupByZip: Error getting location from DB:", err)
		fmt.Println("lookupByZip: Error getting location from DB:", err) // Debugging line
		return nil, err

	}
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
	// Begin a transaction
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to begin transaction: %v", err), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // Ensure the transaction is rolled back in case of an error.

	location, err := app.lookupByZip(r.Context(), tx, zip)
	if err != nil {
		fmt.Println("zipLookupHandler: Error looking up ZIP code:", err)
	}
	if location == nil {
		http.Error(w, "No location found", http.StatusNotFound)
		return
	}
	utils.JsonResponse(w, http.StatusOK, location)
	// location, err := app.store.Locations.GetGeneralLocationByZip(ctx, zip)
	// if err != nil {
	// 	if err == store.ErrNotFound {
	// 		fmt.Println("zipLookupHandler: No record found for ZIP:", zip) // Debugging line
	// 		result, err := app.lookupByZip(ctx, r, zip)
	// 		if err != nil {
	// 			http.Error(w, fmt.Sprintf("Lookup error: %v", err), http.StatusInternalServerError)
	// 			return
	// 		}
	// 		utils.JsonResponse(w, http.StatusOK, result)
	// 	} else {
	// 		fmt.Println("zipLookupHandler: Error getting location from DB:", err) // Debugging line
	// 		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
	// 		return
	// 	}
	// } else {
	// 	fmt.Println("ZipCode Non-Precise Location Already Exists in DB: NOT Calling Nominatim API") // Debugging line
	// 	utils.JsonResponse(w, http.StatusOK, location)

	// }
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
