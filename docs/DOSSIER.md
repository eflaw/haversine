# Project context (dossier excerpt)

The still-relevant parts of the project dossier, trimmed to what a reader of *this
repo* needs. The full dossier lives with the papers.

## Through-line

One object — the unit circle, written as `e^{iλ}` or read as two perpendicular
projections — generates spherical geometry two ways and rests on a single rigorous
root (the complex exponential). Everything is **classical**; the value is unified
**exposition**, not new theorems. Genre: expository/pedagogical mathematics
(*Monthly*, *Mathematics Magazine*, *College Mathematics Journal*; arXiv `math.HO`).

## The three papers

- **P1 — *The Sphere on the Complex Disk.*** `w = R(φ)e^{iλ}` is a homeomorphism
  `S²∖{S pole} → 𝔻`; the antipode is the one-point compactification; de Moivre makes
  axial rotation `= ×e^{iα}`. Verified, standalone.
- **P2 — *From the Unit Circle to the Haversine.*** sin/cos as projections →
  haversine = squared half-chord → (3D) great-circle distance formula. Verified,
  standalone. **This repo is the code it cites.**
- **P3 — *Six Roads to Sine and Cosine.*** Six equivalent definitions of sin/cos.
  In progress (needs writing + gap-patching).

Reading order: P3 → P2 → P1. Submit-first order: P2, then P1, then P3.

## What changed between the draft and the verified papers

Relevant because this repo carries the difference:

- An earlier P1 draft had a **bridge** section (haversine = stereographic chordal
  metric) and a **computational appendix** (the benchmark). Both were **cut** from
  the verified papers to keep them lean and expository. That correct-but-cut material
  is now `docs/MATH.md` and `docs/BENCHMARK.md` here.
- P2 gained a one-line pointer to a companion code repo
  (`github.com/USERNAME/haversine-check`) — this repo.

## Standing caveats (still apply)

- **References.** The verified papers already fixed the previously shaky citations
  (de Moivre now attributed to the 1722 *Phil. Trans.* "De sectione anguli", with the
  general form credited to Euler 1749; Inman 1835 given as the 3rd edition that coined
  "haversine"). If you touch citations, keep to primary/textbook sources.
- **AI-use disclosure.** Single human author; name Claude (Anthropic) as a
  drafting/typesetting tool in an acknowledgements line; author affirms verifying all
  mathematics and references. Do **not** list AI as an author.
- **Framing.** Introduce everything as exposition/synthesis, not priority or novelty —
  including the identity in `docs/MATH.md`, which is a corollary of classical facts.
- **Author line.** "Independent researcher" is acceptable for these venues.

## Not carried over

The dossier's figure inventory, the "square construction" note, and the detailed
publication timeline are about the papers, not the code, and are omitted here. See
the full dossier for those.
