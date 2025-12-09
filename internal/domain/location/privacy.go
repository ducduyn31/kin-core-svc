package location

import (
	"math"
)

func haversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371000 // meters

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func FuzzLocation(lat, lng float64, precision Precision) (float64, float64) {
	switch precision {
	case PrecisionExact:
		return lat, lng
	case PrecisionNeighborhood:
		return math.Round(lat*1000) / 1000, math.Round(lng*1000) / 1000
	case PrecisionCity:
		return math.Round(lat*100) / 100, math.Round(lng*100) / 100
	case PrecisionCountry:
		return math.Round(lat), math.Round(lng)
	default:
		return lat, lng
	}
}

type Precision string

const (
	PrecisionExact        Precision = "exact"
	PrecisionNeighborhood Precision = "neighborhood"
	PrecisionCity         Precision = "city"
	PrecisionCountry      Precision = "country"
)

func IsValidPrecision(p Precision) bool {
	switch p {
	case PrecisionExact, PrecisionNeighborhood, PrecisionCity, PrecisionCountry:
		return true
	default:
		return false
	}
}
