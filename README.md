# Clear Songs (complete)

**Languages:** [English](#english) · [Italiano](#italiano)

<a id="english"></a>

## English

App to inspect and manage your **saved** Spotify library: dashboard with per-artist summaries (including genre filters), bulk deletion by artist or by track-count ranges, and playlist operations. **Go** (Gin) backend, **Angular** frontend, session and cache in **Redis**, optional backup and persistence in **PostgreSQL**.

### Repository layout

| Folder | Contents |
|--------|----------|
| `clear-songs/` | HTTP API, Spotify OAuth, optional Gemini integration for genre fallback |
| `clear-songs-front/` | Angular SPA (build uses `API_URL` pointing at the backend) |
| `docker-compose.yml` (repo root) | Full stack: API, frontend, PostgreSQL, Redis |

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose (for the full stack), or Go 1.25+, a Node.js version compatible with Angular 20, and local PostgreSQL and Redis if you run the API outside containers.

### Configuration

1. Copy the example env file and fill in real values:

   ```bash
   cp .env.example .env
   ```

2. In the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard), create an app and set **Redirect URIs** to match **exactly** what you use in `.env`:

   - Backend: `REDIRECT_URI` (e.g. `http://127.0.0.1:3000/auth/callback`).
   - Frontend: `SPOTIFY_REDIRECT_URI` (e.g. `http://127.0.0.1:4200/callback`).

3. Optional: a **Google Gemini** API key (`GEMINI_API_KEY`) to classify artists whose Spotify tags do not map to your genre model in the summary. See comments in `.env.example` for batching, timeouts, and per-artist Redis cache.

### Run with Docker (recommended)

From the repository root (where `docker-compose.yml` lives):

```bash
docker compose up --build
```

Typical services:

- API: `http://127.0.0.1:3000`
- Frontend: `http://127.0.0.1:4200` (host port mapped to the container)
- PostgreSQL: port `5432` on the host
- Redis: port `6379` on the host

`docker-compose.yml` uses `env_file: .env` for the API; ensure `CLIENT_ID`, `CLIENT_SECRET`, and redirect URIs match your Spotify app.

Stop:

```bash
docker compose down
```

### Local development (without rebuilding the whole stack)

**Backend** (`clear-songs/`):

```bash
cd clear-songs
go mod download
go run ./cmd/server/main.go
```

Requires environment variables consistent with `.env` (including `DB_*`, `REDIS_*` if you use local database and Redis). Server entrypoint: `cmd/server/main.go`.

**Frontend** (`clear-songs-front/`):

```bash
cd clear-songs-front
npm install
npm start
```

The `prestart` script generates config from `.env` in the frontend folder (see `tools/generate-env.js`). Set `API_URL` to your backend (e.g. `http://127.0.0.1:3000`).

### Tests

```bash
cd clear-songs && go test $(go list ./... | grep -v '/test/e2e')
cd clear-songs-front && npm test
```

End-to-end tests under `clear-songs/test/e2e` may need real services and auth; skipping them in local CI is normal if not configured.

### Further documentation

- Environment template: [`.env.example`](.env.example)
- Older backend API notes: [`clear-songs/README.md`](clear-songs/README.md) (some sections may describe legacy layout or commands; prefer this README and the source for current paths)

---

<a id="italiano"></a>

## Italiano

Applicazione per analizzare e gestire la libreria Spotify salvata: dashboard con riepilogo per artista (anche per genere), eliminazione bulk per artista o per intervalli di conteggio, operazioni sulle playlist. Backend in **Go** (Gin), frontend **Angular**, dati di sessione e cache in **Redis**, backup/opzionale persistenza in **PostgreSQL**.

### Struttura del repository

| Cartella | Contenuto |
|----------|-----------|
| `clear-songs/` | API HTTP, OAuth Spotify, integrazione opzionale Gemini per il fallback genere |
| `clear-songs-front/` | SPA Angular (build con `API_URL` verso il backend) |
| `docker-compose.yml` (root) | Stack: API, frontend, PostgreSQL, Redis |

### Prerequisiti

- [Docker](https://docs.docker.com/get-docker/) e Docker Compose (per lo stack completo), oppure Go 1.25+, Node.js compatibile con Angular 20 ed istanze locali di PostgreSQL e Redis per lo sviluppo senza container API.

### Configurazione

1. Copia le variabili d’esempio e compila i valori reali:

   ```bash
   cp .env.example .env
   ```

2. Nel [Dashboard Spotify](https://developer.spotify.com/dashboard) crea un’app e configura i **Redirect URI** in modo che coincidano **esattamente** con quanto usato in `.env`:

   - Backend: `REDIRECT_URI` (es. `http://127.0.0.1:3000/auth/callback`).
   - Frontend: `SPOTIFY_REDIRECT_URI` (es. `http://127.0.0.1:4200/callback`).

3. Opzionale: chiave **Google Gemini** (`GEMINI_API_KEY`) per classificare gli artisti senza tag Spotify mappabili nel riepilogo. Vedi commenti in `.env.example` per batch, timeout e cache Redis per artista.

### Avvio con Docker (consigliato)

Dalla root del repository (dove si trova `docker-compose.yml`):

```bash
docker compose up --build
```

Servizi tipici:

- API: `http://127.0.0.1:3000`
- Frontend: `http://127.0.0.1:4200` (porta host mappata sul container)
- PostgreSQL: porta `5432` sul host
- Redis: porta `6379` sul host

Il file `docker-compose.yml` carica `env_file: .env` per l’API; assicurati che `CLIENT_ID`, `CLIENT_SECRET` e gli URI di redirect siano coerenti con l’app Spotify.

Arresto:

```bash
docker compose down
```

### Sviluppo locale (senza ricostruire tutto lo stack)

**Backend** (`clear-songs/`):

```bash
cd clear-songs
go mod download
go run ./cmd/server/main.go
```

Richiede variabili d’ambiente allineate a `.env` (inclusi `DB_*`, `REDIS_*` se usi database e Redis locali). Il percorso del server è `cmd/server/main.go`.

**Frontend** (`clear-songs-front/`):

```bash
cd clear-songs-front
npm install
npm start
```

Lo script `prestart` genera la configurazione da `.env` nella cartella frontend (vedi `tools/generate-env.js`). Imposta `API_URL` al backend (es. `http://127.0.0.1:3000`).

### Test

```bash
cd clear-songs && go test $(go list ./... | grep -v '/test/e2e')
cd clear-songs-front && npm test
```

I test end-to-end in `clear-songs/test/e2e` possono richiedere servizi reali e autenticazione; escluderli è normale in CI locale se non configurati.

### Documentazione aggiuntiva

- Esempio variabili: [`.env.example`](.env.example)
- Dettagli storici sull’API backend: [`clear-songs/README.md`](clear-songs/README.md) (alcune sezioni possono riferirsi a layout o comandi legacy; fare riferimento a questo README e al codice per i path attuali)
