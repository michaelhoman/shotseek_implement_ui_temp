package store_test

import (
	"testing"

	"github.com/michaelhoman/ShotSeek/internal/location"
	"github.com/michaelhoman/ShotSeek/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestLocationFromNominatim(t *testing.T) {
	tests := []struct {
		name     string
		input    *location.NominatimResult
		expected *store.Location
		wantErr  bool
	}{
		{
			name: "valid Nominatim result",
			input: &location.NominatimResult{
				Lat: "38.8142294",
				Lon: "-94.9308321",
				Address: struct {
					City        string `json:"city"`
					Town        string `json:"town"`
					Village     string `json:"village"`
					State       string `json:"state"`
					County      string `json:"county"`
					Postcode    string `json:"postcode"`
					Country     string `json:"country"`
					CountryCode string `json:"country_code"`
				}{
					City:        "Gardner",
					State:       "Kansas",
					County:      "Johnson County",
					Postcode:    "66030",
					Country:     "United States",
					CountryCode: "US",
				},
			},
			expected: &store.Location{
				Street:      "",
				City:        "Gardner",
				State:       "Kansas",
				County:      "Johnson County",
				ZIPCode:     "66030",
				Country:     "United States",
				CountryCode: "US",
				Latitude:    38.8142294,
				Longitude:   -94.9308321,
			},
			wantErr: false,
		},
		{
			name: "invalid latitude and longitude",
			input: &location.NominatimResult{
				Lat: "invalid_lat",
				Lon: "invalid_lon",
				Address: struct {
					City        string `json:"city"`
					Town        string `json:"town"`
					Village     string `json:"village"`
					State       string `json:"state"`
					County      string `json:"county"`
					Postcode    string `json:"postcode"`
					Country     string `json:"country"`
					CountryCode string `json:"country_code"`
				}{
					City:        "Gardner",
					State:       "Kansas",
					County:      "Johnson County",
					Postcode:    "66030",
					Country:     "United States",
					CountryCode: "US",
				},
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := location.LocationFromNominatim(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
