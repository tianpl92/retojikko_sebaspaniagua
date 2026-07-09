package handler

import (
	"context"
	"net/http"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/service"
)

// contextKey type for context values
type contextKey string

const userIDKey contextKey = "userID"

// contextWithUserID stores the userID in context.
func contextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// AuthMiddleware checks token from Authorization header and injects userID into context.
func AuthMiddleware(authService *service.AuthService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			writeError(w, http.StatusUnauthorized, "Credentials invalid")
			return
		}
		// Strip "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		} else {
			writeError(w, http.StatusUnauthorized, "Credentials invalid")
			return
		}

		userID, err := authService.ValidateToken(token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "Credentials invalid")
			return
		}

		// Store userID in context
		ctx := r.Context()
		ctx = contextWithUserID(ctx, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
