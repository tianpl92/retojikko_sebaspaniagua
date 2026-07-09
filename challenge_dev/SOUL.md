# 🏗️ Public Calls Portal — Executive & Technical Report

> **Project:** Public Calls Portal — SECOP II Convocatorias Públicas
> **Repository:** `retojikko_sebaspaniagua/challenge_dev`
> **Branch:** `feature/challenge_dev_public_proposals`
> **Tech Stack:** Go 1.26.4 · PostgreSQL 16 · HTML5/CSS/JS · REST/JSON
> **External API:** datos.gov.co — SECOP II Socrata Open Data API (SODA v2.1)
> **Auth:** JWT (HMAC HS256, 1h expiry) · bcrypt password hashing
> **Last Updated:** 2026-07-09

---

## Executive Summary

The **Public Calls Portal** connects to Colombia's SECOP II public procurement platform via `datos.gov.co`, enabling users to **browse, filter, save, and review public calls for proposals** (convocatorias públicas).

The project was delivered in **three development phases**, each following a spec-driven workflow:

| Phase | Scope | Status |
|-------|-------|--------|
| **Database** | Schema design (4 tables), migration SQL, FK constraints, indexes | ✅ Complete |
| **Backend** | Layered Go service: 9 REST endpoints, JWT auth, datos.gov.co integration, mocks + Postgres repos | ✅ Complete |
| **Frontend** | 4 HTML/CSS/JS pages, fetch-based API client, responsive UI, modals | ✅ Complete |
| **Bugfix/Refinement** | Specs 10–19: endpoint refinements, response shape fixes, frontend integration fixes | ✅ Complete |

---

## Workflow Applied

Every phase followed the same **spec-driven, subagent-delegation pattern**:

```
  ┌─────────────────────────────────────────────────────────┐
  │  1. Spec document (numbered .md in specs_project/)      │
  │     - Context, Role, Requirements, Actions, Expected     │
  ├─────────────────────────────────────────────────────────┤
  │  2. Plan → Subagent delegation (Hermes)                 │
  │     - Principal: deepseek/deepseek-v4-flash (orchestrate)│
  │     - Subagent:   deepseek/deepseek-v4-pro (codegen)     │
  ├─────────────────────────────────────────────────────────┤
  │  3. Build → Code generation + manual review             │
  ├─────────────────────────────────────────────────────────┤
  │  4. Verify → go build, go vet, go test, gofmt           │
  ├─────────────────────────────────────────────────────────┤
  │  5. Commit → GitHub via SSH                              │
  └─────────────────────────────────────────────────────────┘
```

---

## Phase 1 — Database Schema Design (Specs 1–6)

| Spec | Description | Outcome |
|------|-------------|---------|
| 1 | Scan datos.gov.co for Convocatorias endpoint | Identified SECOP II SODA endpoint `p6dx-8zbt.json` |
| 2–5 | Schema iterations | 4 tables with FK constraints, indexes, triggers, comments |
| 6 | Migration execution | `database/public_calls_database.sql` — production-ready |

### Schema: `portal_plan_public_app`

