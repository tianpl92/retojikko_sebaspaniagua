package handler

import (
	"net/http"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/service"
)

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

	// Read token from context (stored by AuthMiddleware)
	token, ok := r.Context().Value(tokenKey).(string)
	if !ok || token == "" {
		writeError(w, http.StatusUnauthorized, "Credentials invalid")
		return
	}

	// Invalidate the session (best-effort: ignore error if already deleted/expired)
	_ = h.authService.Logout(r.Context(), token)

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Sesion cerrada exitosamente",
	})
}