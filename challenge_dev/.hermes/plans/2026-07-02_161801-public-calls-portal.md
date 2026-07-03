# Public Calls Portal Implementation Plan

> **For Hermes:** Use subagent-driven-development skill to implement this plan task-by-task.

**Goal:** Build a web application for browsing, filtering, and saving public calls for proposals from datos.gov.co, with a Go backend, PostgreSQL persistence, and an HTML5/JavaScript frontend.

**Architecture:** Start with a backend-first MVP that exposes REST JSON endpoints for authentication, public-call browsing, and saving favorites. Keep the backend layered as described in `hermes.md` (database/repository/service/handler/router) and serve a separate static frontend that consumes those endpoints with `fetch`. Use PostgreSQL for persistence and keep the frontend dependency-light unless a future requirement justifies a build tool.

**Tech Stack:** Go, PostgreSQL, HTML5, CSS, vanilla JavaScript, REST/JSON.

---

## Current context

- The repository currently only has `hermes.md` as product guidance.
- `hermes.md` defines a two-phase workflow:
  1. Backend first
  2. Frontend second
- The target app is the “Public Calls Portal”, which should connect to datos.gov.co services and allow users to browse, filter, and save public calls.
- The repo structure described in `hermes.md` should be created from scratch.

## Assumptions

- The first deliverable is an MVP, not the full production portal.
- Authentication will be implemented minimally unless the external requirements specify a stronger identity provider.
- The frontend will be a static HTML/CSS/JS app that talks directly to the backend REST API.
- The datos.gov.co response shape may need to be adapted once the exact external API endpoint is confirmed.

## Proposed milestones

0. Scan `datos.gov.co` and identify public service documentation for available endpoints related to **Convocatorias publicas**.
1. Create the backend project skeleton and database schema.
2. Implement core backend endpoints and business logic.
3. Add backend unit tests and verify the API works end-to-end.
4. Create the frontend pages and wire them to the backend.
5. Validate the complete flow: login, browse, filter, save, and retrieve saved calls.

---

### Task 1: Scan datos.gov.co

**Objective:** Discover and document the official public-service endpoint documentation on `datos.gov.co` that can be used to search and retrieve **Convocatorias publicas** information.

**Files:**
- Modify: `.hermes/plans/2026-07-02_161801-public-calls-portal.md` (planning only; no code changes)
- Potentially create later: `backend/internal/repository/datosgov_repository.go`
- Potentially create later: `backend/internal/service/public_call_sync_service.go`

**Best-practice guide:**
1. Start from the official `datos.gov.co` site and look for documentation pages, API reference pages, or developer/help sections.
2. Search for terms such as `Convocatorias publicas`, `convocatorias`, `public calls`, `API`, `endpoint`, `servicio`, `documentación`, and `datos abiertos`.
3. Prefer official documentation over blog posts, mirrors, or third-party summaries.
4. Capture for each candidate endpoint:
   - base URL
   - HTTP method
   - query parameters
   - pagination scheme
   - authentication requirements, if any
   - response format and example fields
5. Confirm whether the endpoint supports filtering by text, category, status, date, or issuer metadata.
6. Note any limits, rate limits, or access restrictions that could affect the backend design.
7. If multiple endpoints exist, rank them by suitability for the MVP and select the simplest one that satisfies the portal needs.
8. Record the findings in the plan before implementing backend integration.

**Expected outcome:**
- A short discovery note identifying the best `datos.gov.co` endpoint candidate for `Convocatorias publicas`.
- Enough API detail to design the repository/service layer without guessing.

---

### Task 2: Create the repository skeleton and backend bootstrap

**Objective:** Establish the folder layout, Go module, and application entrypoint so backend development can proceed cleanly.

