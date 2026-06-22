# latlon2h3

Go library to convert latitude/longitude coordinates to [Uber H3](https://h3geo.org/) hexagon indices at resolution 8.

Resolution 8 cells average ~0.74 km² — great for neighborhood-level spatial aggregation.

For full H3 functionality (cell properties, neighbors, reverse lookup), use [uber/h3-go](https://github.com/uber/h3-go) directly.

## Prerequisites

This library depends on the H3 C library via CGo. You need `libh3` installed in your build environment.

### macOS

```bash
brew install h3
```

### Linux (Debian/Ubuntu)

```bash
sudo apt-get install libh3-dev
```

### Docker

If you prefer not to install native dependencies, use the included Docker image:

```bash
# Build the image
docker build -t latlon2h3 .

# Run tests
docker run --rm -v "$PWD":/src latlon2h3 go test -v -race ./...

# Run benchmarks
docker run --rm -v "$PWD":/src latlon2h3 \
    go test -bench='BenchmarkToH3$' -benchtime=100000x -benchmem ./...
```

### Docker for apps importing this library

For apps that import `latlon2h3` (or `uber/h3-go` directly), use a multi-stage build:

```dockerfile
# ---- build stage ----
FROM golang:1.22-bookworm AS build

RUN apt-get update && apt-get install -y --no-install-recommends \
    libh3-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o /bin/app ./cmd/app

# ---- runtime stage ----
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    libh3-4 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /bin/app /usr/local/bin/app
ENTRYPOINT ["/usr/local/bin/app"]
```

Key points:
- Build stage needs `libh3-dev` (headers + shared lib)
- Runtime stage only needs `libh3-4` (shared lib only, no headers)
- `CGO_ENABLED=1` must be set explicitly (some Go images default to 0)

## Installation

```bash
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
