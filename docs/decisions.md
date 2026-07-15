## 2026-07-02 — Folder naming: snake_case, not kebab-case

**Context:** An early draft of `frontend/src/README.md` documented a target
convention of kebab-case feature directories (`features/return-points/`)
while existing code already used snake_case (`features/return_points/`).
That file has been removed.
**Decision:** Folder naming under `frontend/src/` is **snake_case**
throughout — matches existing code, no migration needed.
**Why:** Avoids an unnecessary rename of working code to match a doc that
was never implemented; snake_case was already the de facto convention.
**Alternatives considered:** Migrate existing code to kebab-case — rejected,
no benefit, extra churn.
