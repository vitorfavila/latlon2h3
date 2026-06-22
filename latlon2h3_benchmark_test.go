package latlon2h3_test

import (
	"math/rand"
	"testing"

	"github.com/vitoravila/latlon2h3"
)

// randomCoord returns a valid random lat/lon pair.
func randomCoord(rng *rand.Rand) (float64, float64) {
	lat := rng.Float64()*180 - 90   // [-90, 90]
	lon := rng.Float64()*360 - 180  // [-180, 180]
	return lat, lon
}

// BenchmarkToH3 measures per-op latency of a single ToH3 call.
// Use -benchtime to control volume:
//
//	go test -bench=BenchmarkToH3$ -benchtime=100000x -benchmem ./...
//	go test -bench=BenchmarkToH3$ -benchtime=1000000x -benchmem ./...
func BenchmarkToH3(b *testing.B) {
	rng := rand.New(rand.NewSource(42))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lat, lon := randomCoord(rng)
		_, _ = latlon2h3.ToH3(lat, lon)
	}
}

// BenchmarkToH3_Parallel measures throughput with all CPU cores.
func BenchmarkToH3_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		rng := rand.New(rand.NewSource(42))
		for pb.Next() {
			lat, lon := randomCoord(rng)
			_, _ = latlon2h3.ToH3(lat, lon)
		}
	})
}
