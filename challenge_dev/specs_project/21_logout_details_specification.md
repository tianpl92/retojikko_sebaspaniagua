#Role
Full-stack developer web services — Go backend + HTML/CSS/JS frontend

#Context
The Public Calls Portal currently has no logout functionality. Once a user logs in, the JWT token remains valid for 1 hour with no way to proactively invalidate it. Users need a "Cerrar sesion" button to end their session and clear local credentials.

#Agents - Subagents
1. Principal agent to run plans, read documents, generate documentation, code review, use deepseek/deepseek-v4-flash model from openrouter.
2. Subagent to generate Go backend code and tests use deepseek/deepseek-v4-pro model from openrouter.
3. Subagent to generate HTML/CSS/JS frontend code use deepseek/deepseek-v4-pro model from openrouter.

---

#Previous requeriments
1. Shut down backend service if it is up and running.
2. Shut down frontend service if it is up and running.

# Architecture & Design

## Backend: POST /logout

**Protocol:** `POST /logout`
**Auth:** `Authorization: Bearer <token>` (required)
**Content-Type:** `application/json`

### Flow

```
Client                          Backend
  │                                │
  │  POST /logout                  │
  │  Authorization: Bearer <jwt>   │
  │ ─────────────────────────────► │
  │                                │ 1. AuthMiddleware validates JWT
  │                                │ 2. Extract token from Authorization header
  │                                │ 3. Delete session from user_sessions table
  │                                │ 4. Clear any server-side state
  │  {"message":"Sesión cerrada    │
  │   exitosamente"}               │
  │ ◄───────────────────────────── │
  │                                │
```

### Input

| Header | Value |
|--------|-------|
| `Authorization` | `Bearer <jwt-token>` |

No body parameters required.

### Response — 200 OK

```json
{
  "message": "Sesion cerrada exitosamente"
}
```

### Error Responses

| Status | Condition | Body |
|--------|-----------|------|
| 401 Unauthorized | Missing/invalid/expired token | `{"error":"Credentials invalid"}` |

### Backend file changes

| File | Action | Description |
|------|--------|-------------|
| `backend/internal/domain/user_session.go` | No change | Existing struct sufficient |
| `backend/internal/repository/session_repository.go` | **Add method** | Add `DeleteByToken(ctx, token) error` to the interface |
| `backend/internal/repository/mock_repository.go` | **Add method** | Implement `DeleteByToken` on `mockSessionRepository` — remove session from in-memory map by token |
| `backend/internal/repository/postgres_session_repository.go` | **Add method** | Implement `DeleteByToken` — `DELETE FROM user_sessions WHERE token = $1` |
| `backend/internal/service/auth_service.go` | **Add method** | Add `Logout(ctx, token) error` — calls `sessionRepo.DeleteByToken(ctx, token)` |
| `backend/internal/handler/auth_handler.go` | **Add handler** | Add `LogoutHandler` or extend existing handler with logout route |
| **OR** `backend/internal/handler/logout_handler.go` | **Create new file** | Dedicated handler for logout (preferred for separation of concerns) |
| `backend/internal/router/router.go` | **Add route** | `mux.Handle("/logout", AuthMiddleware(authService, logoutHandler))` |
| `backend/internal/handler/handler_test.go` | **Add tests** | Test logout: success, invalid token, missing token |
| `backend/internal/service/service_test.go` | **Add tests** | Test Logout service method |

### Handler implementation (pseudocode)

```go
// LogoutHandler handles POST /logout.
type LogoutHandler struct {
    authService *service.AuthService
}

func NewLogoutHandler(authService *service.AuthService) *LogoutHandler {
    return &LogoutHandler{authService: authService}
}

func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        writeError(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }

    // Extract token from Authorization header
    token := extractBearerToken(r)
    if token == "" {
        writeError(w, http.StatusUnauthorized, "Credentials invalid")
        return
    }

    if err := h.authService.Logout(r.Context(), token); err != nil {
        writeError(w, http.StatusInternalServerError, err.Error())
        return
    }

    writeJSON(w, http.StatusOK, map[string]string{
        "message": "Sesion cerrada exitosamente",
    })
}

func extractBearerToken(r *http.Request) string {
    auth := r.Header.Get("Authorization")
    if !strings.HasPrefix(auth, "Bearer ") {
        return ""
    }
    return strings.TrimSpace(auth[7:])
}
```

