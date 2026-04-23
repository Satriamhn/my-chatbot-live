# my Chatbot Life — Multi-Tenant SaaS Chatbot Platform Execution Plan

## TL;DR
> **Summary**: Bootstrap a new monorepo (`/frontend` + `/backend`) from the current TailAdmin template baseline, then deliver a production-oriented multi-tenant chatbot system with Gin, PostgreSQL/GORM, Supabase pgvector RAG, WebSocket inbox, and strict tenant isolation.
> **Deliverables**:
> - Backend architecture artifacts (Mermaid flowchart + ERD) and GORM model layer
> - Frontend pages: `/settings/bot`, `/knowledge`, `/inbox` with pixel-consistent TailAdmin adaptation
> - API service layer (Axios JWT + `X-Org-ID`) and owner-contact menu/profile
> - Wave-1 testing infrastructure + CI gates + security guardrails
> **Effort**: XL
> **Parallel**: YES - 5 waves
> **Critical Path**: T1 → T4 → T5 → T6 → T7 → T11 → T12 → T13 → T18 → T19

## Context
### Original Request
User meminta eksekusi arsitektur lengkap “SaaS AI Chatbot Platform” multi-tenant dengan standar UI/UX ketat (pixel-perfect), memanfaatkan template TailAdmin yang sudah ada, serta output backend/frontend rinci termasuk Mermaid flowchart, ERD, GORM structs, halaman UI inti, integrasi API, dan menu “Hubungi Saya”.

### Interview Summary
- Product name dikunci: **my Chatbot Life** (logo akan ditambahkan owner belakangan).
- Backend default dipilih: **Gin**.
- Vector DB default fase awal: **Supabase pgvector**.
- Tenant propagation dikunci: **JWT claim + `X-Org-ID` header**.
- Test strategy dikunci: **bootstrap test infra di wave pertama**.
- Owner profile final untuk wave pertama:
  - Name: Satria Mahendra
  - Role: Admin my Chatbot Life
  - Email: satriamahendra1406@gmail.com
  - WhatsApp: +6281226251707

### Metis Review (gaps addressed)
Metis menandai risiko tenant leakage, idempotency, prompt injection, dan ambiguity tenancy/auth contract. Plan ini sudah mengunci default berikut tanpa pertanyaan tambahan:
- JWT authority: backend memverifikasi claim tenant, `X-Org-ID` hanya active-tenant selector.
- Mismatch `org_id` claim vs `X-Org-ID` wajib **403 org_mismatch** + audit log.
- RLS PostgreSQL diterapkan di wave foundational data security (bukan ditunda).
- Scope exclusions eksplisit: **tidak** mencakup billing, SSO enterprise, marketing automation, dan mobile app.

## Work Objectives
### Core Objective
Menghasilkan fondasi dan fitur inti platform chatbot multi-tenant my Chatbot Life yang siap dieksekusi agent tanpa keputusan tambahan, dengan isolasi tenant ketat, UI presisi berbasis TailAdmin, dan jalur backend-frontend terintegrasi.

### Deliverables
1. Monorepo root dengan `/frontend` dan `/backend`.
2. Backend Gin scaffold + config + auth tenant middleware + database layer + RLS + GORM models.
3. Artifact arsitektur: Mermaid flowchart E2E + Mermaid ERD.
4. API domain: bot settings, knowledge base, chat sessions/messages, websocket messaging, RAG retrieval pipeline.
5. Frontend integration layer (Axios + org context) dan tiga halaman inti (`/settings/bot`, `/knowledge`, `/inbox`).
6. Owner contact/profile menu untuk Satria Mahendra.
7. Test infra + CI baseline + security and reliability gates.

### Definition of Done (verifiable conditions with commands)
- `npm --prefix "frontend" run lint` exits 0.
- `npm --prefix "frontend" run test` exits 0.
- `npm --prefix "frontend" run build` exits 0.
- `go test ./... -race -cover` from `backend` exits 0.
- `curl` org mismatch test returns HTTP 403 with `org_mismatch` error code.
- WebSocket authenticated chat flow returns AI response payload containing source/confidence metadata.
- RLS verification tests prove cross-tenant records are not accessible.

### Must Have
- `organization_id` enforced in all tenant-owned entities and access paths.
- JWT claim + `X-Org-ID` validation contract implemented consistently.
- Pixel-consistent reuse of TailAdmin primitives (`ComponentCard`, `Button`, table primitives, form inputs).
- “Take Over (Human Mode)” disables bot input in `/inbox`.
- Live widget preview synced with bot configuration form.

### Must NOT Have (guardrails, AI slop patterns, scope boundaries)
- No tenant authority from header-only context.
- No unscoped GORM queries on tenant tables.
- No direct SQL string interpolation for tenant filters.
- No hardcoded UI offsets that break template spacing conventions unless explicitly justified.
- No addition of out-of-scope domains (billing/SSO/mobile/marketing modules).

## Verification Strategy
> ZERO HUMAN INTERVENTION — all verification is agent-executed.
- Test decision: **TDD** with wave-1 bootstrap.
- Frontend framework: Vitest + React Testing Library + Playwright smoke.
- Backend framework: Go `testing` + `httptest` + DB integration/policy checks.
- QA policy: Every task includes happy + failure/edge scenario.
- Evidence: `.sisyphus/evidence/task-{N}-{slug}.{ext}`.

## Execution Strategy
### Parallel Execution Waves
> Target: 5-8 tasks per wave. Dependencies extracted into early waves.

Wave 1 (Foundation): T1, T2, T3, T4, T8
Wave 2 (Security/Data Core): T5, T6, T7, T9
Wave 3 (Conversation & RAG Core): T10, T11, T12, T13
Wave 4 (Frontend Experience): T14, T15, T16, T17, T18
Wave 5 (Hardening & Release Readiness): T19, T20

