# 🏗️ Public Calls Portal — Executive & Technical Report

> **Project:** Public Calls Portal — SECOP II Convocatorias Públicas
> **Repository:** `retojikko_sebaspaniagua/challenge_dev`
> **Branch:** `feature/challenge_dev_public_proposals`
> **Tech Stack:** Go 1.26.4 · PostgreSQL 16 · HTML5/CSS/JS · REST/JSON
> **External API:** datos.gov.co — SECOP II Socrata Open Data API (SODA v2.1)
> **Last Updated:** 2026-07-03

---

## Executive Summary

The **Public Calls Portal** is a web application that connects to Colombia's SECOP II public procurement platform via `datos.gov.co` to allow users to **browse, filter, and save public calls for proposals** (convocatorias públicas).

The project follows a **backend-first, two-phase approach**:
1. **Phase 1 (Complete):** PostgreSQL database schema + Go backend API skeleton
2. **Phase 2 (Pending):** HTML5/CSS/JS frontend

The backend is fully implemented with a **layered architecture** (domain → repository → service → handler → router), **22 Go source files**, **29 passing tests**, and a fully deployed PostgreSQL database (`portal_plan_public_app`).

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend (Phase 2)                        │
│            HTML5 · CSS · Vanilla JavaScript                  │
├─────────────────────────────────────────────────────────────┤
│                     REST / JSON                              │
├─────────────────────────────────────────────────────────────┤
│                        Router                                │
│               http.ServeMux (net/http)                       │
├───────────┬──────────┬──────────────────────────────────────┤
│  Handler  │  Handler  │            Handler                   │
│   Auth    │ Proposals│         Saved Proposals               │
├───────────┴──────────┴──────────────────────────────────────┤
│                        Service                               │
│    Auth · Proposal Listing/Filter · Saved Proposal CRUD      │
├─────────────────────────────────────────────────────────────┤
│                       Repository                             │
│       User · PublicCallProposal · SavedProposal              │
├─────────────────────────────────────────────────────────────┤
│                        Domain                                │
│        Structs: User · PublicCallProposal · SavedProposal    │
├─────────────────────────────────────────────────────────────┤
│                      Database (pgx)                          │
│         portal_plan_public_app (PostgreSQL 16)               │
└─────────────────────────────────────────────────────────────┘
```

---

## Completed Tasks

### ✅ Task 0: Discover datos.gov.co API

| Item | Detail |
|------|--------|
| **Endpoint** | `GET https://datos.gov.co/resource/p6dx-8zbt.json` |
| **Dataset** | SECOP II — Procesos de Contratación (Dataset ID: `p6dx-8zbt`) |
| **Protocol** | Socrata Open Data API (SODA) v2.1 |
| **Auth** | None required (public data) |
| **Documentation** | `integration_doc/datosgovco/documentation.md` |
| **Test PoC** | `integration_doc/datosgovco/test.go` |

Supports full-text search (`$q`), column filtering (`$where`), pagination (`$limit`, `$offset`), and sorting (`$order`).

### ✅ Task 1: Database Schema & PostgreSQL Setup

- **PostgreSQL 16** installed on local machine (port 5432)
- Database `portal_plan_public_app` created
- User `sebasdb` created with `CREATEDB` + full CRUD on all schema objects
- **166-line SQL schema** at `database/public_calls_database.sql`
- Schema includes: `pgcrypto` extension, `set_updated_at()` auto-trigger function, 3 tables, indexes, foreign keys, comments

**Tables:**

| Table | Purpose | Key Columns |
|-------|---------|-------------|
| `users` | Portal user accounts | id, email, password, status, timestamps |
| `public_calls_proposals` | Cached SECOP II calls (50+ fields) | id, nombre, entidad, fase, modalidad, precio, fechas, estado, proveedores |
| `public_call_user_associations` | User ↔ Saved call links | id (uuid), public_call_id, user_id, association_date |

### ✅ Task 2: Go Backend Skeleton

**22 Go source files** across 7 packages:

```
challenge_dev/backend/
├── cmd/api/main.go                        # Entrypoint, DI wiring, server startup
├── go.mod / go.sum                        # Module: github.com/tianpl92/.../backend (go 1.26.4)
├── internal/
│   ├── config/
│   │   ├── config.go                      # Loads PORT, DB_USER, DB_PASSWORD, SECRET_KEY from env
│   │   └── config_test.go                 # 2 tests (defaults + env override)
│   ├── database/
│   │   └── postgres.go                    # pgx v5 connection pool (ready for prod swap)
│   ├── domain/
│   │   ├── user.go                        # User struct with JSON tags
│   │   ├── public_proposal.go             # PublicCallProposal struct (50+ fields mapped)
│   │   └── saved_proposal.go              # SavedProposal struct
│   ├── repository/
│   │   ├── user_repository.go             # Interface + mock: FindByEmail, Create
│   │   ├── proposal_repository.go         # Interface + mock: List, Filter, FindByID
│   │   ├── saved_proposal_repository.go   # Interface + mock: Save, FindByUserID, Delete
│   │   ├── mock_repository.go             # In-memory mock implementations
│   │   └── mock_repository_test.go        # 8 tests (CRUD edge cases)
│   ├── service/
│   │   ├── auth_service.go                # Login with HMAC token, ValidateToken
│   │   ├── proposal_service.go            # ListProposals, FilterProposals
│   │   ├── saved_proposal_service.go      # SaveProposal, GetSaved, RemoveSaved
│   │   └── service_test.go                # 10 tests (auth, proposals, saved)
│   ├── handler/
│   │   ├── health_handler.go              # GET /health → {"status":"ok"}
│   │   ├── auth_handler.go                # POST /login → {token, user}
│   │   ├── proposal_handler.go            # GET /public-proposals?query=&fase=&entidad=
│   │   ├── saved_proposal_handler.go      # POST + GET /saved-proposals (auth required)
│   │   └── handler_test.go                # 9 tests (httptest.Server)
│   └── router/
│       └── router.go                      # http.ServeMux wiring for all 5 routes
└── README.md                              # Env vars documentation
```

