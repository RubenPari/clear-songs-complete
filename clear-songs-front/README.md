# ClearSongsFront

Angular 20 frontend for Clear Songs, providing a UI to manage your Spotify library.

This project was generated using [Angular CLI](https://github.com/angular/angular-cli) version 20.1.6.

## Environment Setup

### Configuration

1. **Copy environment example**
   ```bash
   cp .env.example .env
   ```

2. **Edit `.env` with your settings**
   ```env
   API_URL=http://127.0.0.1:3000
   SPOTIFY_CLIENT_ID=your_spotify_client_id
   SPOTIFY_REDIRECT_URI=http://127.0.0.1:4200/callback
   ```

   The build scripts automatically generate `src/environments/environment.auto.ts` from `.env`.

## Development server

To start a local development server, run:

```bash
npm install
npm start
```

Once the server is running, open your browser and navigate to `http://127.0.0.1:4200/`. The application will automatically reload whenever you modify any of the source files.

The `start` script:
1. Generates environment config from `.env` using `tools/generate-env.js`
2. Runs Angular dev server on `http://127.0.0.1:4200`

## Code scaffolding

Angular CLI includes powerful code scaffolding tools. To generate a new component, run:

```bash
ng generate component component-name
```

For a complete list of available schematics (such as `components`, `directives`, or `pipes`), run:

```bash
ng generate --help
```

## Building

To build the project run:

```bash
npm run build
```

This will compile your project and store the build artifacts in the `dist/` directory. By default, the production build optimizes your application for performance and speed.

The `build` script:
1. Generates environment config from `.env`
2. Runs Angular build

## 🐳 Docker

### Build and Run Docker Image

1. **Build the Docker image**
   ```bash
   docker build -t clear-songs-front:latest .
   ```

2. **Run the container**
   ```bash
   docker run --rm -p 4200:80 clear-songs-front:latest
   ```

   The app will be available at `http://127.0.0.1:4200`

### Docker Image Details

- **Build stage**: Node 20 Alpine
  - Installs dependencies (`npm ci`)
  - Runs `tools/generate-env.js` to generate config from `.env`
  - Builds the Angular app (`npm run build`)
  - Output: `dist/clear-songs-front`

- **Runtime stage**: Nginx Alpine
  - Serves the built app via Nginx
  - SPA routing configured (fallback to `index.html`)
  - Exposes port `80`

### Multi-stage Build Benefits

- Final image is lightweight (Nginx only, no Node)
- Fast production deployment
- Optimized for container orchestration (Kubernetes, etc.)

### Running with Full Stack (Docker Compose)

If you want to run backend and frontend together:

```bash
# From the root (clear-songs-full-stack)
cd clear-songs
docker compose up --build

# In another terminal, from frontend folder
cd clear-songs-front
docker build -t clear-songs-front:latest . && docker run --rm --network host -p 4200:80 clear-songs-front:latest
```

Then access:
- **Frontend**: `http://127.0.0.1:4200`
- **API**: `http://127.0.0.1:3000`

## Running unit tests

To execute unit tests with the [Karma](https://karma-runner.github.io) test runner, use the following command:

```bash
ng test
```

## Running end-to-end tests

For end-to-end (e2e) testing, run:

```bash
ng e2e
```

Angular CLI does not come with an end-to-end testing framework by default. You can choose one that suits your needs.

## Additional Resources

For more information on using the Angular CLI, including detailed command references, visit the [Angular CLI Overview and Command Reference](https://angular.dev/tools/cli) page.