### Dependency Matrix (full, all tasks)
| Task | Depends On | Blocks |
|---|---|---|
| T1 | - | T2,T3,T4,T14-T18 |
| T2 | T1 | T5,T6,T9-T13,T19 |
| T3 | T1 | T16,T17,T18,T19 |
| T4 | T1 | T5,T6,T7,T9-T13,T19 |
| T5 | T4,T2 | T9,T10,T11,T12,T13,T19 |
| T6 | T4,T2 | T7,T9,T10,T11,T12,T13,T19 |
| T7 | T6 | T9,T10,T11,T12,T13,T19 |
| T8 | T4,T6,T7 | T20 |
| T9 | T5,T6,T7 | T16,T19 |
| T10 | T5,T6,T7 | T17,T19 |
| T11 | T5,T6,T7 | T12,T13,T18,T19 |
| T12 | T5,T6,T7,T11 | T13,T18,T19 |
| T13 | T11,T12 | T18,T19 |
| T14 | T1,T2 | T15,T16,T17,T18,T19 |
| T15 | T1,T14 | T16,T17,T18,T19 |
| T16 | T3,T9,T14,T15 | T19 |
| T17 | T3,T10,T14,T15 | T19 |
| T18 | T3,T11,T12,T13,T14,T15 | T19 |
| T19 | T2,T3,T5,T6,T7,T9-T18 | T20,F1-F4 |
| T20 | T8,T19 | F1-F4 |

### Agent Dispatch Summary (wave → task count → categories)
- Wave 1 → 5 tasks → implementation(2), testing(1), security(1), writing(1)
- Wave 2 → 4 tasks → security(1), implementation(2), refactorer/data(1)
- Wave 3 → 4 tasks → implementation(3), deep(1)
- Wave 4 → 5 tasks → visual-engineering(4), implementation(1)
- Wave 5 → 2 tasks → security-auditor(1), writing(1)

## TODOs
> Implementation + Test = ONE task. Never separate.
> EVERY task includes agent profile, parallelization, references, acceptance criteria, and QA scenarios.

<!-- TASKS_INSERT_HERE -->
- [ ] 17. Build `/knowledge` page with action buttons and trained-data table

  **What to do**: Implement knowledge management page with three actions (`Upload Document`, `URL Sync`, `Manual Q&A`) and tailored trained-data table columns (`Nama File`, `Tipe`, `Status`, `Tanggal`) using TailAdmin table primitives.
  **Must NOT do**: Do not keep irrelevant template columns (e.g., product qty/signature fields).

  **Recommended Agent Profile**:
  - Category: `visual-engineering` — Reason: table precision and data UX adaptation.
  - Skills: [`ocs-delegation-gate`] — maintain strict scope and quality.
  - Omitted: [`ocs-release-integrity`] — no release actions in task.

  **Parallelization**: Can Parallel: YES | Wave 4 | Blocks: T19 | Blocked By: T3,T10,T14,T15

  **References**:
  - Pattern: `frontend/src/components/ui/table/index.tsx` — base table primitives.
  - Pattern: `frontend/src/components/tables/BasicTables/BasicTableOne.tsx` — row/cell styling pattern.
  - Pattern: `frontend/src/components/ecommerce/RecentOrders.tsx` — card+table composition.

  **Acceptance Criteria**:
  - [ ] Three actions render with consistent button styles.
  - [ ] Table renders only chatbot-relevant columns and data.
  - [ ] Status badges map correctly (`queued`,`processing`,`ready`,`failed`).

  **QA Scenarios**:
  ```
  Scenario: Knowledge table data integrity
    Tool: Playwright
    Steps: Load seeded knowledge dataset and inspect table headers/rows
    Expected: Only required four columns visible with accurate values
    Evidence: .sisyphus/evidence/task-17-knowledge-page.png

  Scenario: Invalid action payload handling
    Tool: Playwright
    Steps: Trigger URL Sync with invalid URL
    Expected: Validation error shown; no invalid row inserted
    Evidence: .sisyphus/evidence/task-17-knowledge-page-error.txt
  ```

  **Commit**: YES | Message: `feat(frontend): add knowledge base page and trained data table` | Files: [`frontend/src/pages/knowledge/**`, `frontend/src/components/knowledge/**`]

- [ ] 18. Build `/inbox` omnichannel chat page with role differentiation and takeover behavior

  **What to do**: Create chat window UI that visually distinguishes User vs AI Bot vs Human Agent messages; show AI metadata (`source file`, `confidence rate`); implement `Take Over (Human Mode)` button to disable bot input behavior and reflect state.
  **Must NOT do**: Do not use dropdown-specific hardcoded offsets from header notification component.

  **Recommended Agent Profile**:
  - Category: `visual-engineering` — Reason: complex chat UX/state rendering.
  - Skills: [`ocs-delegation-gate`] — keeps behavior aligned with backend state machine.
  - Omitted: [`ocs-openai-multi-account`] — not required for inbox UI.

  **Parallelization**: Can Parallel: NO | Wave 4 | Blocks: T19 | Blocked By: T3,T11,T12,T13,T14,T15

  **References**:
  - Pattern: `frontend/src/components/header/NotificationDropdown.tsx` — message-row visual seed.
  - Pattern: `frontend/src/components/common/ComponentCard.tsx` — container baseline.
  - API/Type: backend message payload from T13 (metadata fields).

  **Acceptance Criteria**:
  - [ ] Message roles are visually distinct with consistent typography and spacing.
  - [ ] AI messages render source/confidence metadata.
  - [ ] Human takeover disables bot-mode input path and updates status indicator.

  **QA Scenarios**:
  ```
  Scenario: Chat role rendering and metadata visibility
    Tool: Playwright
    Steps: Open inbox with seeded mixed-role conversation
    Expected: User/AI/Human styles differ; AI metadata visible under AI messages
    Evidence: .sisyphus/evidence/task-18-inbox-ui.png

  Scenario: Takeover disables bot input
    Tool: Playwright
    Steps: Click `Take Over (Human Mode)` and attempt bot input action
    Expected: Bot input disabled/block message shown until mode switched back
    Evidence: .sisyphus/evidence/task-18-inbox-ui-error.txt
  ```

  **Commit**: YES | Message: `feat(frontend): add omnichannel inbox with takeover mode` | Files: [`frontend/src/pages/inbox/**`, `frontend/src/components/inbox/**`]

