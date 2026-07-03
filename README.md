# haversine-check

Companion code for the *Circle &amp; Sphere* papers. It computes great-circle
distance **two independent ways** and checks them against a **third** independent
reference:

- **Haversine** — the formula derived in *From the Unit Circle to the Haversine*
  (P2): `a = sin²(Δφ/2) + cos φ₁ cos φ₂ sin²(Δλ/2)`, `d = 2·atan2(√a, √(1−a))`.
- **Chordal** — the *same* distance read off the stereographic disk image of each
  point, via the Riemann-sphere chordal metric of *The Sphere on the Complex Disk*
  (P1): `a = |w₁−w₂|² / ((1+|w₁|²)(1+|w₂|²))`, with
  `w = tan(π/4 − φ/2)·e^{iλ}`.
- **Reference** — the independent central angle `d = arccos(u₁·u₂)` used by P2's
  numerical-check remark.

The point of the repo: the two formulas are provably the *same number*
(see [`docs/MATH.md`](docs/MATH.md)), and this code is the reproducible evidence —
the check P2's remark points to, plus the worked "both methods" example the papers
reference but don't print in full.

## Quick start

```sh
# the paper's numerical check (both methods vs the reference, over random pairs)
go test ./...

# the worked example — distance computed BOTH ways, side by side
go run ./cmd/distance
#   -> London -> New York, ~5570.22 km, haversine == chordal to ~1e-13

# any two points (lat1 lon1 lat2 lon2), optional radius with -r
go run ./cmd/distance 51.5074 -0.1278 40.7128 -74.0060
go run ./cmd/distance -r 6378.137 -33.8688 151.2093 40.7128 -74.0060

# performance comparison (see docs/BENCHMARK.md)
go test -bench=. -benchmem -run=^$ ./...
```

## What the test asserts

`TestAgreement` is the paper's check: over 1024 random coordinate pairs, both the
haversine and the chordal formula match the independent `arccos` reference (and
each other) to well within `1e-9` rad — the paper reports `~1e-13`.
`TestLondonNewYork` pins the reference distance (5570.22 km at R = 6371 km).

## Layout

```
geo/geodist.go        Haversine, Chordal, CentralAngle (the three methods)
geo/batch.go          precomputable per-point coordinates (the batch story)
geo/geodist_test.go   the numerical check + benchmarks
cmd/distance/main.go  CLI: one distance, computed both ways
docs/MATH.md          why Haversine == Chordal (the stereographic-chordal identity)
docs/BENCHMARK.md     batched-distance timings and what they mean
docs/DOSSIER.md       project context: the three papers and their status
CLAUDE_HANDOFF.md     orientation for a future AI/human session
```

## Notes and caveats

- **Radius / units.** Distances come back in the same unit as the `radius`
  argument. `geo.EarthRadiusKm = 6371.0` reproduces the paper's reference value;
  pass any radius you like.
- **Chordal near the south pole.** The stereographic image blows up at the south
  pole (the one point P1's flattening excludes), so `Chordal` loses precision only
  in a tiny neighbourhood of it. For general use prefer `Haversine`; `Chordal`
  exists to demonstrate the identity and to expose the precomputable coordinate.
- **Sphere model only.** This is spherical distance. Ellipsoidal models
  (Vincenty, Karney) and road-network travel time are out of scope; on those the
  spherical value is at best a fast approximation.
- **Module path / the paper's URL.** P2's remark cites
  `https://github.com/USERNAME/haversine-check`. Replace `USERNAME` there and,
  if you want `go install` from that path to work, change the `module` line in
  `go.mod` (and the import in `cmd/distance`) to the full GitHub path.

## Provenance and disclosure

Drafted with Claude (Anthropic); the mathematics and the numerical agreement were
verified (Python cross-check and the Go test here). Per the papers' policy, name a
human author and add an AI-use acknowledgement; do not list AI as an author.
See [`LICENSE`](LICENSE) (MIT, fill in the year/name).
