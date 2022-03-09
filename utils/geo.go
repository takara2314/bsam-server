package utils

import "math"

const (
	EarthRadiusKm = 6371
)

// CalcDistanceAtoBEarth calculates [A] to [B] distance in the Earth.
func CalcDistanceAtoBEarth(latA float64, lngA float64, latB float64, lngB float64) float64 {
	dLat := DegToRad(latB - latA)
	dLng := DegToRad(lngB - lngA)

	a := math.Pow(math.Sin(dLat/2), 2.0) +
		math.Pow(math.Sin(dLng/2), 2.0)*
			math.Cos(latA)*
			math.Cos(latB)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadiusKm * c * 1000
}