- [ ] 19. Run cross-layer hardening tests and reliability gates

  **What to do**: Implement and execute integration/E2E suites for tenant isolation, auth mismatch, websocket auth, RAG safety, and takeover flow; wire CI gates for frontend+backend+security checks.
  **Must NOT do**: Do not mark task complete with partial green checks.

  **Recommended Agent Profile**:
  - Category: `security-auditor` — Reason: cross-domain validation and risk closure.
  - Skills: [`ocs-delegation-gate`] — strict verification and evidence quality.
  - Omitted: [`ocs-runtime-validation`] — optional add-on, not mandatory if criteria met.

  **Parallelization**: Can Parallel: NO | Wave 5 | Blocks: T20,F1-F4 | Blocked By: T2,T3,T5,T6,T7,T9-T18

  **References**:
  - External: `https://cheatsheetseries.owasp.org/cheatsheets/WebSocket_Security_Cheat_Sheet.html`
  - External: `https://genai.owasp.org/llmrisk/llm01-prompt-injection/`
  - Pattern: T5/T6 contract tests for org mismatch and RLS.

  **Acceptance Criteria**:
  - [ ] `go test ./... -race -cover` passes.
  - [ ] `npm --prefix "frontend" run test` and smoke E2E pass.
  - [ ] CI config enforces blocking checks for lint/test/build on both stacks.

  **QA Scenarios**:
  ```
  Scenario: Full reliability gate pass
    Tool: Bash
    Steps: Run backend and frontend full test commands plus CI-local validation script
    Expected: All checks pass with exit code 0
    Evidence: .sisyphus/evidence/task-19-hardening.txt

  Scenario: Forced tenant-leak test
    Tool: Bash
    Steps: Execute adversarial test attempting cross-tenant resource access
    Expected: Access denied and test marked pass for security behavior
    Evidence: .sisyphus/evidence/task-19-hardening-error.txt
  ```

  **Commit**: YES | Message: `test(security): add cross-layer tenant and reliability gates` | Files: [`backend/tests/**`, `frontend/tests/**`, `.github/workflows/**`]

- [ ] 20. Produce implementation runbook and evidence index for handoff

  **What to do**: Write final technical runbook covering setup, command matrix, troubleshooting, timeout recovery ladder, monitor control-loop operations, and evidence index mapping for all tasks.
  **Must NOT do**: Do not leave undocumented manual-only steps.

  **Recommended Agent Profile**:
  - Category: `writing` — Reason: high-precision operational documentation.
  - Skills: [`ocs-delegation-gate`] — consistent and verifiable documentation output.
  - Omitted: [`ocs-installer-copy-seo`] — non-marketing technical docs.

  **Parallelization**: Can Parallel: NO | Wave 5 | Blocks: F1-F4 | Blocked By: T8,T19

  **References**:
  - Pattern: `.sisyphus/evidence/**` — evidence archive outputs.
  - External: Metis recovery directives from session review.

  **Acceptance Criteria**:
  - [ ] Runbook includes exact commands for local run/test/build/migrate.
  - [ ] Incident response section documents resume/relaunch/isolate decision ladder.
  - [ ] Evidence index maps each task to generated proof files.

  **QA Scenarios**:
  ```
  Scenario: Clean-room runbook execution
    Tool: Bash
    Steps: Follow runbook commands in fresh environment
    Expected: Environment can be setup and validated without missing steps
    Evidence: .sisyphus/evidence/task-20-runbook.txt

  Scenario: Recovery procedure simulation
    Tool: Bash
    Steps: Simulate stalled task condition and follow documented ladder
    Expected: Procedure reaches deterministic resume/relaunch outcome
    Evidence: .sisyphus/evidence/task-20-runbook-error.txt
  ```

  **Commit**: YES | Message: `docs(ops): add execution runbook and evidence index` | Files: [`docs/runbook/**`, `.sisyphus/evidence/index.md`]

<!-- TASKS_INSERT_HERE -->
- [ ] 13. Implement WebSocket delivery pipeline for AI responses

  **What to do**: Add authenticated WebSocket hub with per-connection tenant/session context; stream chat responses including `source_file` and `confidence_rate`; follow read/write pump concurrency model.
  **Must NOT do**: Do not allow unauthenticated socket upgrades or mixed-tenant channels.

  **Recommended Agent Profile**:
  - Category: `implementation` — Reason: realtime transport implementation.
  - Skills: [`ocs-delegation-gate`] — enforce safety constraints.
  - Omitted: [`ocs-runtime-validation`] — runtime deep validation in T19.

  **Parallelization**: Can Parallel: YES | Wave 3 | Blocks: T18,T19 | Blocked By: T11,T12

  **References**:
  - External: `https://github.com/gorilla/websocket/blob/e064f32e3674d9d79a8fd417b5bc06fa5c6cad8f/doc.go#L122-L134`
  - External: `https://cheatsheetseries.owasp.org/cheatsheets/WebSocket_Security_Cheat_Sheet.html`

  **Acceptance Criteria**:
  - [ ] WS upgrade requires valid JWT and tenant context.
  - [ ] AI response payload includes message text + metadata (`source_file`, `confidence_rate`).
  - [ ] Ping/pong and deadline handling tested.

  **QA Scenarios**:
  ```
  Scenario: Authenticated websocket message flow
    Tool: Bash
    Steps: Connect WS with valid token; send user message; wait for AI response
    Expected: Receive response event with metadata fields present
    Evidence: .sisyphus/evidence/task-13-websocket.txt

  Scenario: Unauthorized websocket rejected
    Tool: Bash
    Steps: Attempt WS upgrade without token
    Expected: Upgrade denied with 401/403
    Evidence: .sisyphus/evidence/task-13-websocket-error.txt
  ```

  **Commit**: YES | Message: `feat(realtime): add authenticated websocket ai response pipeline` | Files: [`backend/internal/realtime/**`, `backend/internal/handlers/ws/**`]

