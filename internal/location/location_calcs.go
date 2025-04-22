package location

import "math"

func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371.0 // Radius of Earth in kilometers

	dLat := (lat2 - lat1) * (math.Pi / 180)
	dLon := (lon2 - lon1) * (math.Pi / 180)

	lat1Rad := lat1 * (math.Pi / 180)
	lat2Rad := lat2 * (math.Pi / 180)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c // in kilometers
}

// GetBoundingBox returns a lat/lon bounding box around a point, radius in miles
func GetBoundingBox(lat, lon, radiusMiles float64) (minLat, maxLat, minLon, maxLon float64) {
	// Convert radius from miles to degrees
	latDelta := radiusMiles / 69.0 // ~69 miles per degree of latitude
	lonDelta := radiusMiles / (math.Cos(lat*math.Pi/180.0) * 69.0)

	minLat = lat - latDelta
	maxLat = lat + latDelta
	minLon = lon - lonDelta
	maxLon = lon + lonDelta

	return
}
