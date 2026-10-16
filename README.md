# Calculator

Go REST backend and React + TypeScript calculator, ready for Docker. Go processes every calculator action through POST /api/calculator/action. React only sends actions, formats text, and renders the returned state. See backend/README.md for the API contract.

## Run with Docker Desktop

Start Docker Desktop, then run from this directory:

```sh
docker compose up --build -d --wait
```

- Frontend: http://localhost:3000
- Backend health check: http://localhost:8080/api/health
- The frontend forwards `/api/*` to the backend, including `/api/health`.
- Local images: `calculator-frontend:local` and `calculator-backend:local`.

Ports bind to the local machine. Set `FRONTEND_PORT` or `BACKEND_PORT` to
override them in your shell or a root `.env` file (see `.env.example`). A local
`.env` is ignored by Git and sits outside both Docker build contexts.
For example, if a development backend already uses port 8080:

```sh
BACKEND_PORT=8081 docker compose up --build -d --wait
```

In PowerShell, set `$env:BACKEND_PORT = '8081'` before the Compose command.
The backend container still listens on port 8080 internally.

Check status and logs with `docker compose ps` and `docker compose logs`.
Stop the containers with `docker compose down`.

### Builds and ignored files

Each service uses its own directory as the build context and its own
`.dockerignore`. These files allow only source, tests, and required build
configuration; local dependencies, compiled output, Git metadata, editor files,
and environment secrets are excluded from the context.

Both images use multiple build stages. Go tests and vet run during the backend
build; frontend tests and the TypeScript/Vite build run during the frontend
build. The backend runtime contains only the compiled binary and runs as a
non-root user. The frontend runtime contains Nginx and the production assets.
Its health check verifies that Nginx can reach the Go health endpoint.

Rebuild after source changes with `docker compose up --build -d --wait`.

## Local development

Use Go 1.27+ and Node.js 24 LTS.

Backend:

```sh
cd backend
go run ./cmd/api
```

Frontend, in another terminal:

```sh
cd frontend
npm ci
npm run dev
```

Open http://localhost:5173. Vite also forwards `/api/*` to the backend on
port 8080.

## Checks

```sh
cd backend
go test ./...
go vet ./...
go build -o bin/server ./cmd/api
```

```sh
cd frontend
npm test
npm run build
```
