package geo

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

var (
	pts   [][2]float64
	pairs [][4]float64
	sink  float64
)

func init() {
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 1024; i++ {
		lat := rng.Float64()*178 - 89 // (-89, 89): avoid the exact poles
		lon := rng.Float64()*360 - 180
		pts = append(pts, [2]float64{lat, lon})
	}
	for i := 0; i < 1024; i++ {
		a := pts[rng.Intn(len(pts))]
		c := pts[rng.Intn(len(pts))]
		pairs = append(pairs, [4]float64{a[0], a[1], c[0], c[1]})
	}
}

// TestLondonNewYork reproduces the reference distance quoted in P2's remark.
func TestLondonNewYork(t *testing.T) {
	const wantKm = 5570.2222
	got := Haversine(51.5074, -0.1278, 40.7128, -74.0060, EarthRadiusKm)
	if math.Abs(got-wantKm) > 0.01 {
		t.Fatalf("London-NY = %.4f km, want %.4f km", got, wantKm)
	}
}

// TestAgreement is the paper's numerical check: haversine and chordal must both
// match the independent arccos reference, and each other, over random pairs.
func TestAgreement(t *testing.T) {
	const tolRad = 1e-9 // paper reports ~1e-13; 1e-9 is a comfortable guard
	for _, p := range pairs {
		ref := CentralAngle(p[0], p[1], p[2], p[3])
		hav := Haversine(p[0], p[1], p[2], p[3], 1) // radius 1 => angle in radians
		cho := Chordal(p[0], p[1], p[2], p[3], 1)
		if math.Abs(hav-ref) > tolRad {
			t.Errorf("haversine vs reference: %.3e rad at %v", hav-ref, p)
		}
		if math.Abs(cho-ref) > tolRad {
			t.Errorf("chordal vs reference: %.3e rad at %v", cho-ref, p)
		}
		if math.Abs(hav-cho) > tolRad {
			t.Errorf("haversine vs chordal: %.3e rad at %v", hav-cho, p)
		}
	}
}

// TestCachedMatch checks the batch representations agree with the scalar ones.
func TestCachedMatch(t *testing.T) {
	const tolKm = 1e-6
	for _, p := range pairs {
		base := Haversine(p[0], p[1], p[2], p[3], EarthRadiusKm)
		ch := ChordDistance(Flatten(p[0], p[1]), Flatten(p[2], p[3]), EarthRadiusKm)
		hc := HavCached(NewGeoPoint(p[0], p[1]), NewGeoPoint(p[2], p[3]), EarthRadiusKm)
		if math.Abs(ch-base) > tolKm {
			t.Errorf("ChordDistance off by %.3e km at %v", ch-base, p)
		}
		if math.Abs(hc-base) > tolKm {
			t.Errorf("HavCached off by %.3e km at %v", hc-base, p)
		}
	}
}

// --- single-pair benchmarks (indexing prevents constant-folding) ---

func BenchmarkPair_Haversine(b *testing.B) {
	var s float64
	for i := 0; i < b.N; i++ {
		p := pairs[i&1023]
		s += Haversine(p[0], p[1], p[2], p[3], EarthRadiusKm)
	}
	sink = s
}

func BenchmarkPair_Chordal(b *testing.B) {
	var s float64
	for i := 0; i < b.N; i++ {
		p := pairs[i&1023]
		s += Chordal(p[0], p[1], p[2], p[3], EarthRadiusKm)
	}
	sink = s
}

// --- all-pairs (O(n^2)) batch benchmarks, end to end incl. precompute ---

func allPairsHaversine(n int) (s float64) {
	p := pts[:n]
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			s += Haversine(p[i][0], p[i][1], p[j][0], p[j][1], EarthRadiusKm)
		}
	}
	return
}

func allPairsChordalNoCache(n int) (s float64) {
	p := pts[:n]
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			s += Chordal(p[i][0], p[i][1], p[j][0], p[j][1], EarthRadiusKm)
		}
	}
	return
}

func allPairsChordalCached(n int) (s float64) {
	p := pts[:n]
	fp := make([]FlatPoint, n)
	for i := 0; i < n; i++ {
		fp[i] = Flatten(p[i][0], p[i][1])
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			s += ChordDistance(fp[i], fp[j], EarthRadiusKm)
		}
	}
	return
}

func allPairsHavCached(n int) (s float64) {
	p := pts[:n]
	gp := make([]GeoPoint, n)
	for i := 0; i < n; i++ {
		gp[i] = NewGeoPoint(p[i][0], p[i][1])
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			s += HavCached(gp[i], gp[j], EarthRadiusKm)
		}
	}
	return
}

func BenchmarkAllPairs(b *testing.B) {
	for _, n := range []int{64, 256, 1024} {
		n := n
		b.Run(fmt.Sprintf("Haversine/N=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			var s float64
			for it := 0; it < b.N; it++ {
				s += allPairsHaversine(n)
			}
			sink = s
		})
		b.Run(fmt.Sprintf("ChordalNoCache/N=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			var s float64
			for it := 0; it < b.N; it++ {
				s += allPairsChordalNoCache(n)
			}
			sink = s
		})
		b.Run(fmt.Sprintf("ChordalCached/N=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			var s float64
			for it := 0; it < b.N; it++ {
				s += allPairsChordalCached(n)
			}
			sink = s
		})
		b.Run(fmt.Sprintf("HavCached/N=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			var s float64
			for it := 0; it < b.N; it++ {
				s += allPairsHavCached(n)
			}
			sink = s
		})
	}
}
