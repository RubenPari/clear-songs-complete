---
name: invasive-cleanup-monorepo
description: Plan and execute invasive cross-stack cleanup in a monorepo while preserving behavior and verifying through tests.
source: auto-skill
extracted_at: '2026-06-17T09:51:52.274Z'
---

# Invasive cross-stack cleanup with test verification

Use this skill when a user asks for aggressive refactoring and cleanup across multiple packages in a monorepo (e.g. Go backend + Angular frontend) and explicitly wants test/build verification.

## Procedure

### 1. Set expectations and confirm scope
- Ask the user to clarify three things:
  1. **Depth**: surface-only or invasive?
  2. **Verification**: run tests/builds/lint?
  3. **Exclusions**: any areas to leave untouched?
- If the user answers in a compact form (e.g. "1 - invasiva, 2 - verifica con i test, 3 - no"), interpret:
  - Position 1 = invasive cleanup.
  - Position 2 = verify via tests/builds.
  - Position 3 = no exclusions.

### 2. Explore the whole repo before touching anything
- List root and key subdirectories.
- Read `README.md`, `AGENTS.md`, `QWEN.md`, and package manifests (`go.mod`, `package.json`).
- Identify entrypoints, DI container, routing files, generated files, env files, and test setup.
- Look for duplicated logic, stale comments, dead functions, inconsistent naming, and generated-but-committed artifacts.

### 3. Plan cleanup by stack

**Backend (Go)**
- Find and remove unused exports/functions using static search / grep first. Confirm no callers across the module.
- Extract duplicated logic into domain helpers (e.g. redirect URL resolution, env loading) so callers share one source of truth.
- Standardize file naming conventions (e.g. snake_case for multi-word files if that is the repo convention).
- Remove redundant doc comments that restate obvious implementation details.
- Update stale README paths when the code has moved.

**Frontend (Angular / TS)**
- Extract duplicated component logic into injectable services or pure helpers.
  - Example: theme/lang init duplicated between `App` and `MainLayoutComponent` → `PreferencesService` with signals.
  - Example: repeated confirm → loading → toast flow → `confirmAndRunWithNotify` helper.
- Remove dead helper functions/operators that were kept "just in case".
- Delete unused imports after refactoring; lint will catch this, but fix proactively.
- Avoid deleting DOM/class-toggling code unless the service already handles it; `Renderer2` effects may still be needed to sync signal state to the DOM.

### 4. Make changes incrementally and verify in batches
- After each logical change group, run the relevant verification command:
  - `cd backend && go test $(go list ./... | grep -v '/test/e2e')` (or equivalent e2e exclusion).
  - `cd frontend && npm run lint`.
  - `cd frontend && npm run build` as a proxy for typechecking when no dedicated `typecheck` script exists.
- Fix compiler/lint errors immediately before moving to the next change group.

### 5. Handle common verification pitfalls
- If `npm test` requires a browser and the user says not to launch Chrome/Chromium, treat the build (`npm run build`) as the typecheck gate and explicitly note the skipped test run.
- If generated files (e.g. `environment.auto.ts`) are produced by `prestart`/`prebuild`, let them regenerate rather than editing by hand.
- Watch out for imports that lint did not catch but the Angular compiler did (e.g. `PLATFORM_ID`, `DOCUMENT`) after moving code into a new file.

### 6. Final state and handoff
- Keep the todo list updated with the verification steps and which ones passed/skipped.
- Summarize the cleanup per stack and state what tests/builds passed.
- Ask the user whether to review diffs, commit, or continue with deeper cleanup.

## Structural cleanup extension

If the user asks to improve the project folder structure after a cleanup pass, extend the procedure with an architectural review before moving files:

### 1. Diagnose structural friction
- Inspect the top-level layout of both stacks (Go `internal/` and Angular `src/app/`).
- Identify packages with generic names (`shared`, `utils`, `helpers`, `common`) that act as dumping grounds.
- Check whether a layer/feature model is applied consistently across all domains (e.g. `auth` should have the same depth as `track`/`playlist`).
- Look for layer violations: controllers importing external SDK types to do mapping, env/config living in `domain`, transport leaking into application logic.

### 2. Choose one structural improvement at a time
- **Relocate configuration/env helpers** from `domain` to `infrastructure/config` and update all callers (`cmd/server/main.go`, DI container).
- **Relocate Spotify-specific helpers** (image selection, ID conversion) to `infrastructure/external/spotify/helpers` so the domain stays pure.
- **Introduce missing domain packages** (e.g. `domain/auth`) and move lightweight entities/mapping there; make application use cases depend on the domain model rather than SDK types.
- **Move SDK-to-DTO mapping out of HTTP controllers** into the application layer, so controllers only handle input binding, error translation, and status codes.

### 3. Verify structural changes
- After every file move or package rename, run `go test $(go list ./... | grep -v '/test/e2e')` and `go vet ./...`.
- Run `go build ./cmd/server/main.go` (or the repo's entrypoint) to catch import/type regressions that tests might miss.
- Run frontend `npm run lint && npm run build` when the backend changes could affect shared contracts.
- Do not commit generated artifacts (e.g. `environment.auto.ts`, `package-lock.json`) if they were produced only by the build.

## Why this works
- Confirming scope in three axes prevents over-cleaning or missing verification the user explicitly requested.
- Exploring before editing avoids deleting functions that only appear unused because a generated caller exists in a non-obvious file.
- Verifying in small batches catches import regressions and dead-symbol removals early instead of at the end.
- Treating structural cleanup as a follow-on, user-driven step avoids unnecessary churn while leaving a clear path to a cleaner architecture.

## How to apply
- Use this skill when a user asks to "analyze for cleanup opportunities", "refactor", "simplify", "remove duplication", or "improve folder structure" across multiple packages in a monorepo.
- Adjust strictness of "invasive" to the repo's tolerance: preserve public contracts (route paths, env variable names, observable behavior) unless the user requests breaking changes.
