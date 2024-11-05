package models

import (
	"fmt"
	"math"
)

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

func NewCoordinates(latitude, longitude float64) Coordinates {
	return Coordinates{
		Latitude:  latitude,
		Longitude: longitude,
	}
}

func (c Coordinates) String() string {
	return fmt.Sprintf("%f,%f", c.Latitude, c.Longitude)
}

// CalculateDistance calculates the geographic distance between two points
// on the earth using the haversine formula.
//
// The function takes two Coordinates, from and to, and returns the distance
// in kilometers as a float64.
func (c Coordinates) KmDistanceTo(to Coordinates) float64 {
	const earthRadiusKm = 6371.01

	latFrom := toRadians(c.Latitude)
	lonFrom := toRadians(c.Longitude)

	latTo := toRadians(to.Latitude)
	lonTo := toRadians(to.Longitude)

	// Haversine formula, https://en.wikipedia.org/wiki/Haversine_formula
	diffLat := latTo - latFrom
	diffLon := lonTo - lonFrom

	a := math.Pow(math.Sin(diffLat/2), 2) +
		math.Cos(latFrom)*math.Cos(latTo)*
			math.Pow(math.Sin(diffLon/2), 2)

	C := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * C
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
