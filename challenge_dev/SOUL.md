# 🏗️ Public Calls Portal — Executive & Technical Report

> **Project:** Public Calls Portal — SECOP II Convocatorias Públicas
> **Repository:** `retojikko_sebaspaniagua/challenge_dev`
> **Branch:** `feature/challenge_dev_public_proposals`
> **Tech Stack:** Go 1.26.4 · PostgreSQL 16 · HTML5/CSS/JS · REST/JSON
> **External API:** datos.gov.co — SECOP II Socrata Open Data API (SODA v2.1)
> **Auth:** JWT (HMAC HS256, 1h expiry) · bcrypt password hashing
> **Last Updated:** 2026-07-06

---

## Executive Summary

The **Public Calls Portal** is a web application that connects to Colombia's SECOP II public procurement platform via `datos.gov.co` to allow users to **browse, filter, and save public calls for proposals** (convocatorias públicas).

**Phase 1 (Backend) is complete** with:
- **JWT authentication** replacing legacy HMAC tokens
- **Live integration** with datos.gov.co SECOP II SODA API
- **8 REST endpoints** fully implemented with input validation
- **46/46 tests passing**
- **API contract** documented at `apicontract.md`
- **PostgreSQL 16** with 4 tables (`users`, `public_calls_proposals`, `public_call_user_associations`, `user_sessions`)

**Phase 2 (Frontend) is pending** — HTML5/CSS/JS frontend pages, end-to-end wiring.

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
├─────────────────────────────────────────────────────────────┤
│        AuthMiddleware (JWT Bearer Token validation)          │
├───────────┬──────────┬──────────┬───────────────────────────┤
│  Handler  │  Handler  │  Handler │     Handler              │
│   Auth    │  Users   │ Proposals│  Saved Proposals          │
├───────────┴──────────┴──────────┴───────────────────────────┤
│                        Service                               │
│  Auth(JWT) · User CRUD · DatosGov(ext API) · Saved Proposal │
├─────────────────────────────────────────────────────────────┤
│                       Repository                             │
│       User · Proposal · SavedProposal · Session (mock)       │
├─────────────────────────────────────────────────────────────┤
│                        Domain                                │
│     Structs: User · PublicCallProposal · SavedProposal ·     │
│               UserSession                                    │
├─────────────────────────────────────────────────────────────┤
│                   PostgreSQL 16                              │
│         portal_plan_public_app (4 tables)                    │
└─────────────────────────────────────────────────────────────┘
```

---

## Completed Tasks

### ✅ Spec 8: JWT Authentication & User Endpoints

| Component | Detail |
|-----------|--------|
| **JWT Library** | `golang-jwt/jwt/v5` with HS256 signing |
| **Password hashing** | `bcrypt` via `golang.org/x/crypto` |
| **Token expiry** | 1 hour, stored in `user_sessions` table |
| **POST /login** | Validates email + bcrypt password → returns JWT |
| **POST /user-create** | Validates duplicate document + email, bcrypt-hashes password |
| **POST /user-info** | Returns authenticated user profile (no password in response) |
| **POST /user-modify** | Updates profile: only `first_name`, `last_name`, `gender`, `phone_number`, `status` |

### ✅ Spec 9: datos.gov.co Integration

| Item | Detail |
|------|--------|
| **Service** | `service/datosgov_service.go` — HTTP client for SECOP II SODA API |
| **Endpoint** | `GET https://datos.gov.co/resource/p6dx-8zbt.json` |
| **Protocol** | SoQL query language via URL params |
| **Filters** | `query` (full-text), `fase` (exact), `entidad` (partial), `limit`, `offset` |
| **Auth** | Public API (no key required), endpoint requires JWT |
| **Response** | Passthrough — returns raw JSON from datos.gov.co with original column names |
| **Env var** | `INTEGRATION_URL` in `backend/.env` |

### ✅ Spec 10: Endpoint Logic Refinements

| Endpoint | Before | After |
|----------|--------|-------|
| **POST /saved-proposals** | Returned saved object | Returns `{"message":"Guardado satisfactoriamente"}` (idempotent) |
| **GET /saved_proposals** | Returned raw records | Returns array (empty `[]` if none) |
| **POST /user-modify** | Accepted all fields | Only `first_name`, `last_name`, `gender`, `phone_number`, `status`. Rejects `email`, `password`, `id`. Requires at least one field |
| **POST /user-info** | Was `GET /user-info` | Changed to POST, no password in response, `"Usuario no registrado"` on not found |

### ✅ API Contract

- **Created** `apicontract.md` — 486 lines, 8 endpoints documented
- Each endpoint includes: protocol, input params (with examples), response format, error codes

### ✅ Database Schema

- **4 tables** in `portal_plan_public_app`:
  - `users` — Portal user accounts (id=varchar(128) as document number)
  - `public_calls_proposals` — SECOP II proposal data (50+ columns)
  - `public_call_user_associations` — User ↔ saved proposal links
  - `user_sessions` — JWT session tokens (user_id FK, token, expires_at)
- Updated `database/public_calls_database.sql` with `user_sessions` table

---

## API Endpoints

| # | Method | Path | Auth | Description |
|---|--------|------|:----:|-------------|
| 1 | `GET` | `/health` | ❌ | Health check |
| 2 | `POST` | `/user-create` | ❌ | Register new user |
| 3 | `POST` | `/login` | ❌ | Authenticate → JWT (1h) |
| 4 | `POST` | `/user-info` | ✅ | Get user profile (no password) |
| 5 | `POST` | `/user-modify` | ✅ | Update allowed profile fields |
| 6 | `GET` | `/public-proposals` | ✅ | Live data from datos.gov.co |
| 7 | `GET` | `/saved_proposals` | ✅ | List user's saved proposals |
| 8 | `POST` | `/saved-proposals` | ✅ | Save proposal (idempotent) |

