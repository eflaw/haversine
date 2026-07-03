package geo

import "math"

// This file holds the *precomputable* per-point representations that make batched
// (many-pairs-from-a-fixed-set) distance cheap. They are the operational content
// of P1's stereographic-chordal identity: because a point's disk image depends on
// that point alone, it can be computed once and reused. See docs/BENCHMARK.md.

// FlatPoint is the precomputed stereographic disk image of a surface point.
// A pairwise ChordDistance then needs no trigonometric call at all.
type FlatPoint struct{ Re, Im, ModSq float64 }

// Flatten precomputes the disk image of a lat/lon point (paid once per point).
func Flatten(latDeg, lonDeg float64) FlatPoint {
	re, im, m := stereo(latDeg, lonDeg)
	return FlatPoint{re, im, m}
}

// ChordDistance returns the great-circle distance between two precomputed points.
// It performs zero forward trigonometry — only arithmetic, a sqrt and an atan2.
func ChordDistance(p, q FlatPoint, radius float64) float64 {
	du, dv := p.Re-q.Re, p.Im-q.Im
	a := (du*du + dv*dv) / ((1 + p.ModSq) * (1 + q.ModSq))
	return radius * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

// GeoPoint is the precomputed trigonometry of a surface point, for a cached
// haversine/law-of-cosines evaluation.
type GeoPoint struct{ SinLat, CosLat, Lon float64 }

// NewGeoPoint precomputes the per-point trig (paid once per point).
func NewGeoPoint(latDeg, lonDeg float64) GeoPoint {
	s, c := math.Sincos(latDeg * deg2rad)
	return GeoPoint{s, c, lonDeg * deg2rad}
}

// HavCached returns the great-circle distance between two precomputed points via
// the spherical law of cosines, kept well-conditioned by routing through
// a = (1-cos d)/2 and atan2. It differs from ChordDistance per pair by exactly
// one cosine, cos(delta-lambda).
func HavCached(p, q GeoPoint, radius float64) float64 {
	cosd := p.SinLat*q.SinLat + p.CosLat*q.CosLat*math.Cos(p.Lon-q.Lon)
	a := 0.5 * (1 - cosd)
	if a < 0 {
		a = 0
	} else if a > 1 {
		a = 1
	}
	return radius * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