### Token extraction note

The `AuthMiddleware` already validates the JWT and sets `userID` in the context. However, for logout we need the **raw token string** to delete the session record. There are two approaches:

**Approach A (recommended):** Pass the token through context. Modify `AuthMiddleware` to also store the raw token string in the context (e.g., `context.WithValue(r.Context(), tokenKey, tokenString)`). Then the handler reads it via `r.Context().Value(tokenKey)`.

**Approach B:** Extract the token again from the `Authorization` header in the handler. Simpler but duplicates header parsing.

Recommended: **Approach A** — add a `tokenKey` context key (similar to existing `userIDKey` in `context.go`), set it in `AuthMiddleware`, read it in the logout handler.

### AuthMiddleware change

In `backend/internal/handler/context.go`:
```go
type contextKey string

const (
    userIDKey contextKey = "userID"
    tokenKey  contextKey = "token"  // NEW
)
```

In `AuthMiddleware`, after extracting the token string:
```go
ctx := context.WithValue(r.Context(), tokenKey, tokenString)
```

---

## Frontend: "Cerrar sesion" button

### File changes

| File | Action | Description |
|------|--------|-------------|
| `frontend/pages/dashboard.html` | **Modify** | Add "Cerrar sesion" button to the top navbar/header area |
| `frontend/scripts/dashboard.js` | **Modify** | Add `logout()` function, wire click event |
| `frontend/pages/saved-proposals.html` | **Modify** | Add "Cerrar sesion" button to the top navbar/header area |
| `frontend/scripts/saved-proposals.js` | **Modify** | Add `logout()` function, wire click event |

### UI Placement

Add a "Cerrar sesion" link/button in the top-right area of the page header, near the user display name. Style it as a subtle text link (not a primary button) using the secondary color palette:

```html
<button onclick="logout()" class="btn-logout" title="Cerrar sesión">
    Cerrar sesión
</button>
```

### CSS (in-page `<style>` block on each page)

```css
.btn-logout {
    background: none;
    border: 1px solid #6EADBC;
    color: #6EADBC;
    padding: 6px 16px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 13px;
    transition: all 0.2s;
}
.btn-logout:hover {
    background: #6EADBC;
    color: #fff;
}
```

### JavaScript `logout()` function (shared logic for both pages)

```javascript
async function logout() {
    const token = localStorage.getItem("jwt_token");
    if (!token) {
        window.location.href = "login.html";
        return;
    }

    try {
        await fetch(BACKEND_URL + "/logout", {
            method: "POST",
            headers: { Authorization: "Bearer " + token },
        });
    } catch (err) {
        // Proceed with local cleanup even if server request fails
    }

    localStorage.removeItem("jwt_token");
    localStorage.removeItem("user_name");
    localStorage.removeItem("user_last_name");
    window.location.href = "login.html";
}
```

### Expected UX flow

1. User clicks "Cerrar sesion"
2. Frontend calls `POST /logout` with JWT token (background, no loading overlay needed)
3. Backend invalidates the session (deletes from `user_sessions`)
4. Frontend clears `localStorage` (JWT + user data)
5. Redirect to `login.html`

---

# Definition of Done

1. `POST /logout` endpoint implemented and returns `200` with `{"message":"Sesion cerrada exitosamente"}`
2. Invalid/expired tokens return `401` with `{"error":"Credentials invalid"}`
3. Frontend "Cerrar sesion" button visible on both dashboard and saved-proposals pages
4. Clicking "Cerrar sesion" clears localStorage and redirects to login
5. After logout, the old token cannot be reused to access protected endpoints (returns 401)
6. All tests pass: `go build && go vet && go test ./... -count=1`
7. `gofmt` clean
8. API contract (`apicontract.md`) updated with new endpoint
9. Run up backend service.
10. Run up frontend service.

---

# Open questions

- Should logout also delete ALL sessions for the user (force re-login everywhere) or just the current session? ->  Just the current session (by token).
- Should the session row be hard-deleted or soft-deleted? -> Hard-delete (`DELETE FROM user_sessions WHERE token = $1`).
- Should the mock repository also validate the token exists before deleting? -> Yes — return `ErrNotFound` if token doesn't exist, but the handler should treat this as success (already logged out).

---

**Do not commit changes yet.**
