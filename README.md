# latlon2h3

Go library to convert latitude/longitude coordinates to [Uber H3](https://h3geo.org/) hexagon indices at resolution 8 and back.

Resolution 8 cells average ~0.74 km² — great for neighborhood-level spatial aggregation.

## Installation

```bash
# System dependency (macOS)
brew install h3

# Go module
go get github.com/vitoravila/latlon2h3
```

## Quick start

```go
package main

import (
    "fmt"
    "github.com/vitoravila/latlon2h3"
)

func main() {
    // São Paulo → H3
    h3Index, err := latlon2h3.ToH3(-23.5505, -46.6333)
    if err != nil {
        panic(err)
    }
    fmt.Println(h3Index) // "88a8100c03fffff"

    // Custom resolution
    h3Res6, _ := latlon2h3.ToH3AtResolution(-23.5505, -46.6333, 6)

    // H3 → lat/lon (cell center)
    lat, lon, _ := latlon2h3.FromH3(h3Index)
    fmt.Printf("%.4f, %.4f\n", lat, lon) // -23.5530, -46.6350

    // Get the 6 neighboring cells
    neighbors, _ := latlon2h3.Neighbors(h3Index)
    fmt.Println(len(neighbors)) // 6

    // MustToH3 panics on invalid input — use when coords are guaranteed valid
    h := latlon2h3.MustToH3(40.7128, -74.0060)
    fmt.Println(h)
}
```

## API

```go
func ToH3(lat, lon float64) (string, error)
func ToH3AtResolution(lat, lon float64, resolution int) (string, error)
func MustToH3(lat, lon float64) string
func FromH3(h3Index string) (lat, lon float64, err error)
func Resolution(h3Index string) (int, error)
func Neighbors(h3Index string) ([]string, error)
func IsValidLatLon(lat, lon float64) bool
```

## Running tests

```bash
go test -v -race ./...
```
