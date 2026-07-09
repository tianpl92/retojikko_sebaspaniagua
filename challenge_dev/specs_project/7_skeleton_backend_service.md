#Context
Now we are build the skeleton backend service using golang.

#Role
Developer engineer api services

#Agents - Subagents
1.  Principal agent to run plans, read documents, generate documentation, read md instructions, use deepseek/deepseek-v4-flash model from openrouter.
2.  Subagent to genete go files , coding scripts and testing backend; use deepseek/deepseek-v4-pro model from openrouter.

#Actions

1. Use Go lastest version language.
2. Install go latest version into local machine if does not exists.
3. run go version command in shell to verify tool.
4. run go mod init using backend folder.
5. Update plan hermes to specify go version to use into tech information.

#Api server creation
**Objective:** Add the core business logic (skeleton) for login, browsing public proposals, filtering proposals, saving users, saving proposals, and retrieving saved proposals by user.

**Files:**
- Modify/Create: `backend/internal/service/*.go`
- Modify/Create: `backend/internal/handler/*.go`
- Modify/Create: `backend/internal/router/router.go`
- Create: `backend/internal/database/postgres.go`
- Create: `backend/internal/domain/*.go`
- Create: `backend/internal/repository/*.go`
- Create: `backend/internal/handler/auth_handler_test.go`
- Create: `backend/internal/handler/public__handler_test.go`
- Create: `backend/internal/service/public_proposals_service_test.go`
- Create: `backend/internal/service/auth_service_test.go`

**Step 1: Write failing tests for each endpoint**

Cover at least these cases:
- `POST /login` accepts valid credentials and rejects invalid ones
- `GET /public-proposals` returns a list
- `GET /public-proposals?query=...&category=...` filters results
- `POST /saved-proposals` stores a selection for the current user
- `GET /saved-proposals` returns the user’s saved items

**Step 2: Implement service logic incrementally**

Add the smallest code needed to pass each test:
- auth service for login/session token handling
- proposal service for listing/filtering
- saved-proposal service for persistence

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


----------


Tell me first of all, just to be sure, what do you plan to do?
