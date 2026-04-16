# AGENTS.md

## Scope and layout
- Monorepo with 2 real projects:
  - `clear-songs/`: Go API (Gin) with Redis-required session/cache, optional Postgres backup, optional Gemini fallback.
  - `clear-songs-front/`: Angular 20 SPA.
- Main backend entrypoint is `clear-songs/cmd/server/main.go` (not `src/main.go` from older docs).
- API routes are wired in `clear-songs/internal/infrastructure/transport/http/routes.go`.

## Environment and config gotchas
- Copy root env first: `cp .env.example .env` (used by root `docker-compose.yml`).
- Backend accepts redirect var as either `REDIRECT_URL` or `REDIRECT_URI`; compose/readmes use both names in different places.
- Redis is effectively required for backend startup (`di.NewContainer()` fails if Redis is unavailable).
- Frontend `npm start` / `npm run build` runs `prestart`/`prebuild` and generates `src/environments/environment.auto.ts` from `.env` via `tools/generate-env.js`.
- `environment.auto.ts` is generated and gitignored; do not hand-edit it.

## Canonical dev commands
- Full stack (preferred):
  - `docker compose up --build` (repo root)
- Backend local:
  - `cd clear-songs && go mod download && go run ./cmd/server/main.go`
- Frontend local:
  - `cd clear-songs-front && npm install && npm start`

## Verification commands (run from changed package)
- Backend tests (skip e2e):
  - `cd clear-songs && go test $(go list ./... | grep -v '/test/e2e')`
- Frontend lint:
  - `cd clear-songs-front && npm run lint`
- Frontend tests:
  - `cd clear-songs-front && npm test`
- Frontend type safety check (no dedicated `typecheck` script):
  - `cd clear-songs-front && npm run build`

## Testing/behavior notes
- Backend e2e tests live in `clear-songs/test/e2e` and are intentionally excluded in normal local verification.
- Frontend dev server proxies `/auth/*`, `/track/*`, `/playlist/*` to `http://127.0.0.1:3000` (`clear-songs-front/proxy.conf.json`).

## Workflow conventions already enforced
- Frontend pre-commit hook runs `lint-staged` (`clear-songs-front/.husky/pre-commit`).
- `lint-staged` runs `eslint --fix` + `prettier --write` on staged TS/JS/JSON and prettier on HTML/SCSS.