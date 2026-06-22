// Package latlon2h3 converts latitude/longitude coordinates to Uber H3 hexagon
// indices at resolution 8.
//
// H3 is a hexagonal grid system created by Uber for geospatial analysis.
// Resolution 8 cells have an average area of ~0.74 km², suitable
// for neighborhood-level spatial aggregation.
//
// For full H3 functionality (cell properties, neighbors, reverse lookup, etc.),
// use the upstream library directly: github.com/uber/h3-go/v4
//
// Usage:
//
//	h3Index, err := latlon2h3.ToH3(-23.5505, -46.6333) // São Paulo, Brazil
//	h3Index := latlon2h3.MustToH3(40.7128, -74.0060)    // panics on invalid input
package latlon2h3

import (
	"fmt"

	"github.com/uber/h3-go/v4"
)

// DefaultResolution is the H3 resolution used by ToH3 and MustToH3
// when no explicit resolution is provided.
const DefaultResolution = 8

// ToH3 converts a latitude/longitude pair to the H3 cell index string
// at the default resolution (8).
//
// Latitude must be in [-90, 90] and longitude in [-180, 180].
// Returns the H3 index as a hexadecimal string (e.g., "88283080b5fffff").
func ToH3(lat, lon float64) (string, error) {
	return ToH3AtResolution(lat, lon, DefaultResolution)
}

// ToH3AtResolution converts a latitude/longitude pair to an H3 cell index
// string at the specified resolution (0-15).
func ToH3AtResolution(lat, lon float64, resolution int) (string, error) {
	if !validCoord(lat, lon) {
		return "", fmt.Errorf("latlon2h3: invalid coordinates lat=%.6f lon=%.6f", lat, lon)
	}

	geo := h3.NewLatLng(lat, lon)
	cell, err := h3.LatLngToCell(geo, resolution)
	if err != nil {
		return "", fmt.Errorf("latlon2h3: %w", err)
	}

	return cell.String(), nil
}

// MustToH3 is like ToH3 but panics on invalid input.
// Use only when coordinates are guaranteed valid.
func MustToH3(lat, lon float64) string {
	h, err := ToH3(lat, lon)
	if err != nil {
		panic(err)
	}
	return h
}

// IsValidLatLon checks whether the given latitude/longitude pair is valid.
func IsValidLatLon(lat, lon float64) bool {
	return validCoord(lat, lon)
}

// validCoord checks latitude in [-90, 90] and longitude in [-180, 180].
func validCoord(lat, lon float64) bool {
	return lat >= -90 && lat <= 90 && lon >= -180 && lon <= 180
}
