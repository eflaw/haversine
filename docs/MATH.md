# Why `Haversine == Chordal`

This is the identity the repo demonstrates. It was a section in an earlier draft of
P1 ("The two descents meet") but was **cut from the verified paper** to keep it
lean. The mathematics is unchanged and consistent with the verified P1's
conventions; it lives here as the reason the two code paths agree.

## Setup

A surface point at latitude `φ`, longitude `λ` has the stereographic disk image
(P1's flattening map with the stereographic radial law)

```
w = R(φ) · e^{iλ},    R(φ) = tan(π/4 − φ/2).
```

North pole (`φ = π/2`) → `0`; equator → `|w| = 1`; south pole (`φ = −π/2`) → `∞`.

The **chordal metric** of the Riemann sphere measures the central angle `d` between
two points from their images by

```
sin²(d/2) = |w₁ − w₂|² / ((1 + |w₁|²)(1 + |w₂|²)).
```

Claim: the right-hand side equals the haversine `a = hav d` of P2. Hence
`Chordal(...) == Haversine(...)`.

## Two facts that make it collapse

**(1) The denominator.** With `t = π/4 − φ/2`, so `R = tan t` and `2t = π/2 − φ`:

```
1 + |w|² = 1 + tan²t = sec²t,     cos²t = ½(1 + cos 2t) = ½(1 + sin φ),
```

therefore

```
1 + |w|² = 2 / (1 + sin φ).
```

**(2) The numerator.** With `w_j = R_j e^{iλ_j}` and `Δλ = λ₁ − λ₂`,

```
|w₁ − w₂|² = R₁² + R₂² − 2 R₁ R₂ cos Δλ.
```

## The collapse

Multiply the ratio through by `cos²t₁ cos²t₂` and use `tan t · cos t = sin t`:

```
|w₁−w₂|² / ((1+|w₁|²)(1+|w₂|²))
   = sin²t₁ cos²t₂ + cos²t₁ sin²t₂ − 2 sin t₁ cos t₁ sin t₂ cos t₂ cos Δλ.
```

Now substitute `sin²t = ½(1 − sin φ)`, `cos²t = ½(1 + sin φ)`, and
`2 sin t cos t = sin 2t = cos φ`. The first two terms combine to

```
¼[(1 − sin φ₁)(1 + sin φ₂) + (1 + sin φ₁)(1 − sin φ₂)] = ½(1 − sin φ₁ sin φ₂),
```

and the third is `½ cos φ₁ cos φ₂ cos Δλ`. So the ratio equals

```
½[ 1 − (sin φ₁ sin φ₂ + cos φ₁ cos φ₂ cos Δλ) ] = ½(1 − cos d) = hav d,
```

the bracket being exactly the **spherical law of cosines**
`cos d = sin φ₁ sin φ₂ + cos φ₁ cos φ₂ cos Δλ` — which is P2's own eq. (slc).
Therefore the chordal ratio **is** the haversine `a`, and the two code paths return
the same distance. ∎

## Convention note (the one a referee would probe)

The `1 + |w|²` normalisation is tied to this projection direction — north pole at
the centre, south pole at infinity, i.e. stereographic projection *from the south
pole*. If you re-derive from a projection with the opposite centre, the algebra is
the mirror image and the same identity results; but state the convention, because
the normalisation depends on it.

## A second, weaker fact (also cut): the equidistant radius

For the **equidistant** radial law `R(φ) = ½ − φ/π`, the disk radius is literally
the great-circle distance from the pole: writing the colatitude `c = π/2 − φ`,

```
|w| = R(φ) = c/π,     so   distance-from-north-pole = π·R_sphere·|w|.
```

So the rim `|w| = 1` is the south pole at `π·R_sphere` (half a great circle), and the
full diameter corresponds to `2π·R_sphere` (a complete great circle). This is only a
distance *from the centre* — for two arbitrary points you still need a spherical
triangle — so it is a curiosity rather than a computational tool, which is why the
paper dropped it. It is recorded here for completeness.

## Numerical confirmation

At `R = 6371 km`, London → New York gives `5570.2222 km` by all three methods
(haversine, chordal, and the `arccos` reference), and haversine vs chordal agree to
`~1e-13 rad` over `1e5` random pairs. Run `go test ./...` (see `geo/geodist_test.go`)
or `go run ./cmd/distance`.
