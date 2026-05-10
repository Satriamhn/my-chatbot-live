# my Chatbot Life

my Chatbot Life is a monorepo for a SaaS chatbot platform.

## What is here

- `frontend/`, React + Vite + TypeScript dashboard and public widget runtime
- `backend/`, Go + Gin + GORM API and WebSocket service
- `docker-compose.yml`, PostgreSQL 16 with `pgvector` for local development
- `RUNBOOK.md`, day to day operations and incident notes
- `docs/runbook/widget-embed.md`, widget embed rollout and contract details

## Prerequisites

- Node.js 20+ and npm
- Go 1.22+
- Docker Desktop or Docker Engine with Docker Compose
- A Gemini API key for AI backed flows

## Local setup

1. Start PostgreSQL and pgvector:

   ```bash
   docker compose up -d
   ```

   This uses the local defaults from `docker-compose.yml`: `postgres` / `postgres` / `chatbot` on port `5432`.

2. Set backend environment variables:

   ```bash
   cp backend/.env.example backend/.env
   ```

   Fill in `backend/.env` before running the backend. Keep secrets in env files and never commit them. Common values are `GEMINI_API_KEY`, `JWT_SECRET`, `DATABASE_URL`, and `WIDGET_RUNTIME_ORIGINS`.

3. Install dependencies:

   ```bash
   npm install --prefix frontend
   cd backend && go mod download
   ```

## Run locally

```bash
npm run dev --prefix frontend
cd backend && go run ./cmd/api
```

The frontend runs the Vite app. The backend runs the Gin API and WebSocket service.

## Test and build

```bash
npm run test:all
npm run build:all
```

Optional root lint check:

```bash
npm run lint:all
```

`npm run test:all` runs frontend tests when present, then backend Go tests. `npm run build:all` builds the frontend and compiles the backend API binary.

## Widget rollout

- Public widget loading is query driven: `/widget?org_id=...`
- The loader lives in `frontend/public/widget-embed.js`
- Before rollout, add the widget browser origin to `WIDGET_RUNTIME_ORIGINS`
- Read `docs/runbook/widget-embed.md` for the embed contract and rollout steps

## Operations notes

- Use `RUNBOOK.md` for setup, day to day commands, and incident handling
- AI replies still depend on configured provider credentials and available quota, so a running service does not guarantee chat responses