| Table | Purpose | Key Columns |
|-------|---------|-------------|
| `users` | Portal accounts | `id` (document#), `email` (unique), `password` (bcrypt) |
| `public_calls_proposals` | SECOP II proposal cache | `id` (PK), 30+ datos.gov.co columns |
| `public_call_user_associations` | User ↔ saved proposal | FK to both, unique constraint on pair |
| `user_sessions` | JWT session tokens | FK to user, token, expires_at |

---

## Phase 2 — Backend Service (Specs 7–9)

Layered architecture with strict separation:

```
Router (http.ServeMux)
  → AuthMiddleware (JWT Bearer validation)
    → Handler (HTTP concern: parse, validate, respond)
      → Service (business logic)
        → Repository (data access interface)
          → MockRepo | PostgresRepo
            → Domain (structs)
```

### Available Endpoints

| # | Method | Path | Auth | Description |
|---|--------|------|:----:|-------------|
| 1 | `GET` | `/health` | ❌ | Health check |
| 2 | `POST` | `/user-create` | ❌ | Register (bcrypt, dedup by doc# + email) |
| 3 | `POST` | `/login` | ❌ | Authenticate → JWT (1h) |
| 4 | `POST` | `/user-info` | ✅ | Profile (no password in response) |
| 5 | `POST` | `/user-modify` | ✅ | Update allowed fields only |
| 6 | `GET` | `/public-proposals` | ✅ | Live data from datos.gov.co (SoQL) |
| 7 | `GET` | `/saved_proposals` | ✅ | Full proposal data shape |
| 8 | `POST` | `/saved-proposals` | ✅ | Save association (idempotent) |
| 9 | `POST` | `/proposal-save` | ✅ | Upsert full proposal locally |

### Key Design Decisions

| Decision | Rationale |
|----------|-----------|
| **Mocks in `main.go`** | Fast iteration without PostgreSQL dependency; Postgres repos compile-ready |
| **Datos.gov.co passthrough** | Return raw JSON from external API — avoids coupling to its schema |
| **`urlproceso` as `interface{}`** | datos.gov.co returns `{"url":"..."}` (object), not a plain string |
| **Full-data response shape** | `GET /saved_proposals` returns 30+ proposal fields via LEFT JOIN, not association metadata |

---

## Phase 3 — Frontend Application (Specs 11–15)

Static HTML5/CSS/JS served via `python3 -m http.server` on port 3000.

| Page | File | Function |
|------|------|----------|
| Login | `pages/login.html` + `login.js` | JWT authentication, localStorage |
| Dashboard | `pages/dashboard.html` + `dashboard.js` | Browse, filter, paginate, save proposals |
| Saved Proposals | `pages/saved-proposals.html` + `saved-proposals.js` | List saved, detail modal with all fields |
| User Create | `pages/user-create.html` + `user-create.js` | Registration form |

### Frontend → Backend Flow

```
Dashboard:
  1. Select proposal(s) → click "Guardar"
  2. POST /proposal-save  { ...proposal, id: id_del_proceso }
  3. POST /saved-proposals { public_call_id }
  4. Loop for each selected proposal

Saved Proposals:
  1. On page load → GET /saved_proposals
  2. Render response (full proposal data shape)
  3. Click "Detalles..." → modal with 30 fields
```

---

## Phase 4 — Bugfix & Refinement (Specs 10–19)

| Spec | Fix |
|------|-----|
| 10–17 | Endpoint refinements, integration behavior fixes, response shape corrections |
| **18** | `SavedProposalFullData` response struct: removes `public_call_id`, `association_date`; adds `id_del_proceso`, 30 full-data fields. Adds `FindByUserIDWithProposals` (JOIN query) and `Upsert` to Postgres repos |
| **19** | Frontend integration: `id_del_proceso` → `id` mapping in dashboard save payload; `urlproceso` object handling in backend; stale modal fields removed; detail modal shows all 30 proposal fields |

---

## Tests & Quality Gates

**37 tests passing** across all packages:

| Package | Tests | Scope |
|---------|-------|-------|
| `internal/config` | 2 | Defaults + env override |
| `internal/handler` | 16 | All endpoints + auth middleware + full-data response shape |
| `internal/repository` | 7 | Mock CRUD: users, proposals, saved proposals |
| `internal/service` | 12 | JWT auth (login, validate, create), proposals, saved proposals |

**Toolchain:**
- `gofmt` — clean (0 files need formatting) ✅
- `go vet ./...` — zero warnings ✅
- `go build ./...` — compiles ✅

---

## Source Files

```
challenge_dev/
├── .hermes/plans/                  # Implementation plan (updated)
├── database/
│   ├── public_calls_database.sql   # Production schema
│   └── compare_rubric/             # Schema comparison artifacts
├── backend/
│   ├── cmd/api/main.go             # Entrypoint (mock repos)
│   ├── go.mod / go.sum
│   └── internal/
│       ├── config/config.go        # Env vars
│       ├── database/postgres.go    # pgx pool
│       ├── domain/                 # 4 domain structs
│       ├── repository/             # Interfaces + mocks + Postgres impls (5 files)
│       ├── service/                # 4 services + tests
│       ├── handler/                # 7 handlers + middleware + tests
│       └── router/router.go       # Route wiring + CORS
├── frontend/
│   ├── pages/                     # login, dashboard, saved-proposals, user-create
│   ├── scripts/                   # 4 JS files
│   └── sheets/styles.css          # Shared styles, palette: #4A4466, #6EADBC, #9FCBAD, #F1F7D4
├── specs_project/                 # Specs 1–19
├── apicontract.md                 # v1.1.0
├── SOUL.md                        # This report
└── hermes.md                      # Product guidance
```

---

## Git Workflow

```
Plan (spec .md) → Subagent delegation → Code generation
  → go build + go vet + go test + gofmt → Commit via SSH
```

Do not commit changes until explicitly requested ("Do not commit changes yet").

---

## Milestone Progress

| # | Milestone | Status |
|---|-----------|--------|
| 0 | Scan datos.gov.co API | ✅ |
| 1 | Database schema (4 tables) | ✅ |
| 2 | Backend skeleton + layered architecture | ✅ |
| 3 | JWT auth + user CRUD endpoints | ✅ |
| 4 | datos.gov.co live integration | ✅ |
| 5 | Saved proposals (associate, list, upsert) | ✅ |
| 6 | Frontend pages (4 pages) | ✅ |
| 7 | Frontend ↔ Backend wiring | ✅ |
| 8 | End-to-end validation | ✅ |
| 9 | Bugfix & refinement (specs 10–19) | ✅ |

---

## Open Questions

- Wire Postgres repos into `main.go` for production (currently using mocks)
- Add data sync/cache strategy for datos.gov.co (currently live-fetch per request)
- Add refresh token mechanism (currently 1h JWT, re-login required)
- CI/CD pipeline for automated deployment