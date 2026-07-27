# KaucjoMapa — Project Router (AGENTS.md)

This file loads automatically whenever an agent (Hermes, Claude Code, Codex...)
opens this repo. It is the **entry point of the second brain** — level 1:
routing by exact folder/file, no semantic search, no knowledge graph. Keep it
that way until routing itself becomes a pain point.

## Purpose

KaucjoMapa is a web app that helps users find the
nearest deposit-return point (butelkomat / shop / other) for Polish kaucja
(deposit) packaging, check its live status, and crowdsource updates. Full
product spec lives in `README.md` (Polish) — read it once for domain
vocabulary, then rely on `docs/` for anything that changes.

Stack today: **Astro 7 + TypeScript + Sass + Leaflet**, `frontend/` only.
No backend yet — `frontend/src/shared/api` talks to either a real API
(`API_URL` env) or a `DummyApiClient` in dev (`USE_DUMMY_API` env). Package
manager: **bun** for `frontend/`, **npm** at repo root (husky/commitlint only).

## Repository map

| Path | What lives there | When to look here |
| --- | --- | --- |
| `README.md` | Product spec, feature list, status enum, roadmap (Polish) | Understanding *what* to build / product questions |
| `CONTRIBUTING.md` | Branch naming, commit convention, PR rules, husky/commitlint | Before any commit/branch/PR |
| `AGENTS.md` (this file) | Router — routing rules + working rules | Every session start |
| `frontend/AGENTS.md` | Frontend-scoped context: folder architecture, path aliases, naming convention | Any frontend code change / deciding where new code goes |
| `docs/project-overview.md` | Condensed product/domain summary + current scope | Fast context without re-reading full README |
| `docs/architecture.md` | Technical architecture: stack, data flow, API client pattern | Understanding how pieces fit together |
| `docs/decisions.md` | Decision log (append-only) | Before re-litigating a past decision; after making a new one |
| `docs/glossary.md` | Domain terms (kaucjomat, punkt zwrotu, status types...) | Unclear domain wording |
| `frontend/src/features/` | Feature code (return_points, later: reports, navigation...) | Implementing/reading a feature |
| `frontend/src/shared/` | Cross-feature code (api, layouts, components) | Reusable pieces, API client |
| `frontend/src/pages/` | Astro routing only — thin, composes features | Adding/inspecting a route |

## Routing rules

- Product/business question ("what should this feature do?") → `README.md`
  then `docs/project-overview.md`.
- Domain word you don't recognize (kaucjomat, status "Do weryfikacji"...) →
  `docs/glossary.md`.
- "How is the frontend structured / where does X go?" → `frontend/AGENTS.md`
  (feature-based layout, snake_case naming, path aliases).
- "Why was it built this way?" / architecture question → `docs/architecture.md`;
  if the *reasoning* behind a specific choice matters → `docs/decisions.md`.
- Git workflow, commit message format, branch names → `CONTRIBUTING.md`.
- Anything specifically inside `frontend/` → load `frontend/AGENTS.md`, it is
  the frontend-level router and takes precedence for frontend implementation
  detail.
- No backend exists yet. Do not invent backend routing — if asked about
  backend, say it doesn't exist yet and check `docs/decisions.md` for any
  plan before assuming.

## Second brain files

Level: **1** (exact-word routing, markdown only, no vector/graph search).
Escalate to level 2 (an LLM wiki under `docs/wiki/`) only once `docs/*.md`
grows past ~10 files and cross-references get hard to keep straight. Don't
build that ahead of the pain.

- `docs/project-overview.md` — always-true background on product scope.
- `docs/architecture.md` — always-true background on technical shape.
- `docs/decisions.md` — append-only decision log, newest entry on top.
- `docs/glossary.md` — domain vocabulary, grows as needed.

## Working rules

- Before creating or modifying any `.md` file in this repo (this file,
  `frontend/AGENTS.md`, anything under `docs/`, `README.md`,
  `CONTRIBUTING.md`) as a result of a decision made during a conversation,
  **always ask the user first** whether to make that edit — do not write
  automatically, even for routine decision-log entries.
- Reuse existing conventions before inventing new ones: feature-based folders
  with **snake_case** naming (`frontend/AGENTS.md`), path aliases
  (`@features/*`, `@shared/*`...), Conventional Commits (`CONTRIBUTING.md`).
- Before adding a new top-level concept to `docs/`, check `docs/decisions.md`
  — it may already be settled or explicitly rejected.
- Keep `README.md` as the single source of truth for product scope; `docs/`
  summarizes and links back to it, it does not fork it.

## Verification

- Frontend changes: `bun --cwd frontend check` (astro check), `bun --cwd
  frontend lint`, `bun --cwd frontend build` before calling something done.
- Commit messages/PR titles must pass commitlint (Conventional Commits) —
  see `CONTRIBUTING.md`. `npm install` at repo root wires up husky hooks.
- When docs are updated, keep this file (`AGENTS.md`) in sync if routing
  changes — stale routing is the main level-1 failure mode.

## Memory policy for Hermes

- Durable facts (stack choices, conventions, folder architecture) belong in
  `docs/architecture.md` / `frontend/AGENTS.md`, not in Hermes cross-session
  memory — this repo is the second brain, keep it in-repo and git-tracked.
- Use Hermes memory only for facts truly external to the repo (e.g. user's
  own workflow preferences), never for project state that belongs here.
