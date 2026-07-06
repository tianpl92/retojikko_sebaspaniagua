package handler

import (
	"context"
)

// contextWithUserID stores the userID in context.
func contextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