---

## Tests & Quality Gates

**46/46 tests passing** across 4 packages:

| Package | Tests | Coverage |
|---------|-------|----------|
| `internal/config` | 2 | Defaults + env override |
| `internal/handler` | 22 | Health, login, user-create, user-info, user-modify, public-proposals, saved-proposals (auth, CRUD, edge cases) |
| `internal/repository` | 10 | Mock CRUD for users, proposals, saved proposals, sessions |
| `internal/service` | 12 | JWT auth (login, validate, create user), proposals, saved proposals |

**Toolchain checks:**
- `gofmt` — 0 files need formatting ✅
- `go vet` — zero warnings ✅
- `go build ./...` — compiles cleanly ✅

---

## Source Files

**27 Go files** across the backend:

```
backend/
├── cmd/api/main.go                        # Entrypoint, DI wiring
├── go.mod / go.sum                        # Go 1.26.4, deps: jwt/v5, pgx/v5, bcrypt
├── internal/
│   ├── config/config.go                   # Env vars: DB, INTEGRATION_URL, SECRET_KEY
│   ├── database/postgres.go               # pgx pool factory
│   ├── domain/
│   │   ├── user.go                        # User struct (ID=string, password hidden)
│   │   ├── public_proposal.go             # 50+ field proposal struct
│   │   ├── saved_proposal.go             # Saved proposal assoc struct
│   │   └── user_session.go               # JWT session struct
│   ├── repository/
│   │   ├── user_repository.go             # Interface: FindByEmail, Create, Update...
│   │   ├── proposal_repository.go         # Interface: List, Filter, FindByID
│   │   ├── saved_proposal_repository.go   # Interface: Save, FindByUserID, Delete
│   │   ├── session_repository.go          # Interface: Create, FindByToken, DeleteExpired
│   │   ├── mock_repository.go             # In-memory mock implementations
│   │   ├── postgres_user_repository.go    # Real pgx implementation
│   │   ├── postgres_proposal_repository.go
│   │   ├── postgres_saved_proposal_repository.go
│   │   └── postgres_session_repository.go
│   ├── service/
│   │   ├── auth_service.go                # JWT + bcrypt auth logic
│   │   ├── datosgov_service.go            # datos.gov.co HTTP client
│   │   ├── proposal_service.go            # Proposal business logic
│   │   └── saved_proposal_service.go      # Saved proposal logic
│   ├── handler/
│   │   ├── health_handler.go              # GET /health
│   │   ├── auth_handler.go                # POST /login, POST /user-create
│   │   ├── user_handler.go                # POST /user-info, POST /user-modify
│   │   ├── proposal_handler.go            # GET /public-proposals
│   │   ├── saved_proposal_handler.go      # POST + GET saved-proposals
│   │   └── context.go                     # AuthMiddleware + context helpers
│   └── router/router.go                   # All route wiring
└── README.md
```

---

## Milestone Progress

| # | Milestone | Status |
|---|-----------|--------|
| 0 | Scan datos.gov.co API | ✅ Done |
| 1 | Backend project skeleton + database schema | ✅ Done |
| 2 | JWT auth + user endpoints | ✅ Done |
| 3 | datos.gov.co live integration | ✅ Done |
| 4 | Endpoint logic refinements | ✅ Done |
| 5 | Frontend pages (login + dashboard) | ⏳ Pending |
| 6 | Frontend ↔ Backend wiring | ⏳ Pending |
| 7 | End-to-end validation | ⏳ Pending |

---

## Git History (Today — 2026-07-06)

```
215fb84 Building logic to endpoints for saved proposals, get user info, modify user and obtains user proposals saved
8e5468b  Build integration with api external and creates first version api contract
4e689f7  feat: add specs 1-9, integration docs, JWT backend
5aa5167  feat: add JWT auth, user endpoints, and real PostgreSQL repos
e1bed16  docs: add executive and technical report to SOUL.md
d68b6b1  feat: add backend skeleton with layered architecture
95256fb  Refactor sql schema
```

---

## Project File Tree

```
challenge_dev/
├── .gitignore
├── hermes.md                          # Product guidance
├── SOUL.md                            # This report
├── apicontract.md                     # API contract (8 endpoints)
├── backend/                           # Go backend (Phase 1 complete)
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
│   ├── public_calls_database.sql      # Production schema (+user_sessions)
│   └── compare_rubric/                # Schema comparison artifacts
├── integration_doc/
│   └── datosgovco/                    # SECOP II API docs + test
├── specs_project/                     # Task specs 1–10
└── .hermes/plans/                     # Implementation plan
```

---

## Next Steps (Phase 2 — Frontend)

1. **Build frontend pages** — `login.html`, `dashboard.html` with CSS sheets
2. **Create JavaScript API client** — `api.js` with `fetch()` wrappers
3. **Wire login flow** — login page → JWT storage → redirect to dashboard
4. **Build proposal browse UI** — search, filter, paginate, save proposals
5. **Connect real PostgreSQL** — swap in-memory mocks for pgx repos in `main.go`
6. **End-to-end validation** — full login → browse → filter → save → retrieve flow

---

## Risks & Open Questions

| Risk | Mitigation |
|------|------------|
| datos.gov.co API rate limits | Service has 30s timeout; pagination with `$limit`/`$offset` ready |
| Real DB not connected yet | `postgres_*.go` repos are ready; `database/postgres.go` has `NewPool()` |
| Static frontend (no framework) | Fine for MVP; can add React/Vue later if needed |
| No refresh token mechanism | Token is 1h; user must re-login after expiry |
| Mock data for saved proposals | Real DB integration needed for persistence across restarts |