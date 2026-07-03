// Package geo computes great-circle distance two independent ways and checks
// them against a third, independent reference.
//
//   Haversine  — the formula derived in the paper "From the Unit Circle to the
//                Haversine" (P2): a = sin^2(dphi/2) + cos(phi1)cos(phi2)sin^2(dlam/2).
//   Chordal    — the same distance read off the stereographic disk image of each
//                point, via the Riemann-sphere chordal metric (P1's construction).
//                Chordal(...) == Haversine(...) exactly; see docs/MATH.md.
//   CentralAngle — the independent reference d = arccos(u1 . u2) used by P2's
//                  numerical-check remark.
//
// All angles are in degrees on input. Distances are returned in the same unit as
// the radius argument; EarthRadiusKm reproduces the paper's reference value.
package geo

import "math"

// EarthRadiusKm is the reference sphere radius used in the companion paper's
// numerical check (London–New York ~ 5570 km at this radius).
const EarthRadiusKm = 6371.0

const deg2rad = math.Pi / 180

// unit returns the 3D unit position vector of a lat/lon point (degrees).
func unit(latDeg, lonDeg float64) (x, y, z float64) {
	phi := latDeg * deg2rad
	lam := lonDeg * deg2rad
	sLam, cLam := math.Sincos(lam)
	sPhi, cPhi := math.Sincos(phi)
	return cPhi * cLam, cPhi * sLam, sPhi
}

// stereo returns the stereographic disk image w = tan(pi/4 - phi/2) e^{i*lambda}
// as (Re w, Im w, |w|^2). North pole -> 0; south pole -> infinity (see caveat in
// Chordal).
func stereo(latDeg, lonDeg float64) (re, im, modSq float64) {
	phi := latDeg * deg2rad
	lam := lonDeg * deg2rad
	r := math.Tan(math.Pi/4 - phi/2)
	sLam, cLam := math.Sincos(lam)
	return r * cLam, r * sLam, r * r
}

// haversineAngle returns the central angle (radians) via the haversine formula,
// in its atan2 form (accurate for all separations including antipodal).
func haversineAngle(lat1, lon1, lat2, lon2 float64) float64 {
	la1 := lat1 * deg2rad
	la2 := lat2 * deg2rad
	sLat := math.Sin((lat2 - lat1) * deg2rad / 2)
	sLon := math.Sin((lon2 - lon1) * deg2rad / 2)
	a := sLat*sLat + math.Cos(la1)*math.Cos(la2)*sLon*sLon
	return 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

// chordalAngle returns the central angle (radians) from the stereographic disk
// images via the chordal metric: a = |w1-w2|^2 / ((1+|w1|^2)(1+|w2|^2)).
func chordalAngle(lat1, lon1, lat2, lon2 float64) float64 {
	u1, v1, m1 := stereo(lat1, lon1)
	u2, v2, m2 := stereo(lat2, lon2)
	du, dv := u1-u2, v1-v2
	a := (du*du + dv*dv) / ((1 + m1) * (1 + m2))
	return 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

// Haversine returns the great-circle distance (in the unit of radius) between two
// lat/lon points, using the haversine formula of P2.
func Haversine(lat1, lon1, lat2, lon2, radius float64) float64 {
	return radius * haversineAngle(lat1, lon1, lat2, lon2)
}

// Chordal returns the great-circle distance (in the unit of radius) via the
// stereographic disk image and the chordal metric of P1. It equals Haversine
// exactly (docs/MATH.md).
//
// Caveat: the stereographic image blows up at the south pole (the projection
// antipode, the one point P1's flattening excludes), so Chordal loses precision
// for points within a tiny neighbourhood of the south pole. For a general-purpose
// distance, prefer Haversine; Chordal is here to demonstrate the identity and to
// expose the precomputable coordinate used by the batch functions.
func Chordal(lat1, lon1, lat2, lon2, radius float64) float64 {
	return radius * chordalAngle(lat1, lon1, lat2, lon2)
}

// CentralAngle returns the great-circle central angle (radians) via the
// independent dot-product route d = arccos(u1 . u2). This is the reference used
// by P2's numerical-check remark. Note arccos itself loses precision for very
// small angles (nearby points) — which is exactly why the haversine form is
// preferred in practice — so this is a reference for well-separated pairs, not a
// recommended production formula.
func CentralAngle(lat1, lon1, lat2, lon2 float64) float64 {
	x1, y1, z1 := unit(lat1, lon1)
	x2, y2, z2 := unit(lat2, lon2)
	dot := x1*x2 + y1*y2 + z1*z2
	if dot > 1 {
		dot = 1
	} else if dot < -1 {
		dot = -1
	}
	return math.Acos(dot)
}
