package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/michaelhoman/ShotSeek/internal/utils"
)

type NominatimResult struct {
	DisplayName string   `json:"display_name"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	Address     Address  `json:"address"`
	BoundingBox []string `json:"boundingbox"`
	OSMType     string   `json:"osm_type"`
	OSMID       string   `json:"osm_id"`
	Class       string   `json:"class"`
	Type        string   `json:"type"`
	Importance  float64  `json:"importance"`
	PlaceID     int      `json:"place_id"`
	License     string   `json:"license"`
}

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

func lookupByZip(zip string) (*NominatimResult, error) {
	fmt.Println("lookupByZip: Starting ZIP lookup for:", zip)

	url := fmt.Sprintf("https://nominatim.openstreetmap.org/search?postalcode=%s&country=USA&format=json", zip)
	fmt.Println("lookupByZip: Constructed URL:", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("lookupByZip: Error creating request:", err)
		return nil, err
	}

	userAgent := "ShotSeek/1.0 (contact@shotseek.app)"
	req.Header.Set("User-Agent", userAgent)
	fmt.Println("lookupByZip: Set User-Agent:", userAgent)

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

	fmt.Println("lookupByZip: Received response with status:", resp.Status)

	if resp.StatusCode != http.StatusOK {
		fmt.Println("lookupByZip: Non-OK HTTP status:", resp.Status)
		return nil, fmt.Errorf("nominatim error: %s", resp.Status)
	}

	var results []NominatimResult
	fmt.Println("lookupByZip: Decoding JSON response...")
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		fmt.Println("lookupByZip: JSON decode error:", err)
		return nil, err
	}

	fmt.Println("lookupByZip: Number of results received:", len(results))
	if len(results) == 0 {
		fmt.Println("lookupByZip: No results found for ZIP:", zip)
		return nil, fmt.Errorf("no results found for ZIP %s", zip)
	}

	fmt.Println("lookupByZip: Returning first result:", results[0])
	return &results[0], nil
}

// LookupByZip godoc
//
//	@Summary		Lookup location by ZIP code
//	@Description	Lookup location by ZIP code
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param			ZIPCode	path		string	true	"ZIP code"
//	@Success		200		{object}	NominatimResult
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

	result, err := lookupByZip(zip)
	if err != nil {
		http.Error(w, fmt.Sprintf("Lookup error: %v", err), http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, http.StatusOK, result)
}