**Files:**
- Create: `backend/go.mod`
- Create: `backend/cmd/api/main.go`
- Create: `backend/internal/config/config.go`
- Create: `backend/internal/router/router.go`
- Create: `backend/internal/handler/health_handler.go`
- Create: `backend/internal/service/health_service.go`
- Create: `backend/internal/repository/`
- Create: `backend/internal/database/`
- Create: `backend/internal/domain/`
- Create: `backend/migrations/`
- Create: `frontend/pages/`
- Create: `frontend/scripts/`
- Create: `frontend/sheets/`

**Step 1: Define the minimum backend bootstrap test target**

Add a simple health check endpoint contract in a future test file:
- `backend/internal/handler/health_handler_test.go`

**Step 2: Implement the minimal backend bootstrap**

Create a Go HTTP server entrypoint that:
- reads the port from config/env
- registers a `/health` route
- starts listening with `net/http`

**Step 3: Verify the server starts**

Run:
- `cd backend && go test ./...`
- `cd backend && go run ./cmd/api`

Expected:
- tests compile even before business logic exists
- `GET /health` returns `200 OK` with a small JSON body

**Step 4: Commit the scaffold**

Commit after the backend skeleton is stable.

---

### Task 3: Design the PostgreSQL schema for users, public calls, and saved items

**Objective:** Create the database model that supports browsing and saving calls.

**Files:**
- Create: `backend/migrations/001_init.sql`
- Create: `backend/internal/domain/user.go`
- Create: `backend/internal/domain/public_call.go`
- Create: `backend/internal/domain/saved_call.go`
- Create: `backend/internal/repository/user_repository.go`
- Create: `backend/internal/repository/public_call_repository.go`
- Create: `backend/internal/repository/saved_call_repository.go`

**Step 1: Write a schema-focused test or review checklist**

Define expected entities and constraints:
- users table with unique email
- public_calls table for cached external calls
- saved_calls join table linking users and public calls
- timestamps for auditing

**Step 2: Implement the migration**

Add SQL for:
- `users`
- `public_calls`
- `saved_calls`
- primary keys, foreign keys, unique constraints, and indexes

**Step 3: Verify the migration is syntactically valid**

Run the migration through your chosen local Postgres workflow or SQL parser.

Expected:
- schema applies cleanly to an empty database
- unique and foreign-key constraints behave as intended

**Step 4: Commit the database foundation**

---

### Task 4: Implement backend services and REST endpoints

**Objective:** Add the core business logic for login, browsing calls, filtering calls, saving calls, and retrieving saved calls.

**Files:**
- Modify/Create: `backend/internal/service/*.go`
- Modify/Create: `backend/internal/handler/*.go`
- Modify/Create: `backend/internal/router/router.go`
- Create: `backend/internal/database/postgres.go`
- Create: `backend/internal/domain/*.go`
- Create: `backend/internal/repository/*.go`
- Create: `backend/internal/handler/auth_handler_test.go`
- Create: `backend/internal/handler/public_call_handler_test.go`
- Create: `backend/internal/service/public_call_service_test.go`
- Create: `backend/internal/service/auth_service_test.go`

**Step 1: Write failing tests for each endpoint**

Cover at least these cases:
- `POST /login` accepts valid credentials and rejects invalid ones
- `GET /public-calls` returns a list
- `GET /public-calls?query=...&category=...` filters results
- `POST /saved-calls` stores a selection for the current user
- `GET /saved-calls` returns the user’s saved items

**Step 2: Implement service logic incrementally**

Add the smallest code needed to pass each test:
- auth service for login/session token handling
- call service for listing/filtering
- saved-call service for persistence

**Step 3: Wire handlers to repositories**

Ensure handlers:
- decode JSON input safely
- validate required fields
- return consistent JSON error shapes
- map domain errors to correct HTTP status codes

**Step 4: Verify with backend tests**

Run:
- `cd backend && go test ./...`

Expected:
- all handler and service tests pass
- `gofmt` produces no changes

**Step 5: Commit the backend feature slice**

---

### Task 5: Integrate datos.gov.co as the external source of calls

**Objective:** Fetch public calls from datos.gov.co and normalize them into the app’s internal model.

