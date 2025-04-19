package nominatim

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type NominatimResult struct {
	DisplayName string `json:"display_name"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	Address     struct {
		City    string `json:"city"`
		State   string `json:"state"`
		Zip     string `json:"postcode"`
		Country string `json:"country"`
	} `json:"address"`
}

func LookupByZip(zip string) (*NominatimResult, error) {
	baseURL := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Add("postalcode", zip)
	params.Add("country", "USA")
	params.Add("format", "json")
	params.Add("addressdetails", "1")
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// 1-second delay to comply with rate limit
	time.Sleep(1 * time.Second)

	// Set a proper User-Agent (required by Nominatim)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ShotSeek/1.0 (you@example.com)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var results []NominatimResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found")
	}

	return &results[0], nil
}
