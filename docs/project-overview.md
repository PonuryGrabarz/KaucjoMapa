# Project overview

Condensed, always-true summary of product scope. Full detail (feature specs,
proposed extras, user flow) lives in root `README.md` (Polish) — this file
exists so an agent doesn't have to re-read the whole thing every session.

## What it is

KaucjoMapa (working name **KaucjoMapa**, may change before launch) — a web
app for finding the nearest deposit-return point for Polish kaucja
(deposit-refund) packaging: butelkomaty (reverse vending machines), shops,
and other collection points. Shows points on a map, lets users navigate to
one, and lets the community keep point status up to date.

## Current status

Concept / MVP stage. No backend exists yet — frontend runs against a dummy
API client in dev. See `docs/architecture.md` for the technical shape.

## MVP scope (from README roadmap)

- Map view of return points
- User geolocation
- Nearest-points list
- Point detail view
- Navigation to a point
- Add a new point
- Update / confirm / reject point status
- Filter points by type and status

Post-MVP (v1.0): user accounts, submission history, reputation system, point
photos, admin/moderation panel, favorites.
v2.0 (later): notifications, point ranking, external operator integration,
PWA, data import, eco stats.

## Status enum for a return point

`Dostępny` (available) · `Zapełniony` (full) · `Nieczynny` (closed) ·
`Awaria` (fault) · `Brak informacji` (no info) · `Do weryfikacji` (pending
verification). Each status carries a last-updated date and confirm/reject
counts (community-verified, not yet implemented in frontend).

## Where to go next

- Full product spec / feature list / proposed extras → `README.md`
- Technical architecture, stack, API client → `docs/architecture.md`
- Domain vocabulary → `docs/glossary.md`
- Past/settled decisions → `docs/decisions.md`
- Frontend folder rules → `frontend/src/README.md`, `frontend/AGENTS.md`
