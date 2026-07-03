# geo distance benchmark — v2026-06-21_1305 - the initial proposal.

Compares great-circle distance implementations from the project discussion. Every
variant returns metres and uses the same Earth radius (6378100 m), so the only
thing that differs is how `a = sin²(d/2)` is computed — the comparison is apples
to apples.

## Files
- `geodist.go` — the five implementations.
- `geodist_test.go` — a correctness test plus the benchmarks.
- `go.mod` — module definition.

## The five implementations

| name | what it is | forward trig, single pair | forward trig, per cached pair |
|---|---|---|---|
| `Distance` | uploaded code, squares with `math.Pow` | 4 + two `Pow` | — |
| `DistanceFast` | same, squares with `s*s` | 4 | — |
| `FlattenDistance` | stereographic flatten + chordal metric, no cache | 2 `tan` + 2 `sincos` | — |
| `ChordDistance` (+`Flatten`) | cached flatten | (setup once per point) | **0** |
| `HavCachedDistance` (+`NewGeoPoint`) | cached law-of-cosines, kept well-conditioned | (setup once per point) | **1** (`cos Δλ`) |

The two cached forms differ per pair by exactly one cosine — that is the whole
amortization argument, made measurable.

## Run it

    cd geo_bench_2026-06-21_1305
    go mod tidy        # no external deps; just initialises the module cache
    go vet ./...       # static check

Correctness (all variants must agree with `Distance`):

    go test -run Test -v

Benchmarks:

    go test -bench=. -benchmem -run=^$

For stable, comparable numbers (recommended):

    go test -bench=. -benchmem -run=^$ -benchtime=2s -count=6 | tee bench.txt
    # optional: go install golang.org/x/perf/cmd/benchstat@latest
    benchstat bench.txt

## How to read it

Two groups of benchmarks:

`BenchmarkPair_*` — one isolated distance per call. This settles the hand-analysis
from the discussion. Expect `DistanceFast` to beat `Distance` (no `Pow`).
`FlattenDistance` vs `DistanceFast` is genuinely close and machine-dependent —
using `math.Sincos` narrows the gap that a naive `sin`+`cos` would open, which is
exactly the kind of constant-factor question only a benchmark settles.

`BenchmarkAllPairs/<method>/N=<n>` — the O(n²) all-pairs task, measured end to end
including each method's own precompute. This is where Big-O eventually bites:

- `FlattenNoCache` should be the slowest — it pays the full flatten on every pair.
- `Haversine` (recompute each pair) sits in the middle.
- `HavCached` and `FlattenCached` pull ahead as `N` grows, because the per-point
  trig is paid once (O(n)) and the O(n²) pairs are nearly trig-free.
- `FlattenCached` should edge out `HavCached` — by that one cosine per pair —
  and both stay well-conditioned for small distances (both route through
  `asin(sqrt(a))`).
- `-benchmem` shows the cached forms make one O(n) slice allocation; the naive
  forms allocate nothing.

The crossover N (where caching starts to win) is the interesting number, and it is
machine-, Go-version-, and CPU-dependent — that is the point of running it rather
than reasoning about it.

## Notes

- Numbers are relative; absolute ns/op vary by hardware and Go version. Compare
  rows within one run, not across machines.
- `ChordDistance`/`HavCached` are the well-conditioned forms (squared half-chord),
  so they keep haversine's small-angle accuracy while being batch-friendly. A plain
  `acos(u·v)` law-of-cosines would be faster still but loses precision near zero
  distance — deliberately not included as the "fast" option for that reason.
- This was verified numerically (all five agree to ~7×10⁻⁷ m over 2×10⁵ random
  pairs) but **not compiled** in the environment it was written in — run `go vet`
  and the `-run Test` correctness pass first.