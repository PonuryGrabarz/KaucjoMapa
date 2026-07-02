# Architecture

Technical shape of the project. Product scope lives in
`docs/project-overview.md` / `README.md`; this file is the "how it's built",
not the "what it does".

## Repos / packages

Single repo, single deployable today: `frontend/` (Astro app). Root
`package.json` only wires husky + commitlint for the whole repo — it is not
a workspace/monorepo build. No `backend/` directory exists yet.

## Frontend stack

- **Astro 6**, TypeScript, **Sass** (CSS modules, one `.module.sass` per
  component), **Leaflet** (map rendering), **ky** (HTTP client wrapper).
- Package manager: **bun** for `frontend/` (see `bun.lock`).
- PWA support present (`public/manifest.webmanifest`, service worker via
  `dist/sw.js` output) — see commit `855b89b` (`feat(pwa)`).

## Folder architecture

Feature-based, not layer-based: everything for one feature (components,
domain logic, types) lives together under `frontend/src/features/<name>/`.
Cross-feature code goes in `frontend/src/shared/`. `frontend/src/pages/` is
routing-only. Folder naming convention: **snake_case** throughout `src/`
(e.g. `features/return_points/`). Full rules and path aliases:
`frontend/AGENTS.md`.

## API layer

`frontend/src/shared/api/`:

- `api-client-factory.ts` — decides at import time which implementation to
  use: real `ApiClient` (prod, reads `API_URL` from `astro:env/client`) or
  `DummyApiClient` (dev, when `USE_DUMMY_API` is set). Feature code always
  imports the `apiClient` singleton from here, never instantiates a client
  itself.
- `abstractions/api-client-base.ts` — interface both implementations satisfy.
- `implementations/` — `api-client.ts` (real, ky-based presumably),
  `dummy-api-client.ts` (static/mock data for local dev without a backend).
- `requests/`, `responses/` — typed request/response shapes per endpoint,
  e.g. `get-map-points-request.ts` / `get-map-points-response.ts`.

This split exists specifically because there is no backend yet — it lets
frontend development proceed against dummy data and swap in the real API
later without touching feature code.

## Current features (frontend)

- `return_points` — map view (`return-points-map`), home page view, map pin
  rendering, floating action button. This is the only implemented feature so
  far; README's other domain areas (status, reports, navigation) are planned
  but not yet built as separate features.

## Known drift / things to watch

- No backend, no tests, no CI beyond commit/branch/PR-name linting yet
  (`.github/workflows/ci.yml`, `branch-name.yml` per `CONTRIBUTING.md`).

## Where to go next

- Product scope → `docs/project-overview.md`
- Frontend folder rules in full → `frontend/AGENTS.md`
- Why a choice was made → `docs/decisions.md`
