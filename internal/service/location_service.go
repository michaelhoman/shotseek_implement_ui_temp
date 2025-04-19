package service

import (
	"github.com/michaelhoman/ShotSeek/internal/location"
	"github.com/michaelhoman/ShotSeek/internal/store"
)

// LocationService handles location-based operations
type LocationService struct{}

// NewLocationService creates a new LocationService instance
func NewLocationService() *LocationService {
	return &LocationService{}
}

// ConvertFromNominatim handles converting a Nominatim result to a store.Location
func (s *LocationService) ConvertFromNominatim(nominatimResult *location.NominatimResult) (*store.Location, error) {
	return location.LocationFromNominatim(nominatimResult)
}

// export the type NominatimResult struct {
// 	DisplayName string `json:"display_name"`
// 	Lat         string `json:"lat"`
// 	Lon         string `json:"lon"`
// 	Address     struct {
// 		City        string `json:"city"`
// 		Town        string `json:"town"`
// 		Village     string `json:"village"`
// 		State       string `json:"state"`
// 		County      string `json:"county"`
// 		Postcode    string `json:"postcode"`
// 		Country     string `json:"country"`
// 		CountryCode string `json:"country_code"`
// 	} `json:"address"`
// } from package location
