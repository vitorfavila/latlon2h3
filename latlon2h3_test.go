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

			// Resolution 8 H3 indices are always 15 hex chars.
			if len(h3Index) != 15 {
				t.Errorf("expected 15-char H3 string at res 8, got %q (%d chars)", h3Index, len(h3Index))
			}

			t.Logf("%s: (%.4f, %.4f) → %s", tt.name, tt.lat, tt.lon, h3Index)
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

func TestToH3AtResolution(t *testing.T) {
	// Same coords should produce different cells at different resolutions.
	res6, err := latlon2h3.ToH3AtResolution(-23.5505, -46.6333, 6)
	if err != nil {
		t.Fatalf("res 6: %v", err)
	}
	res8, err := latlon2h3.ToH3AtResolution(-23.5505, -46.6333, 8)
	if err != nil {
		t.Fatalf("res 8: %v", err)
	}
	if res6 == res8 {
		t.Errorf("cells at different resolutions should differ, got same: %s", res6)
	}
}

func TestIsValidLatLon(t *testing.T) {
	if !latlon2h3.IsValidLatLon(0, 0) {
		t.Error("0,0 should be valid")
	}
	if !latlon2h3.IsValidLatLon(90, 180) {
		t.Error("90,180 should be valid")
	}
	if !latlon2h3.IsValidLatLon(-90, -180) {
		t.Error("-90,-180 should be valid")
	}
	if latlon2h3.IsValidLatLon(90.1, 0) {
		t.Error("90.1,0 should be invalid")
	}
	if latlon2h3.IsValidLatLon(0, 180.1) {
		t.Error("0,180.1 should be invalid")
	}
}

func TestMustToH3_PanicsOnInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustToH3 should panic on invalid coords")
		}
	}()
	latlon2h3.MustToH3(91, 0)
}

func TestMustToH3_Valid(t *testing.T) {
	h := latlon2h3.MustToH3(-23.5505, -46.6333)
	if h == "" {
		t.Error("MustToH3 returned empty string")
	}
}
