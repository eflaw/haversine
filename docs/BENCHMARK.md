# Batched distance: benchmark notes

This documents the batch story behind `geo/batch.go`. It was an appendix in an
earlier P1 draft and was **cut from the verified paper**; it is illustrative
engineering, not a paper claim.

## The idea

Because a point's stereographic image depends on that point alone, it is a
**precomputable coordinate** `c = (Re w, Im w, |w|²)`. For an all-pairs distance
over `N` points you flatten each point once (`O(N)` trig) and then every pair is
`O(1)` **arithmetic** (a subtraction, a few multiplies, a divide, `sqrt`, `atan2`)
with **no forward trigonometry**. The haversine cannot amortise as cleanly: it uses
sines of coordinate *differences*, so reducing it to per-point data forces the
spherical-law-of-cosines form, which still costs one `cos(Δλ)` per pair.

`go test -bench=. -benchmem -run=^$ ./...` benchmarks four ways of filling an
all-pairs matrix:

| benchmark | forward trig per pair |
|---|---|
| `ChordalCached`  | 0 |
| `HavCached`      | 1 (`cos Δλ`) |
| `Haversine`      | 4 (recomputed each pair) |
| `ChordalNoCache` | 4 (flatten recomputed each pair) |

## Representative numbers (illustrative)

From one run on an Intel i5-4460 (earlier prototype; **your numbers will differ** —
this is a constant-factor, machine-dependent effect). Per-pair cost, normalised by
`N(N−1)/2`, was roughly:

| method | ns/pair |
|---|---|
| ChordalCached  | ~46 |
| HavCached      | ~66 |
| Haversine      | ~120 |
| ChordalNoCache | ~150 |

Reading: there is a shared "tail" (`sqrt` + `atan2`/`asin` + arithmetic, ~46 ns);
each forward trig call adds ~20 ns; the cached chordal coordinate removes them all,
the cached haversine keeps exactly one cosine (~22 ns), and the recomputed forms pay
four. At `N = 1024` the cached coordinate was ~2.6× faster than recomputed haversine.

## Honest limits (why this is an appendix, not a paper)

- **Constant factor, not asymptotic.** Every method is `Θ(1)` per pair and `Θ(N²)`
  all-pairs; caching moves the constant, not the complexity class.
- **Conditional on reuse.** For a *single* distance the coordinate is built and never
  reused, so plain `Haversine` wins. The benefit needs many pairs from a fixed set
  (distance matrix, nearest-point, clustering).
- **Bigger levers exist.** In a real system a spatial index (k-d tree, R-tree,
  S2/H3) avoids most pairs entirely, and SIMD/GPU vectorises the rest. The
  precompute-then-arithmetic shape here *composes* with both but does not replace
  them.
- **Conditioning.** `ChordDistance` and `HavCached` route through
  `a = sin²(d/2)` and `atan2`, so they keep the haversine's small-angle accuracy — a
  bare `arccos(u·v)` would be faster still but lose precision for nearby points.

For load-bearing numbers, re-run with `-count=6 -benchtime=2s` and reduce with
`benchstat` (`go install golang.org/x/perf/cmd/benchstat@latest`), and record the
Go version and CPU alongside the output.
