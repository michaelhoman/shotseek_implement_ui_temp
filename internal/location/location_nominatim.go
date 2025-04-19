package location

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/michaelhoman/ShotSeek/internal/store" // Correct path for store package
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// var results []NominatimResult

// NominatimResult represents the structure of a Nominatim API result
type NominatimResult struct {
	DisplayName string `json:"display_name"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	Address     struct {
		City        string `json:"city"`
		Town        string `json:"town"`
		Village     string `json:"village"`
		State       string `json:"state"`
		County      string `json:"county"`
		Postcode    string `json:"postcode"`
		Country     string `json:"country"`
		CountryCode string `json:"country_code"`
	} `json:"address"`
}

// LocationFromNominatim converts a NominatimResult into a Location struct
func LocationFromNominatim(n *NominatimResult) (*store.Location, error) {
	fmt.Println("LocationFromNominatim: Converting Nominatim result to Location struct\n\n", n)
	lat, err := strconv.ParseFloat(n.Lat, 64) // 🔧 was `result.Lat` — now fixed to `n.Lat`
	if err != nil {
		return nil, fmt.Errorf("invalid latitude: %w", err)
	}

	lon, err := strconv.ParseFloat(n.Lon, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid longitude: %w", err)
	}

	// Fallback logic: use City, then Town, then Village
	city := n.Address.City
	if city == "" {
		city = n.Address.Town
	}
	if city == "" {
		city = n.Address.Village
	}

	return &store.Location{
		Street:      "", // Not returned from zip-only lookups
		City:        city,
		State:       n.Address.State,
		County:      n.Address.County,
		ZIPCode:     n.Address.Postcode,
		Country:     n.Address.Country,
		CountryCode: strings.ToUpper(n.Address.CountryCode),
		Latitude:    lat,
		Longitude:   lon,
	}, nil
}