- [ ] 14. Build frontend API client layer (Axios + JWT + `X-Org-ID`)

  **What to do**: Create `frontend/src/services` architecture with Axios instance, auth/org interceptors, typed service modules, and error normalization for API consumption.
  **Must NOT do**: Do not duplicate token/header injection in each page component.

  **Recommended Agent Profile**:
  - Category: `implementation` — Reason: integration foundation.
  - Skills: [`ocs-delegation-gate`] — consistency and guardrails.
  - Omitted: [`ocs-installer-copy-seo`] — unrelated.

  **Parallelization**: Can Parallel: YES | Wave 4 | Blocks: T15,T16,T17,T18,T19 | Blocked By: T1,T2

  **References**:
  - Pattern: `frontend/src/main.tsx` — provider composition insertion anchor.
  - Pattern: `frontend/src/components/auth/SignInForm.tsx` — currently presentational form requiring service integration.

  **Acceptance Criteria**:
  - [ ] Single shared Axios client handles Authorization and `X-Org-ID`.
  - [ ] Service modules exist for settings, knowledge, and inbox APIs.
  - [ ] Interceptor tests cover missing token and org mismatch responses.

  **QA Scenarios**:
  ```
  Scenario: Interceptor injects auth and org headers
    Tool: Bash
    Steps: Run unit test mocking request config and session state
    Expected: Outbound config contains Bearer token + `X-Org-ID`
    Evidence: .sisyphus/evidence/task-14-api-client.txt

  Scenario: Expired token path
    Tool: Bash
    Steps: Mock 401 response and execute request through client
    Expected: Client triggers configured auth failure handling path
    Evidence: .sisyphus/evidence/task-14-api-client-error.txt
  ```

  **Commit**: YES | Message: `feat(frontend): add axios services with jwt and x-org-id` | Files: [`frontend/src/services/**`, `frontend/src/context/**`]

- [ ] 15. Register routes, sidebar navigation, branding, and owner contact profile

  **What to do**: Add route entries (`/settings/bot`, `/knowledge`, `/inbox`) in `App.tsx`; add sidebar items in `AppSidebar.tsx`; add "Hubungi Saya" entry + owner profile in `UserDropdown.tsx`; apply product name **my Chatbot Life** and placeholder logo slot.
  **Must NOT do**: Do not hardcode owner data in multiple files—centralize config source.

  **Recommended Agent Profile**:
  - Category: `visual-engineering` — Reason: navigation UX and consistency implementation.
  - Skills: [`ocs-delegation-gate`] — keeps injection points consistent.
  - Omitted: [`frontend-ui-ux`] — optional if base TailAdmin pattern already sufficient.

  **Parallelization**: Can Parallel: YES | Wave 4 | Blocks: T16,T17,T18,T19 | Blocked By: T1,T14

  **References**:
  - Pattern: `frontend/src/App.tsx` — route registry.
  - Pattern: `frontend/src/layout/AppSidebar.tsx` — nav item arrays.
  - Pattern: `frontend/src/icons/index.ts` — icon imports/exports.
  - Pattern: `frontend/src/components/header/UserDropdown.tsx` — owner/contact menu insertion.

  **Acceptance Criteria**:
  - [ ] New routes accessible from sidebar.
  - [ ] Owner contact section displays Name/Role/Email/WhatsApp exactly as provided.
  - [ ] Product name updated to “my Chatbot Life” with placeholder logo component.

  **QA Scenarios**:
  ```
  Scenario: Navigation routes and owner menu visible
    Tool: Playwright
    Steps: Login, open sidebar and user dropdown, click each new route and owner contact item
    Expected: Correct route transitions and owner info rendering
    Evidence: .sisyphus/evidence/task-15-routing-owner.png

  Scenario: Broken route entry detection
    Tool: Playwright
    Steps: Click each new nav entry with route assertion enabled
    Expected: Test fails if any route unresolved/404
    Evidence: .sisyphus/evidence/task-15-routing-owner-error.txt
  ```

  **Commit**: YES | Message: `feat(frontend): wire chatbot routes sidebar and owner profile` | Files: [`frontend/src/App.tsx`, `frontend/src/layout/AppSidebar.tsx`, `frontend/src/components/header/UserDropdown.tsx`]

- [ ] 16. Build `/settings/bot` page with synchronized live widget preview

  **What to do**: Create bot configuration page using TailAdmin form primitives for `Bot Name`, `Welcome Message`, `System Prompt`; implement right-side live preview that updates on form input changes.
  **Must NOT do**: Do not use custom spacing/typography outside template token conventions.

  **Recommended Agent Profile**:
  - Category: `visual-engineering` — Reason: pixel-consistent form+preview UI.
  - Skills: [`ocs-delegation-gate`] — keep task bounded and verifiable.
  - Omitted: [`ocs-markdown-autofix`] — not relevant for UI code.

  **Parallelization**: Can Parallel: YES | Wave 4 | Blocks: T19 | Blocked By: T3,T9,T14,T15

  **References**:
  - Pattern: `frontend/src/pages/Forms/FormElements.tsx` — two-column form composition.
  - Pattern: `frontend/src/components/common/ComponentCard.tsx` — card container and spacing.
  - Pattern: `frontend/src/components/form/input/InputField.tsx`, `frontend/src/components/form/Label.tsx`.

  **Acceptance Criteria**:
  - [ ] Form inputs persist/load via bot settings API.
  - [ ] Live preview text reflects current form state without page reload.
  - [ ] Responsive layout matches template breakpoints.

  **QA Scenarios**:
  ```
  Scenario: Live preview sync
    Tool: Playwright
    Steps: Type into Bot Name/Welcome/System Prompt fields
    Expected: Preview panel updates text in real time
    Evidence: .sisyphus/evidence/task-16-settings-preview.mp4

  Scenario: API save failure handling
    Tool: Playwright
    Steps: Simulate 500 on save action
    Expected: Error state shown, form data retained for retry
    Evidence: .sisyphus/evidence/task-16-settings-preview-error.txt
  ```

  **Commit**: YES | Message: `feat(frontend): add bot configuration page with live preview` | Files: [`frontend/src/pages/settings/BotConfiguration.tsx`, `frontend/src/components/settings/**`]

