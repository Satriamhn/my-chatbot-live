# Work Plan: SaaS AI Chatbot Platform Multi-Tenant

## 1. Goal
Membangun arsitektur backend GORM dan merombak template React+Tailwind (`/frontend`) menjadi UI pixel-perfect untuk Dashboard Chatbot.

## 2. Scope
**IN SCOPE:**
- [Backend] Pembuatan `models.go` untuk entitas multi-tenant: Organizations, Users, Bot_Settings, Knowledge_Bases, Chat_Sessions, Messages.
- [Backend] Dokumentasi `architecture.md` (Flowchart & ERD).
- [Frontend] Instalasi Axios dan setup interceptor `api.ts`.
- [Frontend] Pembuatan halaman `/signin` & `/signup` (SaaS Authentication UI).
- [Frontend] Pembuatan halaman `/settings/bot` (Bot Config & Live Preview dengan Mock Data statis).
- [Frontend] Pembuatan halaman `/knowledge` (Knowledge Base & Data Table dengan Mock Data).
- [Frontend] Pembuatan halaman `/inbox` (Omnichannel Chat Window & Human Takeover dengan Mock Data).
- [Frontend] Tambah menu sidebar "Hubungi Saya (Owner)".

**OUT OF SCOPE:**
- Implementasi API endpoint lengkap (Controller/Route) di Golang (hanya GORM Models & Docs yang diminta).
- Backend business logic untuk RAG.

## 3. Technical Approach

### Backend
- Buat folder `backend/internal/models/` dan `backend/docs/`.
- Gunakan `google/uuid` dan `gorm.io/gorm`.
- Tambahkan struct `User` untuk autentikasi dan relasi multi-tenant.
- Terapkan Soft Deletes (`gorm.DeletedAt`) dan Composite Index `(org_id, id)` untuk keamanan multi-tenant yang ketat.
- Tulis ERD dan Flowchart format Mermaid.js ke dalam file `backend/docs/architecture.md`.

### Frontend
- Direktori kerja: `frontend/`
- Install: `npm install axios react-router-dom lucide-react`
- Setup API: `frontend/src/services/api.ts` dengan header `Authorization` dan `X-Org-Id` (Gunakan static mock ID untuk sekarang).
- Gunakan komponen auth bawaan template (Sign In / Sign Up) dan modifikasi copy/teks-nya.
- Gunakan struktur UI `AppLayout.tsx` dan `AppSidebar.tsx` bawaan template untuk menu dashboard.
- **PENTING**: Gunakan *Mock Data* JSON statis untuk mengisi tabel dan chat window.

## 4. Final Verification Wave
- [ ] Pastikan model Go (termasuk User) dapat di-compile.
- [ ] Pastikan frontend lolos TypeScript check.
- [ ] Pastikan ERD dan Flowchart valid saat di-render.
- [ ] Pastikan routing /signin, /signup, /settings/bot, /knowledge, dan /inbox berjalan sempurna.

[DECISION NEEDED] Tidak ada. Rencana sudah final.
