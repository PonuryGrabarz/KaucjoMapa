# KaucjoMapa Frontend Context

Loads automatically when working inside `frontend/`. This is the level-1
router scoped to frontend code — see root `AGENTS.md` for the project-wide
router and `docs/` for background knowledge.

## Scope

Astro 7 + TypeScript app. Renders the map of deposit-return points, status,
navigation, and (later) reports/community features. No backend yet — see
`src/shared/api`.

## Stack

- **Astro 7** (`astro.config.mjs`), TypeScript, **Sass** (`.module.sass` per
  component), **Leaflet** for the map, **ky** for HTTP.
- Package manager: **bun** (`bun.lock`) — run everything with `bun --cwd
  frontend <script>` or `cd frontend && bun run <script>`.
- Env vars: `astro:env/client` — `API_URL`, `USE_DUMMY_API` (see `.env`,
  `astro.config.mjs`).

## Where things live (routing)

| Need | Path |
| --- | --- |
| Routes (thin, compose only) | `src/pages/*.astro` |
| Return-points feature (map, list, pin, FAB) | `src/features/return_points/` |
| Shared layouts | `src/shared/layouts/` (`app_layout`, `main_layout`) |
| Shared UI (topbar, navigation) | `src/shared/components/` |
| API client + factory (prod/dummy switch) | `src/shared/api/api-client-factory.ts` |
| API request/response types | `src/shared/api/requests/`, `src/shared/api/responses/` |
| API client implementations | `src/shared/api/implementations/` (`api-client.ts` = real, `dummy-api-client.ts` = dev) |

## Conventions

- Feature-based structure: a feature owns its `components/`, `lib/`,
  `data/`, `types.ts`. Features never import from each other directly —
  promote shared code to `src/shared/`. `shared/` never imports from
  `features/`.
- **Folder naming: snake_case** for all `src/` folders (e.g.
  `features/return_points/`, `shared/layouts/app_layout/`). This is the
  canonical convention — do not switch to kebab-case.
- Path aliases (`tsconfig.json`): `@features/*`, `@shared/*`, `@assets/*`,
  `@layouts/*`, `@components/*` — use them instead of relative `../../..`.
- Component pattern in use: `<name>.astro` (markup) + `<name>.ts` (logic) +
  `<name>.module.sass` (scoped styles), one folder per component.
- API access always goes through `apiClient` from `api-client-factory.ts` —
  never call `fetch`/`ky` directly from feature code. In dev it can be backed
  by `DummyApiClient`; don't special-case that in feature code.

## Verification

Run from repo root or `frontend/`:

```bash
bun --cwd frontend check   # astro check (types)
bun --cwd frontend lint    # eslint
bun --cwd frontend format:check  # prettier
bun --cwd frontend build   # production build
```

All four should pass before a frontend change is considered done.
