# QWEN.md — Clear Songs (complete)

This file captures the essential context for working in this repository. Read it first when returning to this project.

## Project overview

Clear Songs is a web application for inspecting and managing a user's **saved Spotify library**.

It provides:

- A dashboard with per-artist summaries (including genre filters).
- Bulk deletion of tracks by artist or by track-count ranges.
- Playlist operations (clear playlist tracks, delete from both playlist and library).
- Optional AI-driven genre fallback via Google Gemini for artists whose Spotify tags do not map to the app's genre model.

The repository is a **monorepo** with two packages:

| Package | Technology | Purpose |
|---------|------------|---------|
| `clear-songs/` | Go 1.25 + Gin | HTTP API, Spotify OAuth, business logic, persistence |
| `clear-songs-front/` | Angular 20 + Tailwind CSS | Single-page application UI |

Supporting infrastructure (declared at the repository root):

- **PostgreSQL** — optional backup/persistence.
- **Redis** — effectively required for session storage and caching.
- **Docker Compose** — recommended way to run the full stack.

## Repository layout

```
clear-songs-complete/
├── .env.example                 # Root environment template (copy to .env)
├── docker-compose.yml           # Full-stack orchestration
├── clear-songs/                 # Go backend
│   ├── cmd/server/main.go       # Server entrypoint
│   ├── internal/
│   │   ├── application/         # Use cases (auth, playlist, track, shared)
│   │   ├── domain/              # Domain models
│   │   └── infrastructure/      # DI, external APIs, logging, persistence, HTTP transport
│   ├── test/
│   │   ├── e2e/                 # End-to-end tests (excluded from normal CI)
│   │   └── mocks/               # Test mocks
│   ├── go.mod
│   └── Dockerfile
└── clear-songs-front/           # Angular frontend
    ├── src/
    │   ├── app/
    │   │   ├── core/            # Guards, interceptors, models, services, stores, utils
    │   │   ├── features/        # auth, dashboard, playlists, tracks
    │   │   ├── layouts/         # Main layout shell
    │   │   ├── app.config.ts
    │   │   ├── app.routes.ts
    │   │   └── ...
    │   ├── environments/        # environment.ts + generated environment.auto.ts
    │   └── main.ts              # Frontend entrypoint
    ├── tools/generate-env.js    # Builds environment.auto.ts from .env
    ├── proxy.conf.json          # Dev-server proxy to backend
    └── Dockerfile
```

## Technologies and dependencies

### Backend (`clear-songs/`)

- **Go 1.25.0**
- **Gin** — HTTP web framework.
- **GORM** + `gorm.io/driver/postgres` — ORM and PostgreSQL driver.
- **go-redis/v9** — Redis client.
- **zmb3/spotify** — Spotify Web API client.
- **golang.org/x/oauth2** — Spotify OAuth2 flow.
- **google/generative-ai-go** + `google.golang.org/api` — Optional Gemini integration.
- **testify** — Testing utilities.
- **Zap** — Structured logging.

### Frontend (`clear-songs-front/`)

- **Angular 20** with the new application builder (`@angular/build`).
- **TypeScript ~5.8.2**
- **Tailwind CSS 4** + `@tailwindcss/postcss`
- **Angular CDK** and **Spartan NG brain** (`@spartan-ng/brain`) — headless UI primitives.
- **ng-icons/lucide** — Icon library.
- **ngx-translate** — Internationalization.
- **D3** — Data visualization (used in dashboards).
- **ngx-toastr** — Toast notifications.
- **Karma + Jasmine** — Unit testing.
- **ESLint + Prettier + lint-staged + Husky** — Linting and formatting.

## Environment configuration

1. Copy the example environment file at the repository root:

   ```bash
   cp .env.example .env
   ```

2. Fill in the real values in `.env`:

   - `CLIENT_ID` and `CLIENT_SECRET` from the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard).
   - `REDIRECT_URI` — backend OAuth callback (e.g. `http://127.0.0.1:3000/auth/callback`).
   - `SPOTIFY_REDIRECT_URI` — frontend callback (e.g. `http://127.0.0.1:4200/callback`).
   - `API_URL` — URL the frontend should call at build time.
   - `FRONTEND_URL` — backend CORS allowed origin.
   - Optional `GEMINI_API_KEY` and related Gemini tuning variables.

   Both redirect URIs must match **exactly** what is configured in the Spotify app settings.

3. The frontend `.env` generation:

   - `npm start` and `npm run build` run `prestart`/`prebuild`, which invokes `tools/generate-env.js`.
   - This script reads `.env` and writes `src/environments/environment.auto.ts`.
   - `environment.auto.ts` is **gitignored** and must **not** be hand-edited.

## Building and running

