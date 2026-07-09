package handler

import (
	"context"
	"net/http"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/service"
)

// contextKey type for context values
type contextKey string

const userIDKey contextKey = "userID"
const tokenKey contextKey = "token"

// contextWithUserID stores the userID in context.
func contextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// contextWithToken stores the raw JWT token string in context.
func contextWithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey, token)
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

		// Store userID and raw token in context
		ctx := r.Context()
		ctx = contextWithUserID(ctx, userID)
		ctx = contextWithToken(ctx, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
