# latlon2h3

Go library to convert latitude/longitude coordinates to [Uber H3](https://h3geo.org/) hexagon indices at resolution 8.

Resolution 8 cells average ~0.74 km² — great for neighborhood-level spatial aggregation.

For full H3 functionality (cell properties, neighbors, reverse lookup), use [uber/h3-go](https://github.com/uber/h3-go) directly.

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
    // São Paulo, Brazil → H3
    h3Index, err := latlon2h3.ToH3(-23.5505, -46.6333)
    if err != nil {
        panic(err)
    }
    fmt.Println(h3Index) // "88a8100c03fffff"

    // Custom resolution (0-15)
    h3Res6, _ := latlon2h3.ToH3AtResolution(-23.5505, -46.6333, 6)

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
func IsValidLatLon(lat, lon float64) bool
```

## Running tests

```bash
go test -v -race ./...
```

## Benchmarks

Measured on Apple Silicon M3 Pro (8 performance cores), Go 1.22, H3 4.5.0.

```bash
# 100K ops — single core
go test -bench='BenchmarkToH3$' -benchtime=100000x -benchmem ./...

# 1M ops — single core
go test -bench='BenchmarkToH3$' -benchtime=1000000x -benchmem ./...

# 1M ops — all cores
go test -bench='BenchmarkToH3_Parallel' -benchtime=1000000x -benchmem ./...
```

| Benchmark | Ops | Time/op | Throughput | Memory |
|---|---|---|---|---|
| ToH3 (single core) | 100K | 865 ns | ~1.16 M ops/s | 40 B / 3 allocs |
| ToH3 (single core) | 1M | 629 ns | ~1.59 M ops/s | 40 B / 3 allocs |
| ToH3 (parallel, 8 cores) | 1M | 149 ns | ~6.71 M ops/s | 40 B / 3 allocs |

The 3 allocations per call come from the CGo bridge (LatLng struct, return string).
Parallel throughput scales nearly linearly with available cores.
