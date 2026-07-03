# CLAUDE_HANDOFF.md

Orientation for a future session (AI or human) picking up this repo. Read this
before changing anything.

## 1. What this repo is

Companion code for a trilogy of short expository maths papers, *Circle &amp; Sphere*.
It exists to make one claim reproducible: the **haversine formula** (paper P2) and
the **stereographic chordal metric** (paper P1) compute the **same** great-circle
distance. The repo verifies that against an independent reference and provides the
worked "both methods" example.

Paper P2's numerical-check remark points at this repo by name
(`github.com/USERNAME/haversine-check`), so it is not optional scaffolding — it is
the code the paper cites.

## 2. State of the papers (as of 2026-07-03)

Three papers, all framed as **exposition/synthesis of classical results** — no new
theorems is claimed, and that framing must be preserved.

| # | Title | Content | Status |
|---|---|---|---|
| P1 | *The Sphere on the Complex Disk* | flattening map `w = R(φ)e^{iλ}`; homeomorphism `S²∖{S}→𝔻`; antipode = one-point compactification; de Moivre → axial rotation | verified, standalone |
| P2 | *From the Unit Circle to the Haversine* | projections → haversine = squared half-chord → (3D) great-circle distance | verified, standalone; **cites this repo** |
| P3 | *Six Roads to Sine and Cosine* | six equivalent definitions of sin/cos | in progress |

The verified `.tex` sources the user holds are the **source of truth**. Filenames
seen this session: `circle_sphere_2026-07-03_1010.tex` (P1),
`haversine_2026-07-03_1010.tex` (P2).

## 3. IMPORTANT — do not "restore" removed material into the papers

Earlier drafts (unverified) contained extra material that this repo now carries:

- a **bridge** section proving haversine = chordal metric (the identity in
  `docs/MATH.md`), plus an **equidistant-radius** proposition;
- a **computational appendix** with the benchmark (now `docs/BENCHMARK.md`).

The verified papers **deliberately removed** all of it to stay lean and purely
expository. That was an editorial choice, not an error. **Do not reinsert it into
the papers.** Its correct home is this repository. If asked to touch the papers,
confirm scope first; default to leaving the verified `.tex` untouched.

## 4. Conventions the code and math depend on (must stay consistent with P1)

- Flattening map: `F(λ, φ) = R(φ)·e^{iλ}`, **north pole → 0**.
- Stereographic radial law: `R(φ) = tan(π/4 − φ/2)`; equator → `|w|=1`,
  **south pole → ∞**. (P1 restricts this law to the northern hemisphere only so it
  lands *in the unit disk*; the distance formula below uses the **unrestricted**
  coordinate, which fills the plane. Both are consistent — the restriction is about
  "landing in the disk", not about the distance identity.)
- Chordal identity (the whole point): with `w_j = R(φ_j)e^{iλ_j}`,
  `hav d = sin²(d/2) = |w₁−w₂|² / ((1+|w₁|²)(1+|w₂|²))`. Derivation in
  `docs/MATH.md`; key step `1+|w|² = 2/(1+sin φ)`.
- Reference radius: **R = 6371 km** (matches P2's London–NY ≈ 5570 km). Do not
  silently switch to 6378.137 (equatorial) or 6371.0088 — it changes the quoted
  number. `geo.EarthRadiusKm = 6371.0`.
- Distance form: `d = 2·atan2(√a, √(1−a))` (robust at antipodes), not `2·asin√a`.

## 5. Verification status

- Math cross-checked in Python and by the Go test `TestAgreement`: haversine and
  chordal each match `arccos(u₁·u₂)` to ~`1e-13` rad over `1e5` random pairs.
- `TestLondonNewYork` pins 5570.22 km at R = 6371 km.
- The Go was **written but not compiled in the environment it was authored in**
  (no Go toolchain there). First action for a new session with Go available:
  `go vet ./... && go test ./...`. It should pass; if not, the likely spots are
  trivial (imports, a Unicode identifier) — the math is confirmed independently.

## 6. Open TODOs

- [ ] Replace `USERNAME` in P2's repo URL, and (if publishing under a module path)
      in `go.mod` + `cmd/distance/main.go`.
- [ ] Fill in author name + year in `LICENSE`, and the `[AUTHOR NAME]` in the papers.
- [ ] Add the AI-use acknowledgement line to each paper (human author; Claude named
      as a drafting/typesetting tool; author affirms verifying all maths + refs).
- [ ] Optional: harden the benchmark (record Go version, run `-count` + `benchstat`,
      a second CPU) before quoting numbers anywhere load-bearing.
- [ ] Optional: if anyone wants the chordal method robust at the south pole, project
      from whichever pole is farther (adaptive centre); not needed for the identity.

## 7. Map of the repo

- `geo/geodist.go` — `Haversine`, `Chordal`, `CentralAngle` (+ internal helpers).
- `geo/batch.go` — `FlatPoint`/`Flatten`/`ChordDistance`, `GeoPoint`/`NewGeoPoint`/
  `HavCached`: the precomputable coordinates behind the batch/benchmark story.
- `geo/geodist_test.go` — the numerical check and the benchmarks.
- `cmd/distance/main.go` — the CLI worked example (both methods side by side).
- `docs/MATH.md` — the stereographic-chordal identity, derived gap-free.
- `docs/BENCHMARK.md` — batched-distance timings + honest limits.
- `docs/DOSSIER.md` — project context and standing caveats carried over from the
  original dossier (only the still-relevant parts).

## 8. Provenance

Repo assembled with Claude (Anthropic, Opus-class) on 2026-07-03, salvaging
correct-but-cut material from earlier paper drafts after checking it against the
verified 2026-07-03 papers. Everything here is a corollary of classical facts;
present it as synthesis, never as a new result.