<!-- TASKS_INSERT_HERE -->
- [ ] 9. Implement Bot Settings API with tenant-safe CRUD

  **What to do**: Create endpoints for bot config (`bot_name`, `welcome_message`, `system_prompt`) scoped by tenant; include validation and optimistic update handling.
  **Must NOT do**: Do not allow cross-tenant read/update paths.

  **Recommended Agent Profile**:
  - Category: `implementation` — Reason: domain API delivery.
  - Skills: [`ocs-delegation-gate`] — guard tenant-safe API behavior.
  - Omitted: [`ocs-openai-multi-account`] — not required for config CRUD.

  **Parallelization**: Can Parallel: YES | Wave 2 | Blocks: T16,T19 | Blocked By: T5,T6,T7

  **References**:
  - API/Type: `backend/internal/models/bot_settings.go` — model contract.
  - Pattern: `frontend/src/components/form/input/InputField.tsx` — UI fields consuming this API.

  **Acceptance Criteria**:
  - [ ] GET/PUT bot settings require auth and tenant context.
  - [ ] Validation errors return structured 422 payload.
  - [ ] Tests cover happy path + org mismatch.

  **QA Scenarios**:
  ```
  Scenario: Bot settings update success
    Tool: Bash
    Steps: Call PUT `/api/v1/settings/bot` with valid payload and matching org context
    Expected: HTTP 200 and persisted values retrievable via GET
    Evidence: .sisyphus/evidence/task-9-bot-settings-api.txt

  Scenario: Cross-tenant update denied
    Tool: Bash
    Steps: Send PUT with mismatched `X-Org-ID`
    Expected: HTTP 403 `org_mismatch`
    Evidence: .sisyphus/evidence/task-9-bot-settings-api-error.txt
  ```

  **Commit**: YES | Message: `feat(api): add tenant-scoped bot settings endpoints` | Files: [`backend/internal/handlers/settings/**`, `backend/internal/routes/**`]

- [ ] 10. Implement Knowledge Base API domain (document/url/manual metadata)

  **What to do**: Build endpoints and service contracts for `Upload Document`, `URL Sync`, `Manual Q&A` metadata entries; include status lifecycle (`queued`, `processing`, `ready`, `failed`) and tenant scoping.
  **Must NOT do**: Do not mix template ecommerce table semantics into knowledge domain.

  **Recommended Agent Profile**:
  - Category: `implementation` — Reason: service/domain modeling and API endpoints.
  - Skills: [`ocs-delegation-gate`] — maintain isolation and contract quality.
  - Omitted: [`ocs-runtime-validation`] — runtime validation later in integrated flow.

  **Parallelization**: Can Parallel: YES | Wave 3 | Blocks: T17,T19 | Blocked By: T5,T6,T7

  **References**:
  - API/Type: `backend/internal/models/knowledge_bases.go` — entity contract.
  - Pattern: `frontend/src/components/tables/BasicTables/BasicTableOne.tsx` — table adaptation reference.

  **Acceptance Criteria**:
  - [ ] Endpoint set supports create/list/update status per tenant.
  - [ ] Response schema includes file/url/manual type and timestamps.
  - [ ] Invalid type payload rejected with 422.

  **QA Scenarios**:
  ```
  Scenario: Knowledge item lifecycle update
    Tool: Bash
    Steps: Create knowledge item then update status from `queued` to `ready`
    Expected: Status transition reflected in subsequent list call
    Evidence: .sisyphus/evidence/task-10-knowledge-api.txt

  Scenario: Invalid source type rejection
    Tool: Bash
    Steps: Submit unsupported source type value
    Expected: HTTP 422 with validation details
    Evidence: .sisyphus/evidence/task-10-knowledge-api-error.txt
  ```

  **Commit**: YES | Message: `feat(api): add knowledge base metadata endpoints` | Files: [`backend/internal/handlers/knowledge/**`, `backend/internal/services/knowledge/**`]

- [ ] 11. Implement Chat Sessions + Messages API with takeover state

  **What to do**: Deliver session/message endpoints with `human_assigned` vs `bot_handling` state; enforce deterministic state transitions and tenant-safe access.
  **Must NOT do**: Do not permit message insert when session tenant does not match request context.

  **Recommended Agent Profile**:
  - Category: `implementation` — Reason: core conversation domain logic.
  - Skills: [`ocs-delegation-gate`] — manage state-transition guardrails.
  - Omitted: [`ocs-release-integrity`] — not release stage.

  **Parallelization**: Can Parallel: NO | Wave 3 | Blocks: T12,T13,T18,T19 | Blocked By: T5,T6,T7

  **References**:
  - API/Type: `backend/internal/models/chat_sessions.go` — session contract.
  - API/Type: `backend/internal/models/messages.go` — message contract.

  **Acceptance Criteria**:
  - [ ] Session state machine enforces valid transitions.
  - [ ] Human takeover endpoint flips mode and blocks bot auto-reply path.
  - [ ] Tests cover invalid transition + unauthorized tenant.

  **QA Scenarios**:
  ```
  Scenario: Human takeover flow success
    Tool: Bash
    Steps: Create session in bot mode then call takeover endpoint
    Expected: Session mode changes to `human_assigned`
    Evidence: .sisyphus/evidence/task-11-chat-domain.txt

  Scenario: Invalid transition blocked
    Tool: Bash
    Steps: Attempt illegal transition sequence defined by state rules
    Expected: HTTP 409 conflict with transition error code
    Evidence: .sisyphus/evidence/task-11-chat-domain-error.txt
  ```

  **Commit**: YES | Message: `feat(chat): add session and message domain with takeover state` | Files: [`backend/internal/handlers/chat/**`, `backend/internal/services/chat/**`]