**Files:**
- Create/Modify: `backend/internal/repository/datosgov_repository.go`
- Create/Modify: `backend/internal/service/public_call_sync_service.go`
- Create/Modify: `backend/internal/domain/public_call.go`
- Create: `backend/internal/repository/datosgov_repository_test.go`
- Create: `backend/internal/service/public_call_sync_service_test.go`

**Step 1: Lock down the external contract in tests**

Use a mocked HTTP response to define:
- the expected datos.gov.co request URL/query parameters
- the fields needed for portal display
- how pagination is handled

**Step 2: Implement normalization**

Convert external API records into the internal `public_calls` shape.

**Step 3: Add sync or fetch behavior**

Decide whether calls are:
- fetched live on request, or
- periodically synced and cached in PostgreSQL

Prefer the simplest approach that still satisfies the portal requirements.

**Step 4: Verify**

Run the repository/service tests and confirm the app handles external API failures gracefully.

**Step 5: Commit the integration layer**

---

### Task 6: Build the frontend login and dashboard pages

**Objective:** Create the visual shell of the portal and the navigation users need to reach core features.

**Files:**
- Create: `frontend/pages/login.html`
- Create: `frontend/pages/dashboard.html`
- Create: `frontend/sheets/login.css`
- Create: `frontend/sheets/dashboard.css`
- Create: `frontend/sheets/styles.css`
- Create: `frontend/scripts/api.js`
- Create: `frontend/scripts/login.js`
- Create: `frontend/scripts/dashboard.js`

**Step 1: Write the page structure**

Implement:
- a login form
- a dashboard with submenu/navigation
- containers for call listings, saved calls, and filters

**Step 2: Wire the frontend to the backend API**

Use `fetch` wrappers in `api.js` for:
- login
- listing calls
- saving a call
- loading saved calls

**Step 3: Add basic styling**

Ensure the layout is readable and responsive enough for the first release.

**Step 4: Verify in a browser**

Manual validation:
- login page loads
- dashboard loads after login
- data renders in the page
- failures show a useful error message

**Step 5: Commit the frontend shell**

---

### Task 6: Connect frontend actions to real backend data flows

**Objective:** Make browse, filter, and save actions work end-to-end.

**Files:**
- Modify: `frontend/scripts/dashboard.js`
- Modify: `frontend/scripts/api.js`
- Modify: `frontend/pages/dashboard.html`
- Modify: `backend/internal/handler/public_call_handler.go`
- Modify: `backend/internal/handler/saved_call_handler.go`

**Step 1: Add UI actions and empty states**

Cover:
- search box
- category/status/date filters
- save and unsave actions
- “no results” state

**Step 2: Implement backend response shapes needed by the UI**

Return payloads that are easy for the frontend to render without extra transformation.

**Step 3: Validate end-to-end manually**

Confirm:
- filters update the list
- saved items persist after refresh
- errors from backend are surfaced in the UI

**Step 4: Commit the integrated slice**

---

## Validation checklist

Before considering the MVP ready, verify all of the following:

- `cd backend && go test ./...` passes
- `gofmt -w` on backend files is clean
- the backend serves `/health` and the public-call endpoints
- PostgreSQL schema applies cleanly to a fresh database
- frontend pages load without console errors
- login, browse, filter, save, and saved-items retrieval work end-to-end

## Risks and tradeoffs

- The external datos.gov.co API shape may differ from what the UI needs, so normalization logic should stay isolated.
- Auth requirements are not fully specified; the first pass may need to be replaced once real identity requirements are known.
- A static frontend keeps the first release simple, but if the app grows, a build tool or framework may become worthwhile.

## Open questions

- What exact datos.gov.co endpoint(s) should be used?
- Should login be real authentication or a minimal portal session for the MVP?
- Do saved calls belong to authenticated users only?
- Should the backend cache public-call data locally or query the external API live on each request?

---

## Delivery order

Implement in this order:
1. Backend scaffold
2. Database schema
3. Core API endpoints
4. External API integration
5. Frontend pages
6. Frontend/backend wiring
7. End-to-end validation
