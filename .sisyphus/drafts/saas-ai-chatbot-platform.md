# Draft: SaaS AI Chatbot Platform

## Requirements (confirmed)
- "SaaS AI Chatbot Platform" multi-tenant architecture with strict pixel-perfect UI/UX
- Monorepo target structure: `/frontend` and `/backend`
- Backend target stack: Golang (Fiber atau Gin), standard layout (`/cmd`, `/internal`), PostgreSQL (GORM), Vector DB (Supabase/Pinecone), OpenAI API
- Frontend target stack: React (Vite) + Tailwind CSS, adapting uploaded TailAdmin template
- Backend outputs requested:
  - Mermaid end-to-end system flowchart (human_assigned vs bot_handling, RAG by organization_id, prompt assembly, WebSocket response)
  - Mermaid ERD for Organizations, Bot_Settings, Knowledge_Bases, Chat_Sessions, Messages (tenant isolation by organization_id)
  - Production-ready GORM structs under `internal/models/` with JSON tags
- Frontend outputs requested:
  - `/settings/bot` page with synchronized live widget preview
  - `/knowledge` page with action buttons + trained-data table tailored for chatbot digital data
  - `/inbox` page with visual role differentiation (User/AI/Human), AI metadata, and Take Over human mode behavior
  - `services/` Axios skeleton carrying JWT token (+ org_id) to Go backend
  - Add menu "hubungi saya" and owner profile info section
- User asks for step-by-step file generation/modification guidance and strict typography/layout consistency
- Product/platform name locked: `my Chatbot Life` (logo asset will be added later by owner)
- Owner contact data confirmed for "Hubungi Saya":
  - Name: Satria Mahendra
  - Role: Admin my Chatbot Life
  - Email: satriamahendra1406@gmail.com
  - WhatsApp: +6281226251707

## Technical Decisions
- Planning-only mode: produce decision-complete execution plan for sub-agent team, not direct source-code implementation in this session
- Use asynchronous sub-agent orchestration for repo discovery and architecture decisions
- Wave priority from prior orchestration setup: Audit & quick wins → Stabilize first
- Target orchestration intensity: Balanced 6-8 agents
- Backend framework selected: Gin
- Vector database selected (phase awal): Supabase pgvector
- Owner section mode selected: real owner data in first wave (not placeholder)
- Transport for flow design anchored to user requirement: WebSocket response path
- Tenant propagation contract selected: JWT claim + `X-Org-ID` header
- Owner profile field set selected: Name + Role + Email + WhatsApp

## Research Findings
- Existing repo currently observed with `template/` as the main code container
- Prior architecture mapping identified TailAdmin app root at `template/free-react-tailwind-admin-dashboard-main`
- Current monorepo state: root-level `/frontend` and `/backend` do not exist yet; only implementation root found is `template/free-react-tailwind-admin-dashboard-main`
- Tooling surface (frontend template): `dev`, `build`, `lint`, `preview` scripts in `template/free-react-tailwind-admin-dashboard-main/package.json`; no CI workflow files and no Docker/compose files detected
- Test infra scan: no test framework config, no test scripts, no test files/mocks/fixtures, no CI test integration detected
- Oracle architecture review delivered defaults and risk controls for multi-tenant isolation, idempotency, ordering, workflow state machine, and UI precision gates
- Frontend pattern mapping complete with concrete references:
  - `src/components/common/ComponentCard.tsx` (card standard)
  - `src/components/ui/button/Button.tsx` (button standard)
  - `src/components/ui/table/index.tsx` + `src/components/tables/BasicTables/BasicTableOne.tsx` (table standard)
  - `src/components/form/Label.tsx`, `InputField.tsx`, `Select.tsx` + `src/pages/Forms/FormElements.tsx` (form composition)
  - `src/components/header/NotificationDropdown.tsx` as closest chat-list seed for inbox UI
- Async orchestration monitor loop blueprint completed with state machine, retry/backoff policy, heartbeat/stall handling, and operator status payload
- Timeout recovery ladder completed with deterministic resume/relaunch/isolate protocol and incident evidence requirements
- API/auth integration scan findings:
  - Current template has no Axios client/interceptor/auth service layer yet
  - Auth pages (`src/components/auth/SignInForm.tsx`, `SignUpForm.tsx`) are presentational (no API wiring)
  - Recommended global injection anchor for auth/org context is provider composition at `src/main.tsx`
- Routing/navigation injection map completed:
  - Route registry: `src/App.tsx`
  - Sidebar data and render: `src/layout/AppSidebar.tsx`
  - Icon registry: `src/icons/index.ts`
  - Owner dropdown insertion point: `src/components/header/UserDropdown.tsx`
- Safety audit findings:
  - Current auth in template is UI scaffold, not real session lifecycle implementation
  - `.env` policy/template and operational runbooks are not yet established in repo
  - Guardrail needs confirmed: command safety gate, session lifecycle gate, recovery loop gate, secrets/env gate, incident gate
- Backend scaffold assessment complete:
  - No existing `/backend`, no `go.mod`, no `*.go` files, no `/cmd` or `/internal` implementation yet
  - Backend architecture work is net-new bootstrap from repository root
- External production guidance (`bg_75dbcdef`) consolidated:
  - Gin remains preferred default for lower integration risk and net/http compatibility
  - Multi-tenant data isolation should combine GORM tenant scopes + Postgres RLS
  - WebSocket implementation should follow read/write pump model with origin/auth/session controls
  - Supabase pgvector first, with Pinecone namespace strategy for high-scale tenant carve-outs
  - Security priorities: tenant leakage prevention, prompt-injection defenses, websocket auth/rate-limiting

## Open Questions
- Owner profile/contact details content to display in menu and profile panel (real data requested)
- Test strategy with currently empty test stack: bootstrap tests now vs phase after core feature delivery

## Scope Boundaries
- INCLUDE: full execution blueprint for backend + frontend + orchestration safety/recovery
- EXCLUDE: direct implementation in source code paths during planning phase
