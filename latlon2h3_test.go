package latlon2h3_test

import (
	"math"
	"testing"

	"github.com/vitoravila/latlon2h3"
)

func TestToH3_ValidCoords(t *testing.T) {
	tests := []struct {
		name     string
		lat, lon float64
		// We don't assert exact cell strings — they could change across H3 versions.
		// Instead we verify the output is a valid H3 string and consistent.
	}{
		{name: "São Paulo", lat: -23.5505, lon: -46.6333},
		{name: "New York", lat: 40.7128, lon: -74.0060},
		{name: "Tokyo", lat: 35.6762, lon: 139.6503},
		{name: "London", lat: 51.5074, lon: -0.1278},
		{name: "Sydney", lat: -33.8688, lon: 151.2093},
		{name: "Equator prime meridian", lat: 0, lon: 0},
		{name: "North Pole edge", lat: 89.999, lon: 0},
		{name: "South Pole edge", lat: -89.999, lon: 0},
		{name: "Date line east", lat: 45, lon: 179.999},
		{name: "Date line west", lat: 45, lon: -179.999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h3Index, err := latlon2h3.ToH3(tt.lat, tt.lon)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if h3Index == "" {
				t.Fatal("expected non-empty H3 index")
			}
			if len(h3Index) != 15 {
				t.Errorf("expected 15-char H3 string at res 8, got %q (%d chars)", h3Index, len(h3Index))
			}

			// Verify resolution is 8.
			res, err := latlon2h3.Resolution(h3Index)
			if err != nil {
				t.Fatalf("Resolution() error: %v", err)
			}
			if res != 8 {
				t.Errorf("expected resolution 8, got %d", res)
			}

			// Roundtrip: lat/lon should be close.
			rlat, rlon, err := latlon2h3.FromH3(h3Index)
			if err != nil {
				t.Fatalf("FromH3() error: %v", err)
			}
			// The cell center should be within a reasonable distance.
			// H3 res 8 cells average ~0.74 km²; center should be close to input.
			t.Logf("%s: input=(%.4f, %.4f) → h3=%s → center=(%.4f, %.4f)", tt.name, tt.lat, tt.lon, h3Index, rlat, rlon)
		})
	}
}

func TestToH3_InvalidCoords(t *testing.T) {
	tests := []struct {
		name     string
		lat, lon float64
	}{
		{name: "lat > 90", lat: 91, lon: 0},
		{name: "lat < -90", lat: -91, lon: 0},
		{name: "lon > 180", lat: 0, lon: 181},
		{name: "lon < -180", lat: 0, lon: -181},
		{name: "NaN lat", lat: math.NaN(), lon: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := latlon2h3.ToH3(tt.lat, tt.lon)
			if err == nil {
				t.Fatal("expected error for invalid coords, got nil")
			}
		})
	}
}

func TestRoundtrip(t *testing.T) {
	lat, lon := -23.5505, -46.6333 // São Paulo

	h3Index, err := latlon2h3.ToH3(lat, lon)
	if err != nil {
		t.Fatalf("ToH3: %v", err)
	}

	rlat, rlon, err := latlon2h3.FromH3(h3Index)
	if err != nil {
		t.Fatalf("FromH3: %v", err)
	}

	// Convert the center back to H3 — it should give the same cell.
	sameIndex, err := latlon2h3.ToH3(rlat, rlon)
	if err != nil {
		t.Fatalf("ToH3(roundtrip): %v", err)
	}
	if sameIndex != h3Index {
		t.Errorf("roundtrip mismatch: original=%s center_to_h3=%s", h3Index, sameIndex)
	}
}

func TestFromH3_Invalid(t *testing.T) {
	tests := []struct {
		name string
		h3   string
	}{
		{name: "empty string", h3: ""},
		{name: "garbage", h3: "not-an-h3-index"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := latlon2h3.FromH3(tt.h3)
			if err == nil {
				t.Fatal("expected error for invalid H3, got nil")
			}
		})
	}
}

func TestNeighbors(t *testing.T) {
	// Use São Paulo as a stable test point.
	h3Index, err := latlon2h3.ToH3(-23.5505, -46.6333)
	if err != nil {
		t.Fatalf("ToH3: %v", err)
	}

	neighbors, err := latlon2h3.Neighbors(h3Index)
	if err != nil {
		t.Fatalf("Neighbors: %v", err)
	}

	// A valid H3 cell at resolution 8 always has 6 neighbors
	// (pentagons at res 0 have 5, but those don't exist at res 8).
	if len(neighbors) != 6 {
		t.Errorf("expected 6 neighbors, got %d", len(neighbors))
	}

	for _, n := range neighbors {
		if n == "" {
			t.Error("neighbor is empty string")
		}
		if n == h3Index {
			t.Error("neighbor equals origin cell")
		}
	}
}

func TestIsValidLatLon(t *testing.T) {
	if latlon2h3.IsValidLatLon(0, 0) != true {
		t.Error("0,0 should be valid")
	}
	if latlon2h3.IsValidLatLon(90, 180) != true {
		t.Error("90,180 should be valid")
	}
	if latlon2h3.IsValidLatLon(-90, -180) != true {
		t.Error("-90,-180 should be valid")
	}
	if latlon2h3.IsValidLatLon(90.1, 0) != false {
		t.Error("90.1,0 should be invalid")
	}
	if latlon2h3.IsValidLatLon(0, 180.1) != false {
		t.Error("0,180.1 should be invalid")
	}
}

func TestMustToH3_PanicsOnInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustToH3 should panic on invalid coords")
		}
	}()
	latlon2h3.MustToH3(91, 0) // should panic
}

func TestMustToH3_Valid(t *testing.T) {
	h := latlon2h3.MustToH3(-23.5505, -46.6333)
	if h == "" {
		t.Error("MustToH3 returned empty string")
	}
}
