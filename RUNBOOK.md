# my Chatbot Life Runbook

This document serves as the technical source of truth for operating the my Chatbot Life SaaS platform. It provides the exact commands for daily operations, an incident response ladder for troubleshooting, and an evidence index mapping build tasks to their generated proof files.

## 0. Environment Setup (First-Time)

Lakukan ini **sekali** sebelum menjalankan aplikasi untuk pertama kali.

### Step 1 — Install Docker Desktop (Windows)
1. Download dari: https://www.docker.com/products/docker-desktop/
2. Install & restart PC jika diminta
3. Buka Docker Desktop, tunggu hingga status **Running** (ikon di system tray)

### Step 2 — Jalankan PostgreSQL via Docker Compose
```bash
# Di root folder proyek (my-chatbot-live/)
docker compose up -d

# Verifikasi postgres sudah berjalan
docker ps
```
Ini akan menjalankan PostgreSQL 16 + pgvector di background. Data tersimpan di volume `pgdata` (persisten).

### Step 3 — Dapatkan Gemini API Key
1. Buka https://aistudio.google.com/app/apikey
2. Login dengan akun Google
3. Klik **"Create API Key"**
4. Copy key tersebut

### Step 4 — Konfigurasi `.env` backend
```bash
# File backend/.env sudah dibuat, tinggal isi GEMINI_API_KEY:
# Edit file: backend/.env
GEMINI_API_KEY=paste-key-kamu-di-sini
```

### Step 5 — Stop & Start Database
```bash
# Stop database
docker compose down

# Start ulang database
docker compose up -d
```

---

## 1. Operational Commands

### Workspace Setup
Install all required dependencies before starting the application.
```bash
npm install --prefix frontend
cd backend && go mod download
```

### Local Development
Run the frontend and backend servers locally.
```bash
# Start the frontend dev server (Vite)
npm run dev --prefix frontend

# Start the backend API server (Gin)
cd backend && go run cmd/api/main.go
```

### Testing
Execute unit and integration tests across both stacks.
```bash
# Run all tests via root package script
npm run test:all

# Run frontend component tests
npm run test --prefix frontend

# Run frontend end-to-end smoke tests
npm run test:e2e --prefix frontend

# Run backend tests with race detection and coverage
cd backend && go test ./... -race -cover
```

### Widget embed rollout
- See `docs/runbook/widget-embed.md` for the tenant embed contract,
  required `org_id`, and production `WIDGET_RUNTIME_ORIGINS` setup.
- The allowlist applies to the browser origin that calls
  `/api/v1/widget/*`, not the tenant site that hosts the script tag.

### Building for Production
Compile the assets for deployment.
```bash
# Build both frontend and backend
npm run build:all

# Build frontend only
npm run build --prefix frontend

# Build backend binary only
cd backend && go build -o api.exe cmd/api/main.go
```

### Database Migration
The backend initializes its schema and Row-Level Security (RLS) policies on startup or via test setup scripts. Ensure PostgreSQL with pgvector is running and the `.env` variables match your local database credentials.

## 2. Incident Response Ladder

When issues occur in production, follow this decision ladder to resolve the problem systematically.

### Level 1: Resume (Soft Recovery)
Used for transient failures like dropped connections or temporary timeouts.
- **WebSocket Connection Drops:** The frontend client will automatically attempt to reconnect with an exponential backoff. Instruct users to refresh their browser if the connection stays dead for over 30 seconds.
- **Database Lock Timeout:** If a transaction times out, the backend returns a 503 error. The retry logic in the Axios client will attempt the request again. Monitor PostgreSQL logs to ensure locks clear.
- **Action:** Wait for automated retries. Check basic network connectivity.

### Level 2: Relaunch (Process Recovery)
Used when a service is unresponsive, leaking memory, or stuck in a bad state.
- **Gin API Unresponsive:** If the `/health` endpoint fails or times out, restart the Go binary process. 
- **Frontend Hangs:** If the React application stops responding due to state corruption in the browser, instruct the user to force reload (Ctrl+F5).
- **Action:** Restart the affected process or container. Verify the service is back up using the health check endpoint.

### Level 3: Isolate (Damage Control)
Used for critical failures like data corruption, security breaches, or major dependency outages.
- **pgvector Extension Fails:** If RAG queries fail because the pgvector extension crashes, pause the AI response capability immediately. Users can still use the `Take Over (Human Mode)` feature in the `/inbox` to respond manually.
- **Tenant Data Leak Suspected:** If an `org_mismatch` 403 error rate spikes unexpectedly, immediately revoke the offending API keys or tokens. The backend middleware and database RLS will block unauthorized access, but high volume indicates an active probing attempt.
- **Action:** Disable the broken feature. Route traffic to human fallbacks. Escalate to the engineering lead for a root cause analysis.

## 3. Verification and Observability

Use these checks when confirming a deployment or debugging a live issue.

- **Start with the standard commands.** Bring the stack up with
  `docker compose up -d`, then run the test commands in section 1 before
  checking a release.
- **Watch the two app logs.** Start the frontend with
  `npm run dev --prefix frontend`. Start the backend with
  `cd backend && go run cmd/api/main.go`. If either fails, fix that service
  first.
- **Check the browser path.** Open the app, sign in, and confirm the inbox,
  knowledge base, and settings screens load without console errors.
- **Review Docker state.** Use `docker ps` to confirm PostgreSQL is running
  and that the app can connect to it.
- **For widget issues, read the dedicated guide.** See
  `docs/runbook/widget-embed.md` for the tenant embed contract, required
  `org_id`, and production `WIDGET_RUNTIME_ORIGINS` setup.
- **Escalate quickly on data or auth anomalies.** If you see unexpected
  `org_mismatch` errors or RAG failures, pause the affected feature, route
  users to the human fallback, and inspect backend logs plus database access.