### ✅ API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/health` | No | Health check |
| `POST` | `/login` | No | Authenticate with email + password → returns HMAC token |
| `GET` | `/public-proposals` | No | List proposals. Filters: `?query=&fase=&entidad=` |
| `POST` | `/saved-proposals` | Yes | Save a proposal (`Authorization: Bearer <token>`) |
| `GET` | `/saved-proposals` | Yes | List user's saved proposals |

### ✅ Tests & Quality Gates

**29/29 tests passing** across 4 packages:

| Package | Tests | Coverage |
|---------|-------|----------|
| `internal/config` | 2 | Defaults + env override |
| `internal/handler` | 9 | Health, login (success/fail/bad-req), proposals (list/filter), saved (auth/create/list/invalid-token) |
| `internal/repository` | 8 | Mock CRUD for users, proposals, saved proposals (create, duplicate, delete, not-found) |
| `internal/service` | 10 | Auth (login, validate, invalid pw, not-found), proposals (list, filter), saved (save, get, remove) |

**Toolchain checks:**
- `gofmt` — 0 files need formatting ✅
- `go vet` — zero warnings ✅
- `go build ./...` — compiles cleanly ✅

### ✅ Environment & Config

| File | Purpose |
|------|---------|
| `backend/.env` (gitignored) | `DB_USER=sebasdb`, `DB_PASSWORD=***` |
| `.gitignore` | Excludes `backend/.env`, `.hermes/` |
| `backend/README.md` | Documents required env vars |
| `SOUL.md` | This report (replaces previous notes) |

---

## Milestone Progress

| # | Milestone | Status |
|---|-----------|--------|
| 0 | Scan datos.gov.co API | ✅ Done |
| 1 | Backend project skeleton + database schema | ✅ Done |
| 2 | Core backend endpoints + business logic | ✅ Done |
| 3 | Backend unit tests (TDD) | ✅ Done |
| 4 | Frontend pages | ⏳ Pending |
| 5 | Frontend ↔ Backend wiring | ⏳ Pending |
| 6 | End-to-end validation | ⏳ Pending |

---

## Git History

```
d68b6b1 feat: add backend skeleton with layered architecture
95256fb Refactor sql schema
18d3ba1 feat (reto jikko) database schema v2 - secop ii columns comparison
4c84bee add git ignore files
a6035d9 feat (reto jikko) challenge dev - database schema creation
cb9c4ee Creacion de specs, resultados y analysis.md usando hermes
a318cb9 Initial commit
```

---

## Project File Tree

```
challenge_dev/
├── .gitignore
├── hermes.md                          # Product guidance
├── SOUL.md                            # This report
├── backend/                           # Go backend (Phase 1)
│   ├── cmd/api/main.go
│   ├── go.mod / go.sum
│   ├── internal/
│   │   ├── config/
│   │   ├── database/
│   │   ├── domain/
│   │   ├── repository/
│   │   ├── service/
│   │   ├── handler/
│   │   └── router/
│   └── README.md
├── database/
│   ├── public_calls_database.sql      # Production schema
│   └── compare_rubric/                # Schema comparison artifacts
├── integration_doc/
│   └── datosgovco/                    # SECOP II API docs + test
├── specs_project/                     # Task specs (1–7)
└── .hermes/plans/                     # Implementation plan
```

---

## Next Steps (Phase 2 — Frontend)

Per the implementation plan, the remaining work is:

1. **Build frontend pages** — login.html, dashboard.html with CSS
2. **Wire frontend to backend** — `fetch()` wrappers in JavaScript
3. **End-to-end validation** — full login → browse → filter → save → retrieve flow
4. **Connect real PostgreSQL** — swap in-memory mocks in `main.go` for actual pgx repository implementations
5. **Integrate datos.gov.co** — fetch live SECOP II data via SODA API and cache in `public_calls_proposals` table

---

## Risks & Open Questions

| Risk | Mitigation |
|------|------------|
| datos.gov.co API shape may differ from UI needs | Normalization layer isolated in repository |
| Auth is minimal (HMAC tokens) | Can be upgraded to real JWT or OAuth later |
| No real DB connection in main.go yet | `database/postgres.go` has `NewPool()` ready; swap mocks once DB integration needed |
| Frontend is static (no build tool) | Fine for MVP; can add framework later if needed |