- [ ] 12. Build RAG retrieval + prompt assembly service (tenant-filtered)

  **What to do**: Implement retrieval service using Supabase pgvector filtered by `organization_id`; assemble AI prompt with bot settings, knowledge snippets, and conversation context.
  **Must NOT do**: Do not pass untrusted retrieved content as system instruction.

  **Recommended Agent Profile**:
  - Category: `deep` — Reason: multi-component retrieval/prompt logic with security constraints.
  - Skills: [`ocs-delegation-gate`] — maintain strict boundaries and evidence.
  - Omitted: [`ocs-openai-multi-account`] — multi-account runtime not requested in this phase.

  **Parallelization**: Can Parallel: NO | Wave 3 | Blocks: T13,T18,T19 | Blocked By: T5,T6,T7,T11

  **References**:
  - External: `https://supabase.com/docs/guides/ai/vector-columns`
  - External: `https://supabase.com/docs/guides/ai/vector-indexes`
  - External: `https://cheatsheetseries.owasp.org/cheatsheets/LLM_Prompt_Injection_Prevention_Cheat_Sheet.html`

  **Acceptance Criteria**:
  - [ ] Retrieval query always scoped by tenant.
  - [ ] Prompt builder separates system instructions from retrieved/user text.
  - [ ] Tests cover injection-like payload handling.

  **QA Scenarios**:
  ```
  Scenario: Tenant-safe retrieval success
    Tool: Bash
    Steps: Seed vectors for two orgs and query as org_a
    Expected: Only org_a chunks returned to prompt builder
    Evidence: .sisyphus/evidence/task-12-rag-service.txt

  Scenario: Prompt injection payload neutralization
    Tool: Bash
    Steps: Include malicious instruction text in knowledge chunk and run prompt assembly
    Expected: Output marks chunk as data, not system override
    Evidence: .sisyphus/evidence/task-12-rag-service-error.txt
  ```

  **Commit**: YES | Message: `feat(rag): add tenant-filtered retrieval and prompt assembly` | Files: [`backend/internal/services/rag/**`, `backend/internal/services/ai/**`]

<!-- TASKS_INSERT_HERE -->
- [ ] 5. Implement JWT + `X-Org-ID` tenant contract middleware

  **What to do**: Build auth middleware that validates JWT claims (`sub`,`org_id`,`role`,`exp`) and checks `X-Org-ID` selector; reject mismatch with `403 org_mismatch`; attach tenant context to request scope.
  **Must NOT do**: Do not treat `X-Org-ID` as authority without JWT validation.

  **Recommended Agent Profile**:
  - Category: `security` — Reason: authn/authz and tenant boundary enforcement.
  - Skills: [`ocs-delegation-gate`] — strict policy and evidence tracking.
  - Omitted: [`ocs-runtime-validation`] — functional correctness first.

  **Parallelization**: Can Parallel: NO | Wave 2 | Blocks: T9,T10,T11,T12,T13,T19 | Blocked By: T4,T2

  **References**:
  - External: `https://gorm.io/docs/security.html` — query safety policy.
  - External: `https://cheatsheetseries.owasp.org/cheatsheets/WebSocket_Security_Cheat_Sheet.html` — auth/session expectations for realtime channel.

  **Acceptance Criteria**:
  - [ ] Missing/invalid JWT returns 401.
  - [ ] JWT org mismatch against `X-Org-ID` returns 403 + `org_mismatch`.
  - [ ] Authorized request exposes tenant context for downstream handlers.

  **QA Scenarios**:
  ```
  Scenario: Valid token + matching org selector
    Tool: Bash
    Steps: Call protected endpoint with valid Bearer token and matching `X-Org-ID`
    Expected: HTTP 200 and handler receives org context
    Evidence: .sisyphus/evidence/task-5-tenant-middleware.txt

  Scenario: Token/header mismatch rejection
    Tool: Bash
    Steps: Call same endpoint with token org_a and `X-Org-ID: org_b`
    Expected: HTTP 403 with `org_mismatch`
    Evidence: .sisyphus/evidence/task-5-tenant-middleware-error.txt
  ```

  **Commit**: YES | Message: `feat(auth): enforce jwt and x-org-id tenant contract` | Files: [`backend/internal/middleware/**`, `backend/internal/auth/**`]

- [ ] 6. Build PostgreSQL + GORM base with tenant scopes and RLS setup

  **What to do**: Add DB config, connection management, migration runner, reusable tenant scope helpers, and SQL migration for enabling RLS policies on tenant tables.
  **Must NOT do**: Do not rely on app-level filters only; DB-level RLS is mandatory.

  **Recommended Agent Profile**:
  - Category: `implementation` — Reason: foundational persistence layer.
  - Skills: [`ocs-delegation-gate`] — safe sequencing of schema/policy changes.
  - Omitted: [`ocs-release-integrity`] — release not yet in scope.

  **Parallelization**: Can Parallel: NO | Wave 2 | Blocks: T7,T9,T10,T11,T12,T13,T19 | Blocked By: T4,T2

  **References**:
  - External: `https://gorm.io/docs/scopes.html` — tenant scoping pattern.
  - External: `https://supabase.com/docs/guides/database/postgres/row-level-security` — RLS policy baseline.

  **Acceptance Criteria**:
  - [ ] DB boot command applies migrations including RLS policies.
  - [ ] Tenant scope helper used in repository layer for tenant-owned entities.
  - [ ] Policy test proves cross-tenant read blocked.

  **QA Scenarios**:
  ```
  Scenario: Migration and RLS initialization succeed
    Tool: Bash
    Steps: Run backend migration command against test DB
    Expected: Tables and RLS policies created successfully
    Evidence: .sisyphus/evidence/task-6-db-rls.txt

  Scenario: Cross-tenant query blocked by policy
    Tool: Bash
    Steps: Execute integration test reading org_b rows with org_a session context
    Expected: Query returns zero rows or permission error per policy
    Evidence: .sisyphus/evidence/task-6-db-rls-error.txt
  ```

  **Commit**: YES | Message: `feat(data): add gorm tenant scopes and postgres rls` | Files: [`backend/internal/db/**`, `backend/migrations/**`]