### Full stack via Docker Compose (recommended)

From the repository root:

```bash
docker compose up --build
```

Typical exposed services:

- API: `http://127.0.0.1:3000`
- Frontend: `http://127.0.0.1:4200`
- PostgreSQL: `127.0.0.1:5432`
- Redis: `127.0.0.1:6379`

Stop:

```bash
docker compose down
```

### Backend local development

```bash
cd clear-songs
go mod download
go run ./cmd/server/main.go
```

Requires local PostgreSQL and Redis (or tunneled equivalents) and the environment variables from `.env`.

### Frontend local development

```bash
cd clear-songs-front
npm install
npm start
```

The dev server runs on `http://127.0.0.1:4200` and proxies `/auth/*`, `/track/*`, and `/playlist/*` to `http://127.0.0.1:3000` per `proxy.conf.json`.

## Testing and verification

Run these from the package that changed:

### Backend

```bash
# Unit/integration tests (skip e2e)
cd clear-songs && go test $(go list ./... | grep -v '/test/e2e')
```

End-to-end tests in `clear-songs/test/e2e` require real services and Spotify auth; they are intentionally excluded from normal local verification.

### Frontend

```bash
# Unit tests
cd clear-songs-front && npm test

# Lint
cd clear-songs-front && npm run lint

# Type-safety check (there is no dedicated typecheck script)
cd clear-songs-front && npm run build
```

## Development conventions

### Backend

- **Entrypoint:** `clear-songs/cmd/server/main.go`. Older documentation may reference `src/main.go`; that path is obsolete.
- **Architecture:** Organized into `internal/application`, `internal/domain`, and `internal/infrastructure` (clean/hexagonal style).
- **Routes:** Wired in `clear-songs/internal/infrastructure/transport/http/routes.go`.
- **DI:** Dependencies are resolved in `clear-songs/internal/infrastructure/di` and injected into controllers/use cases.
- **Logging:** Structured with Zap; request IDs are attached via middleware.
- **Redis is effectively required:** `di.NewContainer()` fails at startup if Redis is unavailable.
- **Redirect variable names:** The backend accepts either `REDIRECT_URL` or `REDIRECT_URI`; root `docker-compose.yml` uses `REDIRECT_URI`.
- **HTTP timeouts:** Configurable via `HTTP_WRITE_TIMEOUT_SEC` and `HTTP_READ_TIMEOUT_SEC` (default 360s, useful for long `/track/summary` calls with many AI genre fallbacks).

### Frontend

- **Standalone components:** Angular application uses standalone components and signals-based reactivity where appropriate.
- **Routes:** Defined in `src/app/app.routes.ts`. Main routes: `/login`, `/callback`, `/dashboard`, `/tracks`, `/playlists`.
- **Auth guard:** `authGuard` protects the main layout routes.
- **Styling:** Tailwind CSS 4 with Spartan NG headless components.
- **i18n:** Uses `ngx-translate`.
- **Icons:** Lucide icons via `@ng-icons/lucide`.
- **Pre-commit:** `clear-songs-front/.husky/pre-commit` runs `lint-staged`, which runs `eslint --fix` + `prettier --write` on staged TS/JS/JSON and `prettier --write` on HTML/SCSS.
- **Do not edit `environment.auto.ts`:** It is generated and ignored by Git.

## Common gotchas

- If the backend container exits immediately, check that Redis is reachable; Redis is required for startup even if Postgres is optional for some flows.
- `REDIRECT_URI` and `SPOTIFY_REDIRECT_URI` must exactly match the values in the Spotify app dashboard, including scheme, host, port, and path.
- The frontend dev server proxies API paths, but the production container talks to the API via an internal Docker service name (`api`) configured in Nginx.
- The backend `clear-songs/README.md` contains older/legacy notes (including `src/main.go` and older route paths). Prefer this `QWEN.md`, `AGENTS.md`, and the actual source code for current paths.

## Useful files to know

| File | Why it matters |
|------|----------------|
| `clear-songs/cmd/server/main.go` | Backend server entrypoint. |
| `clear-songs/internal/infrastructure/transport/http/routes.go` | Where all API routes are registered. |
| `clear-songs/internal/infrastructure/di/*.go` | Dependency injection wiring. |
| `clear-songs-front/src/main.ts` | Frontend entrypoint. |
| `clear-songs-front/src/app/app.routes.ts` | Frontend routing. |
| `clear-songs-front/tools/generate-env.js` | Generates `environment.auto.ts`. |
| `clear-songs-front/proxy.conf.json` | Dev-server API proxy rules. |
| `.env.example` | Template for required environment variables. |
| `docker-compose.yml` (root) | Full-stack orchestration. |
| `AGENTS.md` | Additional agent-oriented project notes. |
