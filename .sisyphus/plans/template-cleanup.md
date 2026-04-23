# Work Plan: Template Cleanup & Refactoring

## 1. Goal
Membersihkan kode bawaan (showcase/junk) dari template TailAdmin React, dengan mempertahankan komponen inti yang dibutuhkan untuk aplikasi Chatbot. Menghilangkan bloated routing dan sidebar items, sambil menjaga kestabilan build (TypeScript & Vite).

## 2. Scope & Guardrails
- **IN SCOPE**: Menghapus rute dan file halaman showcase (Forms, Tables, Charts, UI Elements, Blank, Calendar). Menghapus item menu dari Sidebar. Mengganti label "Ecommerce" menjadi "Home".
- **OUT OF SCOPE**: Membuat fitur baru. Menghapus folder `src/components/ui/`, `src/components/form/`, atau direktori dasar komponen lainnya.
- **GUARDRAIL (Metis)**: DILARANG KERAS menghapus isi dari `frontend/src/components/`. Semua penghapusan fisik hanya difokuskan pada `frontend/src/pages/`.
- **VERIFIKASI MUTLAK**: Harus memastikan `npm run build` dan `tsc --noEmit` sukses (Exit code 0).

## 3. Tasks

### Task 1: Delete Showcase Page Directories
**File/Dir to Delete**:
- `frontend/src/pages/UiElements/`
- `frontend/src/pages/Chart/`
- `frontend/src/pages/Form/`
- `frontend/src/pages/Tables/`
- `frontend/src/pages/Blank.tsx` (jika tidak dalam folder)
- `frontend/src/pages/Calendar.tsx` (sesuai rekomendasi Metis, tidak relevan dengan Chatbot)

### Task 2: Refactor Routing (`frontend/src/App.tsx`)
**Modifications**:
- Hapus semua statement `import` yang merujuk ke halaman-halaman yang dihapus di Task 1.
- Hapus komponen `<Route />` untuk path: `/calendar`, `/profile` (jika tidak dipakai, tapi sesuai persetujuan kita simpan profil, jadi hapus yang berikut:), `/form-elements`, `/basic-tables`, `/alerts`, `/avatars`, `/badge`, `/buttons`, `/images`, `/videos`, `/line-chart`, `/bar-chart`, `/blank`.

### Task 3: Refactor Sidebar (`frontend/src/layout/AppSidebar.tsx`)
**Modifications**:
- Pada array `navItems`:
  - Ubah label `Ecommerce` menjadi `Home` (Rekomendasi Metis).
  - Hapus item `Calendar`.
  - Hapus sub-menu `Forms`, `Tables`, dan `Pages` (termasuk link 404 Error di dalamnya).
- Pada array `othersItems`:
  - Hapus sub-menu `Charts`.
  - Hapus sub-menu `UI Elements`.
- Pertahankan `Authentication` dan `Hubungi Saya`.

## 4. Final Verification Wave
- [ ] Jalankan `cd frontend && npx tsc --noEmit` untuk memastikan tidak ada *import* yang *error/dangling*.
- [ ] Jalankan `cd frontend && npm run build` untuk memverifikasi proses build Vite sukses 100%.
- [ ] Jalankan pengecekan file untuk memastikan `frontend/src/components/` tidak tersentuh.