- [ ] 7. Create production-ready GORM models under `internal/models`

  **What to do**: Implement `Organizations`, `BotSettings`, `KnowledgeBases`, `ChatSessions`, `Messages` structs with `OrganizationID` (where tenant-owned), JSON tags, relation keys, and indexes for tenant-safe querying.
  **Must NOT do**: Do not omit `OrganizationID` on tenant-owned tables.

  **Recommended Agent Profile**:
  - Category: `implementation` — Reason: schema contract implementation.
  - Skills: [`ocs-delegation-gate`] — ensures consistency and guardrails.
  - Omitted: [`ocs-runtime-validation`] — runtime check occurs later.

  **Parallelization**: Can Parallel: YES | Wave 2 | Blocks: T9,T10,T11,T12,T13,T19 | Blocked By: T6

  **References**:
  - Requirement anchor: user-specified entities and `organization_id` lock.
  - Pattern: `https://gorm.io/docs/models.html` — model tags/basics.

  **Acceptance Criteria**:
  - [ ] Model package compiles and migrations generate expected columns.
  - [ ] Tenant tables include indexed `organization_id`.
  - [ ] JSON serialization matches API contract naming.

  **QA Scenarios**:
  ```
  Scenario: Models migrate and relation constraints valid
    Tool: Bash
    Steps: Run migration + schema inspection test
    Expected: All five entities created with expected FKs/indexes
    Evidence: .sisyphus/evidence/task-7-models.txt

  Scenario: Missing organization_id guard
    Tool: Bash
    Steps: Run unit test asserting tenant models include `OrganizationID` and tags
    Expected: Test fails if any model lacks required field/tag
    Evidence: .sisyphus/evidence/task-7-models-error.txt
  ```

  **Commit**: YES | Message: `feat(models): add multi-tenant gorm entities` | Files: [`backend/internal/models/**`]

- [ ] 8. Author backend architecture artifacts (Mermaid flowchart + ERD)

  **What to do**: Create Mermaid flowchart for end-to-end chat flow (human_assigned vs bot_handling, RAG by org, prompt assembly, WebSocket response) and Mermaid ERD for requested entities.
  **Must NOT do**: Do not publish diagram without explicit tenant isolation paths.

  **Recommended Agent Profile**:
  - Category: `writing` — Reason: architecture artifact precision.
  - Skills: [`ocs-delegation-gate`] — output quality guard.
  - Omitted: [`ocs-markdown-autofix`] — optional unless markdown lint is enabled.

  **Parallelization**: Can Parallel: YES | Wave 1 | Blocks: T20 | Blocked By: T4,T6,T7

  **References**:
  - External: `https://supabase.com/docs/guides/ai/rag-with-permissions` — RAG isolation guidance.
  - External: `https://docs.pinecone.io/guides/index-data/implement-multitenancy` — namespace model for future extension.

  **Acceptance Criteria**:
  - [ ] `backend/docs/architecture/system-flowchart.mmd` exists and renders.
  - [ ] `backend/docs/architecture/erd.mmd` exists and renders.
  - [ ] Diagrams include all requested entities/flow branches.

  **QA Scenarios**:
  ```
  Scenario: Mermaid diagrams render successfully
    Tool: Bash
    Steps: Run mermaid CLI render for both `.mmd` files
    Expected: Rendered outputs generated without parse errors
    Evidence: .sisyphus/evidence/task-8-mermaid.txt

  Scenario: Missing node/entity validation
    Tool: Bash
    Steps: Run diagram lint script that checks required node/entity names
    Expected: Script fails if required labels are absent
    Evidence: .sisyphus/evidence/task-8-mermaid-error.txt
  ```

  **Commit**: YES | Message: `docs(backend): add system flowchart and tenant erd` | Files: [`backend/docs/architecture/**`]

<!-- TASKS_INSERT_HERE -->
- [ ] 1. Bootstrap monorepo structure and migrate template into `/frontend`

  **What to do**: Create root-level `/frontend` and `/backend`; copy/adapt current app from `template/free-react-tailwind-admin-dashboard-main` into `/frontend` with no behavior change; preserve TailAdmin structure and assets.
  **Must NOT do**: Do not alter UI behavior/functionality in this migration step.

  **Recommended Agent Profile**:
  - Category: `implementation` — Reason: filesystem/project-structure bootstrap.
  - Skills: [`ocs-delegation-gate`] — enforce safe delegation and evidence discipline.
  - Omitted: [`ocs-runtime-validation`] — not yet runtime-proof stage.

  **Parallelization**: Can Parallel: NO | Wave 1 | Blocks: T2,T3,T4,T14-T18 | Blocked By: -

  **References**:
  - Pattern: `template/free-react-tailwind-admin-dashboard-main/package.json` — current frontend manifest baseline.
  - Pattern: `template/free-react-tailwind-admin-dashboard-main/src/layout/AppLayout.tsx` — layout shell baseline.

  **Acceptance Criteria**:
  - [ ] `frontend/package.json` exists and preserves scripts (`dev`,`build`,`lint`,`preview`).
  - [ ] `npm --prefix "frontend" run build` succeeds.

  **QA Scenarios**:
  ```
  Scenario: Frontend migrated and buildable
    Tool: Bash
    Steps: Run `npm --prefix "frontend" install` then `npm --prefix "frontend" run build`
    Expected: Exit code 0 for install and build
    Evidence: .sisyphus/evidence/task-1-monorepo-migration.txt

  Scenario: Missing source path protection
    Tool: Bash
    Steps: Validate `template/free-react-tailwind-admin-dashboard-main` still exists and compare against `/frontend` copy structure
    Expected: Source remains intact; `/frontend/src` present
    Evidence: .sisyphus/evidence/task-1-monorepo-migration-error.txt
  ```

  **Commit**: YES | Message: `chore(repo): bootstrap frontend and backend roots` | Files: [`frontend/**`, `backend/.gitkeep`]

- [ ] 2. Define workspace scripts and environment/secrets policy

  **What to do**: Add root scripts (frontend/backend test/build/lint orchestration), create `.env.example` contracts for frontend/backend, and update ignore rules to protect `.env*` secrets.
  **Must NOT do**: Do not commit real secrets.

  **Recommended Agent Profile**:
  - Category: `security` — Reason: secret-handling and operational guardrails.
  - Skills: [`ocs-delegation-gate`] — required for safe policy delegation.
  - Omitted: [`ocs-release-integrity`] — release flow not yet active.

  **Parallelization**: Can Parallel: YES | Wave 1 | Blocks: T5,T6,T9-T13,T19 | Blocked By: T1

  **References**:
  - Pattern: `template/free-react-tailwind-admin-dashboard-main/.gitignore` — baseline ignore rules requiring extension for `.env` handling.
  - Pattern: `template/free-react-tailwind-admin-dashboard-main/README.md` — setup guidance anchor.

  **Acceptance Criteria**:
  - [ ] Root `.gitignore` blocks `.env`, `.env.*`, and keeps `.env.example` tracked.
  - [ ] `npm run lint:all` and `npm run test:all` scripts defined at repo root.

  **QA Scenarios**:
  ```
  Scenario: Env policy enforced
    Tool: Bash
    Steps: Create temp `.env.local`, run `git status --short`
    Expected: `.env.local` not listed as tracked file
    Evidence: .sisyphus/evidence/task-2-env-policy.txt

  Scenario: Secret leak guard failure path
    Tool: Bash
    Steps: Temporarily add fake key in tracked file lint step, run secret scan command configured in task
    Expected: Scan reports violation and exits non-zero
    Evidence: .sisyphus/evidence/task-2-env-policy-error.txt
  ```

  **Commit**: YES | Message: `chore(repo): add workspace scripts and env policy` | Files: [`.gitignore`, `.env.example`, `package.json`]

- [ ] 3. Bootstrap frontend test stack (Vitest + RTL + Playwright smoke)

  **What to do**: Configure Vitest and React Testing Library for component tests and Playwright smoke for route-level checks; add baseline tests for app shell and nav rendering.
  **Must NOT do**: Do not gate on flaky snapshot-only tests.

  **Recommended Agent Profile**:
  - Category: `testing` — Reason: fresh test infrastructure design/implementation.
  - Skills: [`ocs-delegation-gate`] — keeps test setup scoped and verifiable.
  - Omitted: [`ocs-runtime-validation`] — runtime validation comes post-feature integration.

  **Parallelization**: Can Parallel: YES | Wave 1 | Blocks: T16,T17,T18,T19 | Blocked By: T1

  **References**:
  - Pattern: `frontend/src/App.tsx` — route shell baseline for smoke tests.
  - Pattern: `frontend/src/layout/AppSidebar.tsx` — navigation render assertions.

  **Acceptance Criteria**:
  - [ ] `npm --prefix "frontend" run test` exits 0.
  - [ ] Playwright smoke verifies `/`, `/signin`, and sidebar navigation render.

  **QA Scenarios**:
  ```
  Scenario: Frontend test baseline passes
    Tool: Bash
    Steps: Run `npm --prefix "frontend" run test` and `npm --prefix "frontend" run test:e2e:smoke`
    Expected: Both commands exit 0
    Evidence: .sisyphus/evidence/task-3-frontend-tests.txt

  Scenario: Regression catch on broken route
    Tool: Bash
    Steps: Intentionally point one smoke test to non-existing path and run smoke suite
    Expected: Smoke suite fails with route assertion error
    Evidence: .sisyphus/evidence/task-3-frontend-tests-error.txt
  ```

  **Commit**: YES | Message: `test(frontend): bootstrap vitest rtl and playwright smoke` | Files: [`frontend/vitest.config.*`, `frontend/tests/**`]

- [ ] 4. Initialize backend Gin service skeleton with test harness

  **What to do**: Create `/backend` Go module, `/cmd/api/main.go`, `/internal` package skeleton, health endpoint, and baseline `httptest` for health route.
  **Must NOT do**: Do not embed business logic in `cmd`; keep it wiring-only.

  **Recommended Agent Profile**:
  - Category: `implementation` — Reason: backend bootstrap and module scaffolding.
  - Skills: [`ocs-delegation-gate`] — preserve structure and verification rigor.
  - Omitted: [`ocs-openai-multi-account`] — unrelated at this stage.

  **Parallelization**: Can Parallel: YES | Wave 1 | Blocks: T5,T6,T7,T9-T13,T19 | Blocked By: T1

  **References**:
  - External: `https://gin-gonic.com/en/docs/introduction/` — framework baseline.
  - External: `https://github.com/gin-gonic/gin/blob/d3ffc9985281dcf4d3bef604cce4e662b1a327a6/gin.go#L243-L250` — net/http compatibility anchor.

  **Acceptance Criteria**:
  - [ ] `backend/go.mod` exists with Gin dependency.
  - [ ] `go test ./...` from `backend` exits 0.
  - [ ] `GET /health` returns 200 JSON.

  **QA Scenarios**:
  ```
  Scenario: Backend starts and serves health check
    Tool: Bash
    Steps: Run `go run ./cmd/api` then `curl -i http://localhost:8080/health`
    Expected: HTTP 200 with expected JSON body
    Evidence: .sisyphus/evidence/task-4-backend-skeleton.txt

  Scenario: Invalid route handling
    Tool: Bash
    Steps: Request `GET /unknown`
    Expected: HTTP 404 with structured error payload
    Evidence: .sisyphus/evidence/task-4-backend-skeleton-error.txt
  ```

  **Commit**: YES | Message: `feat(backend): scaffold gin api with health endpoint` | Files: [`backend/cmd/**`, `backend/internal/**`, `backend/go.mod`]

## Final Verification Wave (MANDATORY — after ALL implementation tasks)
> 4 review agents run in PARALLEL. ALL must APPROVE. Present consolidated results to user and get explicit "okay" before completing.
> **Do NOT auto-proceed after verification. Wait for user's explicit approval before marking work complete.**
> **Never mark F1-F4 as checked before getting user's okay.**
- [ ] F1. Plan Compliance Audit — oracle
- [ ] F2. Code Quality Review — unspecified-high
- [ ] F3. Real Manual QA — unspecified-high (+ playwright if UI)
- [ ] F4. Scope Fidelity Check — deep

## Commit Strategy
- Gunakan atomic commits per task (atau pair task yang tightly-coupled), mengikuti urutan dependency matrix.
- Format pesan: `type(scope): description`.
- Semua commit wajib menyertakan bukti verifikasi command di task evidence.
- Dilarang squash lintas wave agar rollback terkontrol.

## Success Criteria
- Semua deliverables backend/frontend tercapai sesuai request awal tanpa gap fungsional.
- Tidak ada keputusan arsitektur tersisa untuk executor.
- Semua acceptance criteria task dan final verification wave lulus.
- User memberikan explicit "okay" setelah hasil final verification dipresentasikan.
